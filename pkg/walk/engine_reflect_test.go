package walk

import (
	"testing"
)

func TestNewEngineReflect_notStruct(t *testing.T) {
	_, err := NewEngineReflect(42)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEngineReflect_fieldValueOnChild(t *testing.T) {
	type block struct {
		Label string
	}
	type root struct {
		Server block
	}
	engRoot, err := NewEngineReflect(&root{})
	if err != nil {
		t.Fatal(err)
	}
	child, err := engRoot.Child(FieldDesc{Name: "Server"})
	if err != nil {
		t.Fatal(err)
	}
	engChild, ok := child.(*EngineReflect)
	if !ok {
		t.Fatal("expected EngineReflect child")
	}
	val, _, err := engChild.FieldValue(FieldDesc{Name: "Label"})
	if err != nil {
		t.Fatal(err)
	}
	if val.String() != "" {
		t.Fatalf("expected empty label, got %q", val.String())
	}
}

func TestInitPointerField_allocatesNilStructPtr(t *testing.T) {
	type block struct {
		Label string
	}
	type cfg struct {
		Server *block
	}
	root := &cfg{}
	eng, err := NewEngineReflect(root)
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.InitPointerField(fields[0]); err != nil {
		t.Fatal(err)
	}
	if root.Server == nil {
		t.Fatal("expected Server to be allocated")
	}
}

func TestInitPointerField_nonStructPtrUnchanged(t *testing.T) {
	type cfg struct {
		N *int
	}
	n := 7
	root := &cfg{N: &n}
	eng, err := NewEngineReflect(root)
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.InitPointerField(fields[0]); err != nil {
		t.Fatal(err)
	}
	if root.N == nil || *root.N != 7 {
		t.Fatalf("N: %v", root.N)
	}
}
