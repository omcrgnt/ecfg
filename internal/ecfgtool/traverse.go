package ecfgtool

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

type visitCtx struct {
	walk      walk.VisitCtx
	value     reflect.Value
	envKey    string
	fieldPath string
	usageText string
	rootGroup string
	isProto   bool
	isLeaf    bool
}

func traverse(engine walk.Engine, opts Options, visit func(visitCtx) error) error {
	path := make([]string, 0, 8)
	registry := newKeyRegistry()

	walkOpts := walk.Options{
		InitPointers: true,
		AfterField: func(wctx walk.VisitCtx) {
			path = path[:len(path)-1]
		},
	}
	return walk.StructWalk(engine, walkOpts, func(wctx walk.VisitCtx) error {
		f := wctx.Field
		ecfgTag := parseEcfgTag(f.Tag)
		seg, err := segment(wctx.Depth, ecfgTag, f.Name)
		if err != nil {
			return err
		}
		path = append(path, seg)

		isProto := isProtoField(f)

		vctx := visitCtx{
			walk:      wctx,
			fieldPath: joinFieldPath(path, f.Name),
		}

		if err := check(vctx, isProto); err != nil {
			return err
		}

		if engR, ok := wctx.Engine.(*walk.EngineReflect); ok {
			if err := engR.InitPointerField(f); err != nil {
				return err
			}
		}

		vctx.isProto = isProto
		vctx.isLeaf = isLeaf(wctx.Depth, f, isProto)

		if vctx.isLeaf {
			vctx.envKey = join(opts.Prefix, path...)
			if err := registry.add(vctx.envKey, vctx.fieldPath); err != nil {
				return err
			}
			if engR, ok := wctx.Engine.(*walk.EngineReflect); ok {
				val, _, err := engR.FieldValue(f)
				if err != nil {
					return err
				}
				vctx.value = val
			}
			usageText, err := resolveUsage(f, isProto)
			if err != nil {
				return err
			}
			vctx.usageText = usageText
			if len(path) > 0 {
				vctx.rootGroup = path[0]
			}
		}

		if err := visit(vctx); err != nil {
			return err
		}

		if vctx.isLeaf && isProto {
			return walk.SkipDescend()
		}
		return nil
	})
}

func isProtoField(f walk.FieldDesc) bool {
	if f.ReflectType != nil {
		return isProtoMessage(f.ReflectType)
	}
	if f.TypesType != nil {
		return isProtoTypesType(f.TypesType)
	}
	return false
}

func joinFieldPath(path []string, field string) string {
	if len(path) == 0 {
		return field
	}
	return fmt.Sprintf("%s.%s", joinFieldPathSegments(path), field)
}

func joinFieldPathSegments(path []string) string {
	out := path[0]
	for i := 1; i < len(path); i++ {
		out += "." + path[i]
	}
	return out
}
