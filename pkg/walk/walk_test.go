package walk

import (
	"errors"
	"testing"
)

type inner struct {
	X int
}

type root struct {
	Nested inner
}

type rootWithValueStruct struct {
	Block inner
}

func TestStructWalk_skipDescend(t *testing.T) {
	eng, err := NewEngineReflect(&root{})
	if err != nil {
		t.Fatal(err)
	}
	var visited []string
	err = StructWalk(eng, Options{}, func(ctx VisitCtx) error {
		visited = append(visited, ctx.Field.Name)
		if ctx.Field.Name == "Nested" {
			return SkipDescend()
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(visited) != 1 || visited[0] != "Nested" {
		t.Fatalf("visited: %v", visited)
	}
}

func TestIsStructField_reflectAndTypes(t *testing.T) {
	eng, err := NewEngineTypes("github.com/omcrgnt/ecfg/internal/testdata", "AppConfig")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if !isStructField(fields[0]) {
		t.Fatal("Server block should be struct field")
	}
	child, err := eng.Child(fields[0])
	if err != nil {
		t.Fatal(err)
	}
	leafFields, err := child.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if isStructField(leafFields[0]) {
		t.Fatal("Label leaf should not be struct field")
	}
	reflectEng, err := NewEngineReflect(&root{})
	if err != nil {
		t.Fatal(err)
	}
	rf, err := reflectEng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if !isStructField(rf[0]) {
		t.Fatal("Nested should be struct via reflect")
	}
}

func TestStructWalk_descendsValueStruct(t *testing.T) {
	eng, err := NewEngineReflect(&rootWithValueStruct{Block: inner{X: 2}})
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
	want := []string{"Block", "X"}
	if len(names) != len(want) {
		t.Fatalf("visited %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("visited %v want %v", names, want)
		}
	}
}

func TestIsStructField_valueStructKind(t *testing.T) {
	eng, err := NewEngineReflect(&rootWithValueStruct{Block: inner{X: 1}})
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if !isStructField(fields[0]) {
		t.Fatal("embedded struct value field should be struct kind")
	}
}

func TestStructWalk_visitError(t *testing.T) {
	eng, err := NewEngineReflect(&root{})
	if err != nil {
		t.Fatal(err)
	}
	want := errors.New("visit failed")
	err = StructWalk(eng, Options{}, func(VisitCtx) error { return want })
	if !errors.Is(err, want) {
		t.Fatalf("got %v", err)
	}
}
