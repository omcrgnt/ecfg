package ecfgtool

import (
	"github.com/omcrgnt/ecfg/pkg/walk"
)

// traverseRootBlocks walks ecfg-tagged first-level blocks. When eng supplies a builder
// spec (or nested config block), rootSeg is the AppResources ecfg tag (e.g. APP, SERVICE_ITEM).
func traverseRootBlocks(rootSeg, rootField string, eng walk.Engine, opts Options, visit func(visitCtx) error) error {
	path := []string{rootSeg}
	registry := newKeyRegistry()
	const depthOffset = 1

	walkOpts := walk.Options{
		InitPointers: true,
		AfterField: func(wctx walk.VisitCtx) {
			if len(path) > 1 {
				path = path[:len(path)-1]
			}
		},
	}
	return walk.StructWalk(eng, walkOpts, func(wctx walk.VisitCtx) error {
		f := wctx.Field
		wctx.Depth += depthOffset

		seg, err := segment(wctx.Depth, parseEcfgTag(f.Tag), f.Name)
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
			vctx.rootGroup = rootSeg
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
