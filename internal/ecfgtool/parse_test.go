package ecfgtool

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/omcrgnt/ecfg/internal/testdata"
)

func TestParse_happyWithPrefix(t *testing.T) {
	t.Setenv("TEST_SERVER_LABEL", "myapp")
	cfg, err := Parse[testdata.AppConfig](Options{Prefix: "TEST"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Label != "myapp" {
		t.Fatalf("label: got %q", cfg.Server.Label)
	}
}

func TestParse_happyWithoutPrefix(t *testing.T) {
	t.Setenv("SERVER_LABEL", "noprefix")
	cfg, err := Parse[testdata.AppConfig](Options{})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Label != "noprefix" {
		t.Fatalf("label: got %q", cfg.Server.Label)
	}
}

func TestParse_missingEnv(t *testing.T) {
	os.Unsetenv("TEST_SERVER_LABEL")
	_, err := Parse[testdata.AppConfig](Options{Prefix: "TEST"})
	if !errors.Is(err, ErrMissingEnv) {
		t.Fatalf("got %v want ErrMissingEnv", err)
	}
}

func TestParse_emptyEnv(t *testing.T) {
	t.Setenv("TEST_SERVER_LABEL", "   ")
	_, err := Parse[testdata.AppConfig](Options{Prefix: "TEST"})
	if !errors.Is(err, ErrEmptyEnv) {
		t.Fatalf("got %v want ErrEmptyEnv", err)
	}
}

func TestParse_duplicateEnvKey(t *testing.T) {
	t.Setenv("T_BLOCK_LABEL", "x")
	_, err := Parse[testdata.BadDup](Options{Prefix: "T"})
	if !errors.Is(err, ErrDuplicateEnvKey) {
		t.Fatalf("got %v want ErrDuplicateEnvKey", err)
	}
}

func TestParse_validateError(t *testing.T) {
	t.Setenv("T_BLOCK_LABEL", "reject")
	_, err := Parse[testdata.StrictCfg](Options{Prefix: "T"})
	if err == nil || !strings.Contains(err.Error(), "validate:") {
		t.Fatalf("got %v", err)
	}
}

func TestParse_wrapFieldErr(t *testing.T) {
	t.Setenv("T_BLOCK_PORT", "not-a-number")
	_, err := Parse[testdata.IntCfg](Options{Prefix: "T"})
	if err == nil {
		t.Fatal("expected error")
	}
	s := err.Error()
	if !strings.Contains(s, "BLOCK.PORT") || !strings.Contains(s, "T_BLOCK_PORT") {
		t.Fatalf("got %q", s)
	}
}

func TestParse_emptyUsage(t *testing.T) {
	t.Setenv("T_BLOCK_LABEL", "ok")
	_, err := Parse[testdata.EmptyUsageCfg](Options{Prefix: "T"})
	if !errors.Is(err, ErrEmptyUsage) {
		t.Fatalf("got %v want ErrEmptyUsage", err)
	}
}

func TestParse_policyErrors(t *testing.T) {
	cases := []struct {
		name string
		run  func() error
		want error
	}{
		{
			name: "root_not_struct",
			run: func() error {
				_, err := Parse[testdata.BadRootInt](Options{Prefix: "T"})
				return err
			},
			want: ErrRootFieldKind,
		},
		{
			name: "missing_ecfg_tag",
			run: func() error {
				_, err := Parse[testdata.BadNoTag](Options{})
				return err
			},
			want: ErrMissingEcfgTag,
		},
		{
			name: "nested_block",
			run: func() error {
				_, err := Parse[testdata.BadNested](Options{Prefix: "T"})
				return err
			},
			want: ErrNestedBlockAtDepth1,
		},
		{
			name: "slice",
			run: func() error {
				_, err := Parse[testdata.BadSlice](Options{Prefix: "T"})
				return err
			},
			want: ErrUnsupportedContainer,
		},
		{
			name: "bare_string_leaf",
			run: func() error {
				t.Setenv("T_BLOCK_NAME", "x")
				_, err := Parse[testdata.BadBare](Options{Prefix: "T"})
				return err
			},
			want: ErrIncompatibleLeaf,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v want %v", err, tc.want)
			}
		})
	}
}

func TestParse_notStructTarget(t *testing.T) {
	_, err := Parse[int](Options{})
	if !errors.Is(err, ErrNotStruct) {
		t.Fatalf("got %v want ErrNotStruct", err)
	}
}

