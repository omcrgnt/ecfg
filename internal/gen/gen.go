package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/omcrgnt/ecfg/internal/ecfgtool"
)

// Options configures code generation output paths.
type Options struct {
	TemplatePath string
	MarkdownPath string
}

// Run generates env files for the given root config type.
func Run(typeName, pkgPath, prefix string, opts Options) error {
	entries, err := ecfgtool.CollectTemplateEntries(pkgPath, typeName, prefix)
	if err != nil {
		return err
	}
	if opts.TemplatePath != "" {
		path := absPath(opts.TemplatePath)
		if err := ecfgtool.WriteEnvFile(path, entries); err != nil {
			return fmt.Errorf("write template: %w", err)
		}
	}
	if opts.MarkdownPath != "" {
		path := absPath(opts.MarkdownPath)
		if err := ecfgtool.WriteEnvMarkdown(path, prefix, entries); err != nil {
			return fmt.Errorf("write markdown: %w", err)
		}
	}
	return nil
}

func absPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if wd, err := os.Getwd(); err == nil {
		return filepath.Join(wd, path)
	}
	return path
}
