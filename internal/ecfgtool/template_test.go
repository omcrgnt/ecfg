package ecfgtool

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectTemplateEntries_andWrite(t *testing.T) {
	entries, err := CollectTemplateEntries(
		"github.com/omcrgnt/ecfg/internal/testdata",
		"AppConfig",
		"TEST",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries: %+v", entries)
	}
	if entries[0].EnvKey != "TEST_SERVER_LABEL" {
		t.Fatalf("env key: %q", entries[0].EnvKey)
	}
	if entries[0].Usage == "" {
		t.Fatal("expected usage text")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "env.template")
	if err := WriteEnvTemplate(path, entries); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(body)
	if !strings.Contains(s, "TEST_SERVER_LABEL=") {
		t.Fatalf("body: %s", s)
	}
	if strings.Contains(s, "# ") {
		t.Fatalf("env file must not contain usage comments: %s", s)
	}
}

func TestWriteEnvMarkdown(t *testing.T) {
	entries, err := CollectTemplateEntries(
		"github.com/omcrgnt/ecfg/internal/testdata",
		"MultiBlock",
		"M",
	)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "env.md")
	if err := WriteEnvMarkdown(path, "M", entries); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(body)
	if !strings.Contains(s, "# Environment configuration") {
		t.Fatalf("body: %s", s)
	}
	if !strings.Contains(s, "| `M_SERVER_LABEL` |") {
		t.Fatalf("missing table row: %s", s)
	}
	if !strings.Contains(s, "## SERVER") {
		t.Fatalf("missing group heading: %s", s)
	}
}

func TestCollectTemplateEntries_protoCfg(t *testing.T) {
	entries, err := CollectTemplateEntries(
		"github.com/omcrgnt/ecfg/internal/testdata",
		"ProtoCfg",
		"T",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].EnvKey != "T_SERVER_PORT" {
		t.Fatalf("entries: %+v", entries)
	}
}

func TestCollectTemplateEntries_multiBlockBlankLine(t *testing.T) {
	entries, err := CollectTemplateEntries(
		"github.com/omcrgnt/ecfg/internal/testdata",
		"MultiBlock",
		"M",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries: %+v", entries)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "env.template")
	if err := WriteEnvTemplate(path, entries); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "\n\n") {
		t.Fatalf("expected blank line between groups: %q", body)
	}
}

func TestCollectTemplateEntries_typeNotFound(t *testing.T) {
	_, err := CollectTemplateEntries("github.com/omcrgnt/ecfg/internal/testdata", "NoSuchType", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCollectTemplateEntries_packageNotFound(t *testing.T) {
	_, err := CollectTemplateEntries("github.com/no/such/package/xyz", "AppConfig", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCollectTemplateEntries_protoMissingUsage(t *testing.T) {
	_, err := CollectTemplateEntries(
		"github.com/omcrgnt/ecfg/internal/testdata",
		"ProtoNoUsageCfg",
		"T",
	)
	if !errors.Is(err, ErrMissingUsage) {
		t.Fatalf("got %v want ErrMissingUsage", err)
	}
}

func TestWriteEnvTemplate_sortWithinGroup(t *testing.T) {
	entries := []TemplateEntry{
		{EnvKey: "P_GROUP_Z", Usage: "z", RootGroup: "GROUP"},
		{EnvKey: "P_GROUP_A", Usage: "a", RootGroup: "GROUP"},
		{EnvKey: "P_OTHER_B", Usage: "b", RootGroup: "OTHER"},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "env.template")
	if err := WriteEnvTemplate(path, entries); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(body)
	if strings.HasPrefix(s, "\n") {
		t.Fatalf("file must not start with blank line: %q", s)
	}
	zi := strings.Index(s, "P_GROUP_Z")
	ai := strings.Index(s, "P_GROUP_A")
	if zi < 0 || ai < 0 {
		t.Fatalf("missing keys in %q", s)
	}
	if ai > zi {
		t.Fatalf("expected A before Z within group: %q", s)
	}
}

func TestCollectTemplateEntries_sortGroupOrder(t *testing.T) {
	entries, err := CollectTemplateEntries(
		"github.com/omcrgnt/ecfg/internal/testdata",
		"SortGroupCfg",
		"P",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries: %+v", entries)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "env.template")
	if err := WriteEnvTemplate(path, entries); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(body)
	ai := strings.Index(s, "P_GROUP_A")
	zi := strings.Index(s, "P_GROUP_Z")
	if ai < 0 || zi < 0 || ai > zi {
		t.Fatalf("expected A before Z in output: %q", s)
	}
}

func TestCollectTemplateEntries_noUsageLeaf(t *testing.T) {
	_, err := CollectTemplateEntries(
		"github.com/omcrgnt/ecfg/internal/testdata",
		"NoUsageCfg",
		"T",
	)
	if !errors.Is(err, ErrMissingUsage) && !errors.Is(err, ErrIncompatibleLeaf) {
		t.Fatalf("got %v", err)
	}
}
