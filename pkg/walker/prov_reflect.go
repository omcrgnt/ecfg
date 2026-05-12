package walker

import "reflect"

type reflectProvider struct {
	t reflect.Type
}

func NewReflectProvider(v interface{}) (Provider, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return &reflectProvider{t: t}, nil
}

func (p *reflectProvider) GetFields() ([]Field, error) {
	var fields []Field
	for i := 0; i < p.t.NumField(); i++ {
		f := p.t.Field(i)
		if !f.IsExported() {
			continue
		}
		fields = append(fields, &reflectField{f: f})
	}
	return fields, nil
}

func (p *reflectProvider) EntryName() string { return p.t.Name() }

type reflectField struct{ f reflect.StructField }

func (f *reflectField) Name() string          { return f.f.Name }
func (f *reflectField) Tag(key string) string { return f.f.Tag.Get(key) }
func (f *reflectField) Kind() reflect.Kind {
	t := f.f.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind()
}
func (f *reflectField) IsStruct() bool { return f.Kind() == reflect.Struct }
func (f *reflectField) IsProto() bool {
	t := f.f.Type
	_, ok := reflect.PtrTo(t).MethodByName("ProtoMessage")
	return ok
}
func (f *reflectField) GetProvider() (Provider, error) {
	return NewReflectProvider(reflect.New(f.f.Type).Elem().Interface())
}
