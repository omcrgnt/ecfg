package ecfg

import (
	"fmt"
	"reflect"
	"strings"
)

// CatalogSegment returns the env block segment from a catalog wire field struct tag.
// The tag name is [TagKey] (e.g. ecfg:"APP", or cfg:"APP" after [SetTagKey]).
func CatalogSegment(sf reflect.StructField) (string, error) {
	key := TagKey()
	tag, ok := sf.Tag.Lookup(key)
	if !ok {
		return "", fmt.Errorf("ecfg: missing %q tag", key)
	}
	seg := strings.ToUpper(strings.Split(tag, ",")[0])
	if seg == "" {
		return "", fmt.Errorf("ecfg: empty %q tag", key)
	}
	return seg, nil
}
