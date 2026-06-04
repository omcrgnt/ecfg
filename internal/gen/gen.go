package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/omcrgnt/ecfg/internal/ecfgtool"
)

// Run generates env.template for the given root config type.
func Run(typeName, pkgPath, prefix, outPath string) error {
	entries, err := ecfgtool.CollectTemplateEntries(pkgPath, typeName, prefix)
	if err != nil {
		return err
	}
	if !filepath.IsAbs(outPath) {
		if wd, err := os.Getwd(); err == nil {
			outPath = filepath.Join(wd, outPath)
		}
	}
	if err := ecfgtool.WriteEnvTemplate(outPath, entries); err != nil {
		return fmt.Errorf("write template: %w", err)
	}
	return nil
}
