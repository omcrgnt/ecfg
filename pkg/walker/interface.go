package walker

import "reflect"

// Field абстрагирует поле структуры для обоих миров
type Field interface {
	Name() string
	Tag(key string) string
	IsStruct() bool
	// IsProto возвращает true, если тип поля — протобаф-сообщение
	IsProto() bool
	// ParentValue() если нужен доступ к реальным данным (только для рефлексии)
	Kind() reflect.Kind
	GetProvider() (Provider, error)
}

// Provider — это то, что должен реализовать reflect.go и ast.go
type Provider interface {
	GetFields() ([]Field, error)
	// EntryName возвращает имя самой структуры (Config, и т.д.)
	EntryName() string
}

// Handler — функция-коллбэк, которая будет одинаковой везде
type Handler func(f Field) error
