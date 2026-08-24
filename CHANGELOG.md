# Changelog

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

See git history for 1.x Go releases.
