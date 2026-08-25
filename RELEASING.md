# Releasing

Tag `vX.Y.Z` (or push to `main` under auto-release) — `.github/workflows/release.yml` will:

1. Build Bun standalone binaries `praxicraft-assess-{os}-{arch}` + `SHA256SUMS.txt` and attach them to the GitHub Release (`install.sh`).
2. Build `dist/index.js` and **publish `@praxicraft/assess-cli` to npm** (`npm publish --provenance --access public`).

## Prerequisites (GitHub repo secrets)

| Secret | Purpose |
|--------|---------|
| `NPM_TOKEN` | npm automation token with publish rights on `@praxicraft` |
| `GITHUB_TOKEN` | Provided by Actions (release assets) |

Trusted Publishing / provenance needs `id-token: write` (already set on the release workflow).

## Manual release

1. Bump `version` in `package.json` (optional — the workflow also sets version from the tag).
2. Tag and push:

```bash
git tag -a v2.0.1 -m "v2.0.1"
git push origin v2.0.1
```

Or use **Actions → release → Run workflow** with a tag input.

## Local checks

```bash
bun install
bun test
bun run build          # npm package entry (dist/index.js)
bun run build-native   # GitHub Release binaries
```

npm consumers:

```bash
npm install -g @praxicraft/assess-cli
# requires Bun on PATH for the interactive TUI (shebang: bun)
```

## Auto-bump

Pushes to `main` that change package source auto-bump the patch version and `CHANGELOG.md` in the runner, create a local `chore(release)` commit, and **push only the annotated tag** (not `main` — branch protection requires PRs). The tag retains the release commit. CI then **dispatches `release.yml`** (tag pushes from `GITHUB_TOKEN` do not trigger workflows) and opens a **sync PR** so `package.json` and `CHANGELOG.md` on `main` match the release.

Skip with `[skip release]` in the commit message.
