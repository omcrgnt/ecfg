// Command ecfg-gen writes env.template from ecfg struct tags (offline, no ENV).
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/omcrgnt/ecfg/internal/gen"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet("ecfg-gen", flag.ExitOnError)
	typeName := fs.String("type", "", "root config struct name (required)")
	pkgPath := fs.String("pkg", "", "package import path (required)")
	prefix := fs.String("prefix", "", "env key prefix")
	outPath := fs.String("o", "env.template", "output file path")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *typeName == "" || *pkgPath == "" {
		fmt.Fprintln(os.Stderr, "usage: ecfg-gen -type AppConfig -pkg github.com/you/app/config [-prefix APP] [-o env.template]")
		fs.PrintDefaults()
		return 2
	}
	if err := gen.Run(*typeName, *pkgPath, *prefix, *outPath); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