func TestParse_intLeaf(t *testing.T) {
	t.Setenv("T_BLOCK_PORT", "8080")
	cfg, err := Parse[testdata.IntCfg](Options{Prefix: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Block.Port != 8080 {
		t.Fatalf("port: %d", cfg.Block.Port)
	}
}

func TestParse_intLeaf_invalid(t *testing.T) {
	t.Setenv("T_BLOCK_PORT", "not-a-number")
	_, err := Parse[testdata.IntCfg](Options{Prefix: "T"})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParse_intLeaf_validate(t *testing.T) {
	t.Setenv("T_BLOCK_PORT", "0")
	_, err := Parse[testdata.IntCfg](Options{Prefix: "T"})
	if err == nil || !strings.Contains(err.Error(), "validate:") {
		t.Fatalf("got %v", err)
	}
}

func TestParse_boolLeaf(t *testing.T) {
	t.Setenv("T_BLOCK_ENABLED", "true")
	cfg, err := Parse[testdata.BoolCfg](Options{Prefix: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if !bool(cfg.Block.Enabled) {
		t.Fatal("expected true")
	}
}

func TestParse_floatLeaf(t *testing.T) {
	t.Setenv("T_BLOCK_SCORE", "3.14")
	cfg, err := Parse[testdata.FloatCfg](Options{Prefix: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if float64(cfg.Block.Score) != 3.14 {
		t.Fatalf("score: %v", cfg.Block.Score)
	}
}

func TestParse_floatLeaf_invalid(t *testing.T) {
	t.Setenv("T_BLOCK_SCORE", "nanana")
	_, err := Parse[testdata.FloatCfg](Options{Prefix: "T"})
	if err == nil || !strings.Contains(err.Error(), "parse float") {
		t.Fatalf("got %v", err)
	}
}

func TestParse_durationLeaf(t *testing.T) {
	t.Setenv("T_BLOCK_TIMEOUT", "30s")
	cfg, err := Parse[testdata.DurationCfg](Options{Prefix: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if time.Duration(cfg.Block.Timeout) != 30*time.Second {
		t.Fatalf("timeout: %v", cfg.Block.Timeout)
	}
}

func TestParse_durationLeaf_invalid(t *testing.T) {
	t.Setenv("T_BLOCK_TIMEOUT", "bad")
	_, err := Parse[testdata.DurationCfg](Options{Prefix: "T"})
	if err == nil || !strings.Contains(err.Error(), "parse duration") {
		t.Fatalf("got %v", err)
	}
}

func TestParse_unsupportedScalarLeaf(t *testing.T) {
	t.Setenv("T_BLOCK_RATIO", "1")
	_, err := Parse[testdata.RatioCfg](Options{Prefix: "T"})
	if err == nil || !strings.Contains(err.Error(), "unsupported scalar kind") {
		t.Fatalf("got %v", err)
	}
}

func TestParse_protoLeaf_happy(t *testing.T) {
	t.Setenv("T_SERVER_PORT", "443")
	_, err := Parse[testdata.ProtoCfg](Options{Prefix: "T"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParse_protoLeaf_invalidValue(t *testing.T) {
	t.Setenv("T_SERVER_PORT", "999999")
	_, err := Parse[testdata.ProtoCfg](Options{Prefix: "T"})
	if err == nil {
		t.Fatal("expected protovalidate error")
	}
	if !strings.Contains(err.Error(), "validate:") {
		t.Fatalf("got %v", err)
	}
}

func TestParse_protoLeaf_missingEnv(t *testing.T) {
	os.Unsetenv("T_SERVER_PORT")
	_, err := Parse[testdata.ProtoCfg](Options{Prefix: "T"})
	if !errors.Is(err, ErrMissingEnv) {
		t.Fatalf("got %v", err)
	}
}

func TestParse_protoLeaf_missingUsage(t *testing.T) {
	t.Setenv("T_SERVER_PORT", "42")
	_, err := Parse[testdata.ProtoNoUsageCfg](Options{Prefix: "T"})
	if !errors.Is(err, ErrMissingUsage) {
		t.Fatalf("got %v want ErrMissingUsage", err)
	}
}

func TestParse_badProtoWrapper_nestedStruct(t *testing.T) {
	t.Setenv("T_SERVER_WRAP", "1")
	_, err := Parse[testdata.BadProtoWrapper](Options{Prefix: "T"})
	if !errors.Is(err, ErrNestedBlockAtDepth1) {
		t.Fatalf("got %v want ErrNestedBlockAtDepth1", err)
	}
}
