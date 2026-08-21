# Releasing

1. Bump `Version` default in `internal/cmdroot/root.go` / `CHANGELOG.md` if needed.
2. Tag `v0.1.0` (or next) and push the tag.
3. GitHub Actions `release.yml` runs GoReleaser and attaches binaries to the GitHub Release.

Manual:

```bash
goreleaser release --clean
```
