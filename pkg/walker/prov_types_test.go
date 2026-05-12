package walker

import (
	"reflect"
	"testing"
)

func TestTypesProvider(t *testing.T) {
	// Путь к пакету внутри твоего проекта
	pkgPath := "github.com/omcrgnt/ecfg/internal/teststruct"

	// Создаем провайдер без скачивания (предполагаем, что код на диске есть)
	p, err := NewTypesProvider(pkgPath, "Config")
	if err != nil {
		t.Fatalf("failed to load types provider: %v", err)
	}

	fields, err := p.GetFields()
	if err != nil {
		t.Fatalf("failed to get fields: %v", err)
	}

	// Сверяем данные с тем, что мы знаем о структуре
	expected := map[string]struct {
		tag   string
		usage string
		kind  reflect.Kind
	}{
		"Port": {"PORT", "Server port", reflect.Int},
		"User": {"USER", "User name", reflect.String},
	}

	for _, f := range fields {
		exp, ok := expected[f.Name()]
		if !ok {
			t.Errorf("unexpected field found: %s", f.Name())
			continue
		}

		if f.Tag("ecfg") != exp.tag || f.Tag("usage") != exp.usage || f.Kind() != exp.kind {
			t.Errorf("field %s: expected %+v, got name=%s, tag=%s, usage=%s, kind=%v",
				f.Name(), exp, f.Name(), f.Tag("ecfg"), f.Tag("usage"), f.Kind())
		}
	}
}
