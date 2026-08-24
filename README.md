# Praxicraft Assess CLI

Interactive terminal UI (OpenTUI) for the [Assess Public API](https://docs.praxicraft.com) — command palette, Ask AI, and scriptable subcommands.

Binary: `praxicraft-assess`  
Repo: [praxicraft-platform/praxicraft-assess-cli](https://github.com/praxicraft-platform/praxicraft-assess-cli)

## Install

### npm (requires [Bun](https://bun.sh) ≥ 1.1 for the TUI)

```bash
npm install -g @praxicraft/assess-cli
# or
bun add -g @praxicraft/assess-cli
```

### Binary installer

```bash
curl -fsSL https://praxicraft.com/install.sh | sh
```

Installs the latest release binary (SHA-256 verified) as `praxicraft-assess` to `~/.local/bin` or `/usr/local/bin`.

```bash
PRAXICRAFT_VERSION=v2.0.0 curl -fsSL https://praxicraft.com/install.sh | sh
```

Or download `praxicraft-assess-{os}-{arch}` from [Releases](https://github.com/praxicraft-platform/praxicraft-assess-cli/releases).

### From source (Bun)

```bash
git clone https://github.com/praxicraft-platform/praxicraft-assess-cli.git
cd praxicraft-assess-cli
bun install
bun run src/index.ts
```

Requires [Bun](https://bun.sh) ≥ 1.1.

## Quickstart

```bash
praxicraft-assess configure   # or /login in the TUI
praxicraft-assess             # interactive shell
praxicraft-assess whoami
praxicraft-assess assessments list
```

Profiles: `~/.config/praxicraft/config.toml`  
Env: `PRAXICRAFT_API_KEY`, `PRAXICRAFT_API_BASE_URL`, `PRAXICRAFT_PROFILE`

## Interactive UI

- **Brand welcome** + tip (`/login` when signed out)
- **Ask anything…** bar — free text runs Ask AI
- **Type `/`** for the command palette (↑↓ navigate, ↵ select, esc cancel)
- Status: `signed out` or `live|test · profile · version`

### Ask AI

Uses your Assess API key against `POST /api/v1/public/assistant/chat/` (Starter+, scope `assistant:write`) and local MCP:

- Docs index: `https://docs.praxicraft.com/llms.txt` (included on each Ask AI turn)
- Knowledge MCP: `https://docs.praxicraft.com/mcp`
- API tools: `@praxicraft/assess-mcp` (stdio)

## Commands

| Group | Examples |
|-------|----------|
| Auth | `/login`, `/logout`, `configure`, `whoami` |
| Org | `/org get`, `/org billing`, `/org stats` |
| Assessments | `/assessments list`, `/assessments get <slug>` |
| Invites | `/invites list`, `/invites create …` |
| More | `/results`, `/cases`, `/pipelines`, `/webhooks`, `/interviews`, `/integrations` |
| AI | free text or `/ai <question>` |

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md).

## License

MIT
