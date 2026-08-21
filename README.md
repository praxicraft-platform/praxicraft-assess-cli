# Praxicraft Assess CLI

Official command-line interface for the **[Praxicraft Assess](https://assess.praxicraft.com)** Public API.

Binary name: **`praxicraft-assess`**

```bash
# Scripted (AWS-style)
praxicraft-assess org get
praxicraft-assess assessments list --output table
praxicraft-assess invites create senior-backend-screen --email candidate@example.com

# Interactive shell (branded prompt — not a MotD tips feed)
praxicraft-assess
```

**Requires** an organisation API key (`ct_live_…` or `ct_test_…`). Docs: [docs.praxicraft.com](https://docs.praxicraft.com)

## Install

### From source

```bash
git clone https://github.com/praxicraft-platform/praxicraft-assess-cli.git
cd praxicraft-assess-cli
go install ./cmd/praxicraft-assess
```

### GitHub Releases

Download the binary for your OS from [Releases](https://github.com/praxicraft-platform/praxicraft-assess-cli/releases).

## Configure

```bash
praxicraft-assess configure
# or non-interactive:
praxicraft-assess configure --name default --api-key ct_test_… --base-url https://assess.praxicraft.com
```

Config file: `~/.config/praxicraft/config.toml`  
Env: `PRAXICRAFT_API_KEY`, `PRAXICRAFT_API_BASE_URL`, `PRAXICRAFT_PROFILE`

## Command groups

| Group | Examples |
|-------|----------|
| `org` | `get`, `stats`, `team`, `audit-log`, `squads` |
| `assessments` | `list`, `get`, `create`, `update`, `duplicate`, `cases`, `results` |
| `invites` | `list`, `create`, `bulk-create`, `remind`, `cancel`, `result` |
| `results` | `list`, `get` |
| `cases` | `platform-list`, `list`, `create`, `get`, `update`, `delete` |
| `pipelines` | `list`, `get`, `enroll`, `bulk-enroll`, `enrollments`, `reject`, `hold`, `unhold` |
| `webhooks` | `list`, `create`, `get`, `update`, `delete`, `test`, `deliveries`, `retry-delivery` |
| `interviews` | `list`, `create`, `bulk-create`, `get`, `cancel`, `reschedule`, `analysis`, `replay`, `share`, `templates` |
| `integrations` | `list`, `connect-url`, `test` |

Global flags: `--profile`, `--api-key`, `--base-url`, `--output json|table|yaml`, `--query` (JMESPath), `--yes`, `--non-interactive`, `--no-banner`

JSON bodies for create/update: `--body '{"title":"…"}'` or `--body-file path.json`

## Branding

On a TTY, the CLI shows a **Praxicraft Assess** banner (name, version, docs link) and a branded `praxicraft-assess>` REPL prompt. This is product identity — **not** a MotD tips feed. Disable with `--no-banner` or when stdout is not a TTY.

## Examples

See [praxicraft-assess-examples](https://github.com/praxicraft-platform/praxicraft-assess-examples).

## License

MIT
