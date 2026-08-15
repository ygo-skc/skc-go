# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Go **multi-module monorepo** holding the shared libraries and micro-services behind the SKC
projects. Two modules, independently versioned and released (see the `releasing-a-module` skill):

| Module | Path | What it is |
| --- | --- | --- |
| `common/` | `github.com/ygo-skc/skc-go/common/v3` | Published library — gRPC clients, models, parsers, logging/concurrency utils. Consumed through the Go proxy by `ygo-service` **and other repos** (e.g. skc-suggestion-engine). |
| `ygo-service/` | `github.com/ygo-skc/skc-go/ygo-service` | gRPC service on port **9020** serving card, product, restriction, and score data out of the SKC MySQL database. |

`ygo-service/go.mod` both `require`s `common/v3` and `replace`s it with `../common`, so local builds
always compile against the working tree — a stale `require` version is invisible until something
builds without the replace.

**The MySQL schema is owned by another repo** (skc-api). There is no DDL here, so an index cannot be
verified from this codebase — only queried in an index-friendly shape.

## Coding priorities

Optimize in this order: **correctness → idiomatic Go → performance → low memory footprint**. The
last two are first-class goals here, not afterthoughts. Full checklist (idiomatic Go / performance /
memory / per-module rules): `.claude/rules/coding-priorities.md`, which loads automatically whenever
you touch a `.go` file.

**Generated code is never hand-edited.** `common/ygo/*.pb.go`, `common/health/*.pb.go`, and
`common/swift-gen/` are protoc output — change the `.proto` and regenerate.

## Commands

There is **no root `go.mod`**, so `go` commands run per module — `./...` from the repo root matches
nothing.

| Command | Purpose |
| --- | --- |
| `gofmt -l common ygo-service` | Format check across both modules (empty output = clean) |
| `cd common && go vet ./...` | Vet `common` |
| `cd ygo-service && go vet ./...` | Vet `ygo-service` |
| `cd common && go test ./...` | Tests — only `common/parser` has any today |
| `make -C ygo-service build` | `go mod tidy` + vet + cross-compile Linux ARM64 static binary |
| `make -C common all` | Regenerate protobuf/gRPC (needs `protoc` + plugins; runs `go install`) |
| `cd ygo-service && go run .` | Run the service — needs `certs/` and `YGO_SERVICE_DOT_ENV_FILE` |

A `PostToolUse` hook (`.claude/settings.json`) runs gofmt + vet + test on the changed package after
every edit, resolving the owning module automatically.

CI (`.github/workflows`) runs **CodeQL only** — there is no unit-test workflow. Renovate automerges
`minor`, `patch`, `pin`, and `digest` updates, so both modules sit on the latest patch by construction.

## Architecture

```
ygo-service/main.go → db.EstablishDBConn() + go api.RunService()

ygo-service/api/  gRPC handlers, one file per service. Opens a flow with util.NewLogger(ctx, ...),
                  delegates to a repository, returns err.Err(). TLS creds + hand-tuned server
                  windows/keepalive in server.go.
ygo-service/db/   MySQL DAO. One ...Repository interface + YGO...Repository struct per domain,
                  wired as package vars in api/server.go. Queries are const blocks at file top.
common/client/    gRPC clients for the above, returning *model.APIError.
common/model/     Shared models + builders. Interface-based (YGOCard, YGOProduct are interfaces).
common/util/      Logging (ctx-scoped slog), env, TLS, concurrency (AtomicWaitGroup, WorkerPool).
common/parser/    Card-effect text parsing.
```

Errors keep their layer's type: `*status.Status` in `db/`, `err.Err()` at the handler boundary,
`*model.APIError` in `common/client`. See `.claude/rules/coding-priorities.md`.

## Skills

- `go-drift-check` — periodic health audit of both modules (idiomatic Go, SQL hygiene, performance
  and memory, adoptable dependency/stdlib additions).
- `releasing-a-module` — cutting and publishing a per-module version tag.
