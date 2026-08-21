# Releasing

## Automatic (preferred)

Pushes to `main` that change CLI source trigger **auto-release**:

- Paths: `cmd/**`, `internal/**`, `go.mod`, `go.sum`, `.goreleaser.yaml`
- Bumps the patch version (`v0.1.0` → `v0.1.1`)
- Prepends notes to `CHANGELOG.md` from commits since the last tag
- Updates default `Version` in `internal/cmdroot/root.go` and `cmd/praxicraft-assess/main.go`
- Commits `chore(release): vX.Y.Z`, creates annotated tag `vX.Y.Z`, pushes
- Runs GoReleaser and attaches binaries to the GitHub Release

Skipped when the head commit message starts with `chore(release):` (avoids loops).

To skip a release for a source change, put `[skip release]` in the commit message.

## Manual

1. Update `CHANGELOG.md` and default `Version` strings if needed.
2. Tag and push:

```bash
git tag -a v0.1.2 -m "Release v0.1.2"
git push origin v0.1.2
```

3. Workflow `release.yml` runs GoReleaser on `v*` tags (also `workflow_dispatch`).

```bash
goreleaser release --clean
```
