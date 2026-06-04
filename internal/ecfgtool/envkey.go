package ecfgtool

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

func segment(depth int, tag, fieldName string) (string, error) {
	if depth == 0 {
		if tag == "" {
			return "", ErrMissingEcfgTag
		}
		return strings.ToUpper(tag), nil
	}
	if tag != "" {
		return strings.ToUpper(tag), nil
	}
	return fieldNameToSegment(fieldName), nil
}

func fieldNameToSegment(name string) string {
	if name == "" {
		return ""
	}
	var b strings.Builder
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(name[i-1])
				if unicode.IsLower(prev) {
					b.WriteByte('_')
				} else if i+1 < len(name) && unicode.IsUpper(prev) && unicode.IsLower(rune(name[i+1])) {
					b.WriteByte('_')
				}
			}
			b.WriteRune(r)
		} else {
			b.WriteRune(unicode.ToUpper(r))
		}
	}
	return b.String()
}

func join(prefix string, segments ...string) string {
	parts := make([]string, 0, len(segments)+1)
	if p := strings.Trim(strings.TrimSpace(prefix), "_"); p != "" {
		parts = append(parts, strings.ToUpper(p))
	}
	for _, s := range segments {
		if s = strings.TrimSpace(s); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "_")
}

type keyRegistry struct {
	keys map[string]string
}

func newKeyRegistry() *keyRegistry {
	return &keyRegistry{keys: make(map[string]string)}
}

func (r *keyRegistry) add(fullKey, fieldPath string) error {
	if prev, ok := r.keys[fullKey]; ok {
		return fmt.Errorf("%w: %s (%s and %s)", ErrDuplicateEnvKey, fullKey, prev, fieldPath)
	}
	r.keys[fullKey] = fieldPath
	return nil
}

func isLeaf(depth int, f walk.FieldDesc, isProto bool) bool {
	if isProto {
		return true
	}
	if isStructField(f) {
		return false
	}
	return depth >= 1
}

func isStructField(f walk.FieldDesc) bool {
	if f.ReflectType != nil {
		k := f.ReflectType.Kind()
		return k == reflect.Struct || k == reflect.Ptr
	}
	if f.TypesType != nil {
		k := walk.ReflectKind(f.TypesType)
		return k == reflect.Struct || k == reflect.Ptr
	}
	return false
}
