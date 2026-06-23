# ecfg

Typed configuration from environment variables for Go.

## Public API

### `ecfg` (primary)

- `ecfg.Parse[T](opts...)` — load config from ENV
- `ecfg.WithPrefix(prefix)` — optional env key prefix
- `ecfg.Usage`, `ecfg.Validator` — implement on Go leaf types
- `ecfg.Err*` — sentinel errors for `errors.Is`

### `pkg/walk` (optional)

Generic struct traversal (reflect or `go/types`). No ecfg policy. Use when you need walking only; see [pkg/walk/doc.go](pkg/walk/doc.go).

Exported: `Engine`, `StructWalk`, `Options`, `VisitCtx`, `FieldDesc`, `NewEngineReflect`, `NewEngineTypes`, `EngineReflect`, `SkipDescend`, `ReflectKind`.

### Internal

`internal/ecfgtool`, `internal/gen` — not stable API.

## Runtime

```go
cfg, err := ecfg.Parse[AppConfig](ecfg.WithPrefix("MYAPP"))
```

- Root struct fields need an `ecfg:"SEGMENT"` tag.
- For **AppResources**: field is the **resource** ([BuildConfiger] or [NewResourceer]); env shape comes from `resource.BuildConfig()` spec ([item.Spec], [app.Spec], …). ecfg does not walk wire/resource fields.
- For **direct config** (`ecfg.Parse[T]`, testdata): the root struct is the config and is walked as-is.
- Nested blocks are one level deep (no nested struct blocks at depth 1).
- Leaves are either a Go named type with `Usage()` and `Validate()`, or a proto wrapper with a single `value` field and `options.v1.usage`.
- Every leaf must be set in the environment (non-empty).

## Codegen

```bash
go generate ./...
```

```bash
go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppConfig -pkg ./config -prefix MYAPP
```

Writes `.env.template` (`KEY=` only) and optional `env.md` (usage tables), grouped by root ecfg block.

## Benchmarks

Fixtures in `internal/testdata`: root configs with **5 / 10 / 15** blocks; each block has **5** leaf fields at depth 1.

| Benchmark | What it measures |
|-----------|------------------|
| `BenchmarkParse_root5/10/15` | Runtime: `ecfg.Parse` (25 / 50 / 75 env vars) |
| `BenchmarkCollectTemplate_root5/10/15` | Offline: `CollectTemplateEntries` (codegen path) |

Запуск через Docker (как `task test`):

```bash
task bench
```

## Mutation testing

[mutest](https://github.com/fchimpan/mutest) mutates relational operators (`==`, `!=`, `<`, …) and runs tests with an overlay. Scope: `internal/ecfgtool`, `pkg/walk`, and the root `ecfg` package (not `cmd/`, `example/`, or heavy `packages.Load` paths).

Запуск через Docker (как `task test`):

```bash
task mutest
```

Minimum mutation score: **65%** (`-threshold 65`). After targeted tests, typical score is ~**88%**; remaining survived mutants are often `usageresolve`/`engine_types` edge paths or equivalent `sort.Slice` comparisons.
