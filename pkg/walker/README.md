# walker

Обход struct-дерева через `Provider` и колбэк `Handler`. Пакет не привязан к env/тегам `ecfg` — его использует `ecfg.Parse` для чтения переменных окружения.

## Модель

```
Walker.Walk(Provider, Handler)
    │
    ├─ RuntimeProvider  → walkValues (reflect, есть Value())
    └─ SchemaProvider   → walkFields  (AST / reflect.Type, Value() → ErrNoRuntimeValue)
```

| Провайдер | Интерфейс | Назначение |
|-----------|-----------|------------|
| `NewReflectProvider(&cfg)` | `RuntimeProvider` | Заполнение struct в runtime |
| `NewTypesProvider(pkg, name)` | `SchemaProvider` | Статическая схема через `go/types` |
| `reflectTypeProvider` (внутренний) | `SchemaProvider` | Схема элемента slice/map из `ElemProvider()` |

`Walk` не делает `switch` по конкретным типам провайдеров: достаточно реализовать `RuntimeProvider` или `SchemaProvider` (маркерные методы `runtimeProvider()` / `schemaProvider()`).

## Slice и map

Два пути:

1. **Schema (AST)** — `ElemProvider()` + `walkFields`: для `[]Struct` обходятся поля элемента без индексов (нет runtime slice).
2. **Runtime (reflect)** — `walkContent`: индексы slice/array и ключи map, `NodeHook` для путей вроде `Items.0`.

Пустой slice/map в runtime не создаёт элементы — env-индексы в `ecfg` не материализуются сами по себе.

## Proto

Поля, реализующие `proto.Message`, считаются **листом**: `Walk` не заходит в `Seconds`, `Nanos` и т.д. См. `TestWalk_protoFieldIsLeaf`.

## Field.Kind()

У `types`-провайдера `Kind()` возвращает точный `reflect.Kind` (`Int32`, `Uint8`, …), а не сводит все int к `reflect.Int`.

На что влияет:

- `IsStruct()` — только `Kind() == reflect.Struct` (вложенный struct vs slice/map).
- `IsProto()` — отдельная проверка по method set / `proto.Message`.
- Обработчик может различать `int` и `int32` при статическом анализе; для reflect-пути `Kind()` берётся из `reflect.Type`.

## Зависимости и подпакеты

Сейчас один импорт: `github.com/omcrgnt/ecfg/pkg/walker`.

- `golang.org/x/tools` нужен только при вызове `NewTypesProvider`.
- `google.golang.org/protobuf` — для `IsProto()` (reflect и types).

Вынесение в `walker/reflect` и `walker/types` **не уменьшит** размер бинарника, если вы всё равно вызываете оба API: линкер тянет только используемые пакеты. Подпакеты имеют смысл, если хотите явно разделить импорты в коде (`import reflectprov "…/walker/reflect"`) — это отдельный рефакторинг, не обязателен для ecfg.

## Пример

```go
w := walker.New(walker.WithInitNilPointers())
p, err := walker.NewReflectProvider(&cfg)
if err != nil { return err }
return w.Walk(p, func(f walker.Field) error {
    rv, sf, err := f.Value()
    if err != nil { return err }
    _ = rv
    _ = sf.Tag.Get("ecfg")
    return nil
})
```

Статическая схема:

```go
p, err := walker.NewTypesProvider("my/app/config", "Config")
w := walker.New()
return w.Walk(p, func(f walker.Field) error {
  // f.Name(), f.Tag("ecfg"), f.Kind() — без Value()
  return nil
})
```
