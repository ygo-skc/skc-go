---
name: releasing-a-module
description: Use when cutting a release for a module in skc-go — creating a version tag, pushing it, or publishing GitHub release notes for common or ygo-service.
---

# Releasing a Module

## Overview

skc-go is a Go multi-module monorepo. Each module has an independent version stream, so the tag
prefix, the version number, and the release notes must all describe **one** module. The prefix is
not a naming convention — the Go module proxy requires it to resolve
`github.com/ygo-skc/skc-go/common` to the `common/` subdirectory.

## Modules

| Directory | Tag prefix | Notes |
|---|---|---|
| `common/` | `common/vX.Y.Z` | Library. Consumed through the proxy by other repos. |
| `ygo-service/` | `ygo-service/vX.Y.Z` | Service. Reports its version from `api/status_handler.go`. |

## Choosing the version

Scope to the module being released — both the tag lookup and the diff.

```bash
MOD=common   # or ygo-service
PREV=$(git tag --list "$MOD/*" --sort=-v:refname | head -1)
git log --oneline "$PREV"..HEAD -- "$MOD/"
git diff --stat "$PREV"..HEAD -- "$MOD/"
```

For `ygo-service`, check the version string in `api/status_handler.go` **first**. If an earlier
commit already bumped it, that is the intended version — match it and skip the table. A `.go` diff
limited to that string is still a patch, not an API change.

| Bump | When |
|---|---|
| Patch | Dependency, docs, or config changes only — no `.go` files touched |
| Minor | New exported API, additive and backward-compatible |
| Major | Breaking change — also requires a `/vN` suffix in the module path and in every importer |

## Release notes

The notes are exactly these parts, in this order:

1. `## Changes`
2. One bullet per user-visible change in that module
3. A blank line, then
   `**Full Changelog**: https://github.com/ygo-skc/skc-go/compare/<PREV>...<NEW>`

The release title is the tag name, verbatim — `common/v3.2.1`.

## Sequence

Show the version, the scoped diff, and the drafted notes, and get approval **once**. Then run all
three steps without stopping again:

```bash
git tag "$MOD/vX.Y.Z" <commit>        # lightweight: no -a, no -m
git push origin "$MOD/vX.Y.Z"
gh release create "$MOD/vX.Y.Z" --repo ygo-skc/skc-go \
  --title "$MOD/vX.Y.Z" --notes-file notes.md
```

Push the tag first. `gh release create` attaches to an existing tag, but invents one from the
default branch when the tag is missing.

## Why approval comes before the push

Pushing the tag publishes the module to `proxy.golang.org`, which caches versions immutably. A
wrong version cannot be unpublished — only superseded, or retracted in `go.mod`. Approval is the
last reversible moment.

## Common mistakes

- **`git tag -a`.** Every tag in this repo is lightweight. An annotated tag carries a message
  nobody reads; the notes belong in the GitHub Release.
- **Unscoped tag lookup.** `git tag --sort=-v:refname | head -1` returns the newest tag in the
  repo, which is usually the *other* module's. Always filter by prefix.
- **Unscoped diff.** Dropping `-- "$MOD/"` pulls the other module's commits into the notes.
- **Latest badge.** A new release takes the repo's **Latest** badge from whichever module held it.
  Restore it with `gh release edit <tag> --latest` when that matters.
