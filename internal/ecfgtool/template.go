package ecfgtool

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// TemplateEntry is one generated env variable.
type TemplateEntry struct {
	EnvKey    string
	Usage     string
	RootGroup string
}

// CollectTemplateEntries walks ecfg-tagged AppResources fields.
// For [builder.BuildConfiger] resources it uses the spec type from [BuildConfig]; otherwise the field type.
func CollectTemplateEntries(pkgPath, typeName, prefix string) ([]TemplateEntry, error) {
	_, st, all, err := loadRootStruct(pkgPath, typeName)
	if err != nil {
		return nil, err
	}

	var entries []TemplateEntry
	opts := Options{Prefix: prefix, SkipUntaggedRoot: true}

	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Exported() {
			continue
		}
		ecfgTag := parseEcfgTag(st.Tag(i))
		if opts.SkipUntaggedRoot && ecfgTag == "" {
			continue
		}

		eng, err := engineForRootField(all, f)
		if err != nil {
			return nil, fmt.Errorf("ecfg: %s: %w", f.Name(), err)
		}

		if err := traverseRootBlocks(ecfgTag, f.Name(), eng, opts, func(ctx visitCtx) error {
			if !ctx.isLeaf {
				return nil
			}
			entries = append(entries, TemplateEntry{
				EnvKey:    ctx.envKey,
				Usage:     ctx.usageText,
				RootGroup: ctx.rootGroup,
			})
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return entries, nil
}

func sortTemplateEntries(entries []TemplateEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].RootGroup != entries[j].RootGroup {
			return entries[i].RootGroup < entries[j].RootGroup
		}
		return entries[i].EnvKey < entries[j].EnvKey
	})
}

// WriteEnvFile writes KEY= lines grouped by root block (no usage comments).
func WriteEnvFile(path string, entries []TemplateEntry) error {
	sortTemplateEntries(entries)
	var b strings.Builder
	prevRoot := ""
	for i, e := range entries {
		if i > 0 && e.RootGroup != prevRoot {
			b.WriteByte('\n')
		}
		prevRoot = e.RootGroup
		fmt.Fprintf(&b, "%s=\n", e.EnvKey)
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// WriteEnvMarkdown writes env.md with usage tables grouped by root block.
func WriteEnvMarkdown(path, prefix string, entries []TemplateEntry) error {
	sortTemplateEntries(entries)
	var b strings.Builder
	fmt.Fprintln(&b, "# Environment configuration")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Copy `.env.template` to `.env` and set values. See `.env.example` for local defaults.")
	fmt.Fprintln(&b)
	if p := strings.TrimSpace(prefix); p != "" {
		fmt.Fprintf(&b, "**Prefix:** `%s`\n\n", strings.ToUpper(p))
	}
	prevRoot := ""
	for _, e := range entries {
		if e.RootGroup != prevRoot {
			if prevRoot != "" {
				fmt.Fprintln(&b)
			}
			fmt.Fprintf(&b, "## %s\n\n", e.RootGroup)
			fmt.Fprintln(&b, "| Variable | Description |")
			fmt.Fprintln(&b, "|----------|-------------|")
			prevRoot = e.RootGroup
		}
		usage := strings.ReplaceAll(e.Usage, "|", "\\|")
		fmt.Fprintf(&b, "| `%s` | %s |\n", e.EnvKey, usage)
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// WriteEnvTemplate is an alias for [WriteEnvFile].
func WriteEnvTemplate(path string, entries []TemplateEntry) error {
	return WriteEnvFile(path, entries)
}
