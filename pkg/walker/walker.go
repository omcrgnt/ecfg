package walker

import (
	"fmt"
	"reflect"
	"strconv"
)

// NodeKind определяет тип узла в иерархии (структура, коллекция или значение)
type NodeKind int

const (
	KindStruct NodeKind = iota
	KindSlice
	KindMap
	KindLeaf // Конечное значение (int, string, time.Duration и т.д.)
)

// NodeInfo несет информацию об узле для хуков OnEnter/OnExit
type NodeInfo struct {
	Name string            // Имя поля, индекс слайса "0" или ключ мапы
	Kind NodeKind          // Тип узла
	Tag  reflect.StructTag // Теги поля (если узел является полем структуры)
}

// FieldContext содержит данные для финального колбэка обработки значения
type FieldContext struct {
	Value reflect.Value       // Значение, в которое можно писать
	Field reflect.StructField // Метаданные поля (имя, теги, тип)
}

// WalkFunc — колбэк для обработки конечных значений ("листьев")
type WalkFunc func(ctx FieldContext) error

// Walker — основной объект обхода
type Walker struct {
	initNilPointers bool
	onEnter         func(NodeInfo)
	onExit          func(NodeInfo)
}

type Option func(*Walker)

// WithInitNilPointers включает автоматическую аллокацию памяти для nil-указателей
func WithInitNilPointers() Option {
	return func(w *Walker) {
		w.initNilPointers = true
	}
}

// Опции для настройки хуков
func WithOnEnter(fn func(NodeInfo)) Option {
	return func(w *Walker) { w.onEnter = fn }
}

func WithOnExit(fn func(NodeInfo)) Option {
	return func(w *Walker) { w.onExit = fn }
}

func New(opts ...Option) *Walker {
	w := &Walker{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Process — Generic-хелпер для создания и обхода нового экземпляра типа T
func Process[T any](w *Walker, fn WalkFunc) (*T, error) {
	var target T
	val := reflect.ValueOf(&target).Elem()
	if err := w.Walk(val, fn); err != nil {
		return nil, err
	}
	return &target, nil
}

// Walk запускает рекурсивный обход переданного значения
func (w *Walker) Walk(v reflect.Value, fn WalkFunc) error {
	return w.recursiveWalk(v, reflect.StructField{}, fn)
}

func (w *Walker) recursiveWalk(v reflect.Value, field reflect.StructField, fn WalkFunc) error {
	// 1. Обработка указателей
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			if !w.initNilPointers {
				return nil
			}
			if v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		}
		return w.recursiveWalk(v.Elem(), field, fn)
	}

	// 2. Обработка контейнеров и листьев
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			fv := v.Field(i)
			if !f.IsExported() {
				continue
			}

			info := NodeInfo{Name: f.Name, Kind: KindStruct, Tag: f.Tag}

			isCont := isContainer(fv)
			if isCont {
				w.enter(info)
			}

			if err := w.recursiveWalk(fv, f, fn); err != nil {
				return err
			}

			if isCont {
				w.exit(info)
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			info := NodeInfo{
				Name: strconv.Itoa(i),
				Kind: KindSlice,
			}

			w.enter(info)

			if err := w.recursiveWalk(v.Index(i), field, fn); err != nil {
				return err
			}

			w.exit(info)
		}

	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			info := NodeInfo{
				Name: fmt.Sprintf("%v", iter.Key().Interface()),
				Kind: KindMap,
			}

			w.enter(info)

			if err := w.recursiveWalk(iter.Value(), field, fn); err != nil {
				return err
			}

			w.exit(info)
		}

	default:
		// Конечное значение
		return fn(FieldContext{
			Value: v,
			Field: field,
		})
	}

	return nil
}

// Вспомогательные методы, которые делают код чище
func (w *Walker) enter(info NodeInfo) {
	if w.onEnter != nil {
		w.onEnter(info)
	}
}

func (w *Walker) exit(info NodeInfo) {
	if w.onExit != nil {
		w.onExit(info)
	}
}

// isContainer проверяет, является ли значение контейнером
func isContainer(v reflect.Value) bool {
	k := v.Kind()
	if k == reflect.Pointer {
		// Если указатель nil, проверяем тип, на который он указывает
		if v.IsNil() {
			return isContainerType(v.Type().Elem().Kind())
		}
		return isContainer(v.Elem())
	}
	return isContainerType(k)
}

// isContainerType проверяет сам Kind
func isContainerType(k reflect.Kind) bool {
	return k == reflect.Struct || k == reflect.Slice || k == reflect.Array || k == reflect.Map
}
