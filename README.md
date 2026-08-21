# Praxicraft Assess CLI

Official command-line tool for the [Praxicraft Assess](https://assess.praxicraft.com) Public API.

Invite candidates, manage assessments and pipelines, work with webhooks and interviews, and automate hiring workflows from your terminal.

```bash
praxicraft-assess org get
praxicraft-assess assessments list --output table
praxicraft-assess invites create senior-backend-screen --email candidate@example.com
```

Run `praxicraft-assess` with no arguments to open an interactive shell.

You need an organisation API key (`ct_live_…` or `ct_test_…`) from **Assess → Developer → API Keys**. Full API docs: [docs.praxicraft.com](https://docs.praxicraft.com).

## Install

### From GitHub Releases

Download a binary for your OS from [Releases](https://github.com/praxicraft-platform/praxicraft-assess-cli/releases), make it executable, and put it on your `PATH`.

### From source

```bash
git clone https://github.com/praxicraft-platform/praxicraft-assess-cli.git
cd praxicraft-assess-cli
go install ./cmd/praxicraft-assess
```

Requires Go 1.22+.

## Configure

```bash
praxicraft-assess configure
```

Or set credentials without the wizard:

```bash
export PRAXICRAFT_API_KEY="ct_test_xxxxxxxxxxxxxxxx"
# optional
export PRAXICRAFT_API_BASE_URL="https://assess.praxicraft.com"
```

Profiles are stored in `~/.config/praxicraft/config.toml`. Use `--profile` to switch between them.

Check that you’re authenticated:

```bash
praxicraft-assess whoami
```

## Common commands

```bash
# Organisation
praxicraft-assess org get
praxicraft-assess org stats

# Assessments
praxicraft-assess assessments list
praxicraft-assess assessments get senior-backend-screen

# Invites
praxicraft-assess invites create senior-backend-screen --email candidate@example.com --name "Jane Doe"
praxicraft-assess invites list
praxicraft-assess invites result <invite-token>

# Results
praxicraft-assess results list senior-backend-screen
praxicraft-assess results get <invite-token>

# Webhooks
praxicraft-assess webhooks list
praxicraft-assess webhooks create --body '{"url":"https://example.com/hooks","events":["candidate.passed"]}'
```

Create and update payloads can be passed as JSON:

```bash
praxicraft-assess assessments create --body '{"title":"Backend screen"}'
praxicraft-assess assessments create --body-file ./assessment.json
```

## What you can manage

| Command group | Purpose |
|---------------|---------|
| `org` | Profile, stats, team, audit log, squads |
| `assessments` | List, create, update, duplicate, cases, results |
| `invites` | Create, bulk create, list, remind, cancel, result |
| `results` | List by assessment, get by invite token |
| `cases` | Platform and organisation cases |
| `pipelines` | Enroll, bulk enroll, hold, reject, enrollments |
| `webhooks` | Endpoints, test pings, deliveries, retries |
| `interviews` | Rooms, templates, analysis, share, cancel |
| `integrations` | Connected ATS providers |

Run `praxicraft-assess <command> --help` for flags on any command.

## Useful flags

| Flag | Description |
|------|-------------|
| `--api-key` | API key (overrides env / config) |
| `--base-url` | API host (default `https://assess.praxicraft.com`) |
| `--profile` | Named config profile |
| `--output` | `json` (default for scripts), `table`, or `yaml` |
| `--query` | JMESPath filter on JSON output |
| `--filter` | API list filter as `key=value` (repeatable; on list commands) |
| `--yes` | Skip confirmation prompts |
| `--non-interactive` | Disable interactive prompts (CI-friendly) |

## Examples

More recipes (SDKs, curl, n8n, Zapier) live in [praxicraft-assess-examples](https://github.com/praxicraft-platform/praxicraft-assess-examples).

## License

MIT
