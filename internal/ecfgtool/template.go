package ecfgtool

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

// TemplateEntry is one env.template block line pair.
type TemplateEntry struct {
	EnvKey    string
	Usage     string
	RootGroup string
}

// CollectTemplateEntries walks types engine and collects template lines.
func CollectTemplateEntries(pkgPath, typeName, prefix string) ([]TemplateEntry, error) {
	eng, err := walk.NewEngineTypes(pkgPath, typeName)
	if err != nil {
		return nil, err
	}
	var entries []TemplateEntry
	opts := Options{Prefix: prefix}
	err = traverse(eng, opts, func(ctx visitCtx) error {
		if !ctx.isLeaf {
			return nil
		}
		entries = append(entries, TemplateEntry{
			EnvKey:    ctx.envKey,
			Usage:     ctx.usageText,
			RootGroup: ctx.rootGroup,
		})
		return nil
	})
	return entries, err
}

// WriteEnvTemplate writes env.template entries grouped by root block.
func WriteEnvTemplate(path string, entries []TemplateEntry) error {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].RootGroup != entries[j].RootGroup {
			return entries[i].RootGroup < entries[j].RootGroup
		}
		return entries[i].EnvKey < entries[j].EnvKey
	})
	var b strings.Builder
	prevRoot := ""
	for i, e := range entries {
		if i > 0 && e.RootGroup != prevRoot {
			b.WriteByte('\n')
		}
		prevRoot = e.RootGroup
		fmt.Fprintf(&b, "# %s\n%s=\n", e.Usage, e.EnvKey)
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
