package walk

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestNewEngineTypes_ok(t *testing.T) {
	eng, err := NewEngineTypes("github.com/omcrgnt/ecfg/internal/testdata", "AppConfig")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || fields[0].Name != "Server" {
		t.Fatalf("fields: %+v", fields)
	}
	if fields[0].ReflectType != nil {
		t.Fatal("expected types-only FieldDesc")
	}
	if fields[0].TypesType == nil {
		t.Fatal("expected TypesType")
	}
	child, err := eng.Child(fields[0])
	if err != nil {
		t.Fatal(err)
	}
	childFields, err := child.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if len(childFields) != 1 || childFields[0].Name != "Label" {
		t.Fatalf("child fields: %+v", childFields)
	}
}

func TestNewEngineTypes_typeNotFound(t *testing.T) {
	_, err := NewEngineTypes("github.com/omcrgnt/ecfg/internal/testdata", "NoSuchType")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("got %v", err)
	}
}

func TestNewEngineTypes_packageNotFound(t *testing.T) {
	_, err := NewEngineTypes("github.com/no/such/package/xyz", "AppConfig")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReflectKind_types(t *testing.T) {
	eng, err := NewEngineTypes("github.com/omcrgnt/ecfg/internal/testdata", "AppConfig")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if k := ReflectKind(fields[0].TypesType); k != reflect.Struct {
		t.Fatalf("block kind: %v want Struct", k)
	}
	child, err := eng.Child(fields[0])
	if err != nil {
		t.Fatal(err)
	}
	leafFields, err := child.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if k := ReflectKind(leafFields[0].TypesType); k != reflect.String {
		t.Fatalf("leaf kind: %v want string", k)
	}
}

func TestStructWalk_typesEngine(t *testing.T) {
	eng, err := NewEngineTypes("github.com/omcrgnt/ecfg/internal/testdata", "MultiBlock")
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	err = StructWalk(eng, Options{}, func(ctx VisitCtx) error {
		names = append(names, ctx.Field.Name)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Server", "Label", "Worker", "Label"}
	if len(names) != len(want) {
		t.Fatalf("visited %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("visited %v want %v", names, want)
		}
	}
}

func TestStructWalk_fieldsError(t *testing.T) {
	err := StructWalk(errFieldsEngine{}, Options{}, func(VisitCtx) error { return nil })
	if err == nil || err.Error() != "boom" {
		t.Fatalf("got %v", err)
	}
}

func TestStructWalk_childError(t *testing.T) {
	eng := childErrEngine{fields: []FieldDesc{{
		Name:        "Nested",
		ReflectType: reflect.TypeOf(struct{}{}),
	}}}
	err := StructWalk(eng, Options{}, func(VisitCtx) error { return nil })
	if err == nil || err.Error() != "child err" {
		t.Fatalf("got %v", err)
	}
}

type errFieldsEngine struct{}

func (errFieldsEngine) Fields() ([]FieldDesc, error) { return nil, errors.New("boom") }
func (errFieldsEngine) Child(FieldDesc) (Engine, error) { return nil, nil }

type childErrEngine struct {
	fields []FieldDesc
}

func (e childErrEngine) Fields() ([]FieldDesc, error) { return e.fields, nil }
func (e childErrEngine) Child(FieldDesc) (Engine, error) { return nil, errors.New("child err") }
