# Changelog


## 3.0.0

### Breaking

- Rename CLI resource from cases to tasks; Public API paths use /tasks/.

## Unreleased

- feat: `--all` on list commands follows cursor pagination and merges rows into one table
- fix: table output shows a pagination hint instead of raw next URLs
- fix: poll npm/GitHub in the background so the update banner appears without restarting the TUI
- chore: sync CHANGELOG and package.json on main for releases 2.0.7–2.0.9
- chore: auto-release opens a sync PR and triggers the release workflow after tagging

## 2.0.9

- fix: visible update banner + brand-only colors (#22)
- fix: single update banner at top of TUI (no duplicate above input or on stderr)

## 2.0.8

- fix: version-check test independent of package.json version (#21)
- fix: align tabular CLI output with proper column padding (#19)
- fix: auto-release skips existing semver tags (#20)
- feat: tabular CLI output by default; whoami shows org name only (#18)

## 2.0.7

- fix: update banner checks GitHub releases and shows at top (#17)
- feat: show update banner in TUI when a newer CLI is on npm (#16)

## 2.0.6

- Ask AI: plain TUI prose (strip Markdown), direct docs URLs — never show llms.txt to users

## 2.0.5

- Fix Ask AI stalling on “I’ll fetch docs…” (nudge + answer SDK catalog without MCP)
- Timeout individual tool executes so docs MCP cannot hang the TUI

## 2.0.4

- Republish of Ask AI speed fix (stale v2.0.3 tag blocked release).

## 2.0.3

- Ask AI: skip cold MCP for product/CLI questions (much faster)
- Soft-fail slow Assess MCP (`npx`); docs MCP only when needed
- Pair with backend prompt so related product questions are not refused

## 2.0.2

- Fix Ask AI in standalone binaries: docs MCP uses HTTP (no mcp-remote / bunfs)
- /help renders a branded command list again

## 2.0.1

- One-line PRAXICRAFT welcome mark (slick/tiny) with Assess CLI label
- Distinct input rail, tips, and footer so the TUI feels on-brand

## 2.0.0

- Rewrite in Bun + TypeScript with OpenTUI interactive shell
- Command palette (`/`), Ask AI (Assess assistant proxy + dual MCP)
- Preserve `~/.config/praxicraft/config.toml` and env var names
- Go sources archived under `legacy/`

## Earlier

See git history and GitHub Releases for pre-2.0.0 notes.
