# ecfg

Typed configuration from environment variables for Go.

## Public API

### `ecfg` (org bootstrap)

- `ecfg.LoadEnv(reg)` — initialize config values in side-registry entries from ENV
- `ecfg.SetPrefix(prefix)` — env key prefix (default `APP`)
- `ecfg.SetTagKey(key)` — registry custom tag key (default `ecfg`)
- `ecfg.TagKey()` — current tag key for `LoadEnv` and `unique.AddWithCustomTag`
- `cmd/ecfg-gen` — codegen from `AppResources` ([BuildConfig] AST)

Pipeline step **Apply** → `ecfg.SetPrefix` + `ecfg.LoadEnv(registrySpecs)`.

### `ecfg/config` (standalone)

- `config.Parse[T](opts...)` — load config from ENV into a new struct
- `config.WithPrefix(prefix)` — per-call env key prefix
- `config.Err*` — sentinel errors for `errors.Is`

No registry, `res`, or `unique` dependency.

### `pkg/walk` (optional)

Generic struct traversal (reflect or `go/types`). No ecfg policy. See [pkg/walk/doc.go](pkg/walk/doc.go).

### Internal

`internal/ecfgtool`, `internal/gen` — not stable API.

## Runtime

**Org (side-registry):**

```go
ecfg.SetPrefix("DEMO")
regSpecs := unique.New()
regSpecs.AddWithCustomTag(cfg, ecfg.TagKey(), "SERVICE_ITEM")
ecfg.LoadEnv(regSpecs)
```

**Standalone:**

```go
cfg, err := config.Parse[AppConfig](config.WithPrefix("MYAPP"))
```

- Registry entries need custom tag `ecfg` (or [SetTagKey]) with segment value (`SERVICE_ITEM`, …).
- Standalone root struct fields need `ecfg:"SEGMENT"` tag.
- Nested blocks are one level deep (no nested struct blocks at depth 1).
- Leaves: Go named type with `Usage()` and `Validate()`, or proto wrapper with `value` + `options.v1.usage`.
- Every leaf must be set in the environment (non-empty).

## Breaking migration (v0.21.0)

| Was | Now |
|-----|-----|
| `ecfg.Parse[T]` | `config.Parse[T]` |
| `ecfg.Apply(reg, &ar, …)` | `ecfg.SetPrefix` + `ecfg.LoadEnv(regSpecs)` |
| `ecfg.Register` | removed |

## Codegen

```bash
go generate ./...
```

```bash
go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppResources -pkg ./cmd/demo -prefix DEMO -template .env.template
```

## Benchmarks

Fixtures in `internal/testdata`: root configs with **5 / 10 / 15** blocks; each block has **5** leaf fields at depth 1.

| Benchmark | What it measures |
|-----------|------------------|
| `BenchmarkParse_root5/10/15` | `config.Parse` (25 / 50 / 75 env vars) |
| `BenchmarkCollectTemplate_root5/10/15` | `CollectTemplateEntries` (codegen path) |

```bash
task bench
```

## Mutation testing

```bash
task mutest
```

Minimum mutation score: **65%** (`-threshold 65`).
