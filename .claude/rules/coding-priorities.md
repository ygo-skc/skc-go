---
paths:
  - "**/*.go"
---

# Coding priorities (read before writing or changing any Go)

When adding or modifying code in this repo, optimize in this order: **correctness → idiomatic Go →
performance → low memory footprint**. The last two are first-class goals here, not afterthoughts —
`common/` is a published library that runs inside every consumer's hot path, and `ygo-service` is a
gRPC service fronting MySQL whose server is hand-tuned down to window sizes and whose binary ships
with `GOEXPERIMENT=greenteagc`.

**Never hand-edit generated code.** `common/ygo/*.pb.go`, `common/health/*.pb.go`, and
`common/swift-gen/` are protoc output. Change the `.proto` and regenerate with `make -C common all`.

**Idiomatic Go**
- Follow Effective Go / Go Code Review Comments: short names in small scopes, `err != nil` handled
  immediately, no needless getters, accept interfaces & return concrete types, keep interfaces small
  and defined at the consumer (as `CardRepository` and friends are in `ygo-service/db/`).
- Match the surrounding code: package-level dependency vars, not framework DI — repositories are
  wired once at `ygo-service/api/server.go:33-38`.
- Logging is request-scoped. A gRPC handler opens the flow with `util.NewLogger(ctx, "Flow Name")`
  and passes the returned ctx down; everything below pulls it back with `util.RetrieveLogger(ctx)`.
  Bare `slog.` is correct **only** in startup/init paths that have no ctx (`main.go`,
  `db/connect.go`, `api/server.go`).
- Errors keep their layer's type — a bare `error` crossing a boundary is wrong:
  - `ygo-service/db/*` returns `*status.Status`, built by `handleQueryError`/`handleRowParsingError`.
  - `ygo-service/api/*` handlers return `err.Err()` (nil-safe on a nil `*status.Status`).
  - `common/client/*` returns `*model.APIError`, converted by `rpcErrorToAPIError`.
- Prefer the standard library over new dependencies. Shared behaviour belongs in `common/util` —
  add it there rather than duplicating it inside `ygo-service`.
- Use `gofmt` semantics and leave no vet warnings. There is no root `go.mod`, so vet **per module**:
  `cd common && go vet ./...`, `cd ygo-service && go vet ./...`.

**Performance**
- Fan out independent downstream/DB work concurrently with `util.AtomicWaitGroup[T]` or
  `util.WorkerPool`; keep genuinely dependent/fail-fast calls sequential. Respect their contracts —
  `AtomicWaitGroup.Store` must be called exactly once or `Load` blocks forever
  (`common/util/atomic.go:13`), and `WorkerPool.Run` closes its own channel and waits.
- Batch instead of looping per item: `common/client` exposes `GetCardsByIDProto`/`GetCardsByNameProto`,
  and the repo layer takes `IN (%s)` placeholder runs. A loop issuing one call per ID is a defect.
- Compute expensive things once, at package or init scope — `quoteRegex`
  (`ygo-service/db/sql_helpers.go:20`), `chicagoLocation` (`ygo-service/api/server.go:24-31`), the
  tuned `sql.DB` pool (`ygo-service/db/connect.go`). Never compile a regex or load a location per call.
- Hoist invariants out of loops, avoid redundant round-trips, and don't re-parse or re-fetch data
  already in hand.
- Let MySQL do the filtering. Project the explicit `cardAttributes` column list rather than
  `SELECT *`, and pass values as `?` placeholders via `variablePlaceholders`/`buildVariableQuerySubjects`
  — never `fmt.Sprintf` a request-derived value into SQL text.

**Low memory footprint**
- Preallocate slices and maps with a known capacity (`make([]T, 0, n)` / `make(map[K]V, n)`) when the
  size is predictable from the input — but not when it isn't, as with a row count before a scan.
- Declare row-scan variables **above** the `for rows.Next()` loop so one set is reused for every row,
  the way `parseCardRows` does (`ygo-service/db/card_repo.go:161-165`).
- Copying is mostly already cheap here: `common/model` is interface- and pointer-based (`YGOCard`
  and `YGOProduct` are interfaces) and the repo layer passes `*ygo.Card`. Don't invent copy-avoidance
  findings. Do range by index if a new **value-type** struct lands in a slice or map on a hot path —
  `[]ProductContent` (`common/model/product.go:25`) is the only one today.
- Don't hold whole result sets longer than needed, and remember anything crossing gRPC must fit
  `MaxSendMsgSize(2<<20)` (`ygo-service/api/server.go:97`).
- Prefer streaming/in-place transforms over building throwaway intermediate slices — see
  `common/parser/text.go`, which scans runes by index with an early exit.

**Module-specific**
- **`common/`** publishes to `proxy.golang.org`, where versions are **immutable**. A changed or
  removed exported symbol forces a major bump and a `/vN` path change in every importer — see the
  `releasing-a-module` skill before touching exported API. An exported helper with no in-repo caller
  is **not** dead code; downstream repos use it.
- **`ygo-service/`** does not own its schema — the SKC MySQL DDL lives in another repo. You cannot
  verify an index from here, so keep queries index-friendly: no function wrapping an indexed column,
  no leading-wildcard `LIKE`, and treat `FORCE INDEX` hints as coupling that breaks at runtime if
  the schema owner renames an index.

Don't add speculative micro-optimizations that hurt readability, and don't add defensive guards for
failure modes that can't occur in this environment — keep changes idiomatic and measured.
