# Praxicraft Assess CLI

Official command-line tool for the [Praxicraft Assess](https://assess.praxicraft.com) Public API.

Invite candidates, manage assessments and pipelines, work with webhooks and interviews, and automate hiring workflows from your terminal.

```bash
praxicraft-assess org get
praxicraft-assess assessments list --output table
praxicraft-assess invites create senior-backend-screen --email candidate@example.com
```

You need an organisation API key (`ct_live_…` or `ct_test_…`) from **Assess → Developer → API Keys**. Full API docs: [docs.praxicraft.com](https://docs.praxicraft.com).

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/praxicraft-platform/praxicraft-assess-cli/main/install.sh | sh
```

The script pulls the **latest** release binary for your OS/arch and installs `praxicraft-assess` to `/usr/local/bin` (or `~/.local/bin`).

```bash
# optional custom destination
curl -fsSL https://raw.githubusercontent.com/praxicraft-platform/praxicraft-assess-cli/main/install.sh | PRAXICRAFT_INSTALL_DIR="$HOME/bin" sh
```

### Manual download

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

## Interactive mode

Run with no arguments (in a terminal), or use `interactive`:

```bash
praxicraft-assess
# or
praxicraft-assess interactive
```

You see the Praxicraft ASCII wordmark, then a **command center**: press **1–9** to open a workspace (Assessments, Pipelines, Webhooks, …), then a **letter** to run an action (L list, G get, …). Esc/← goes back. After a command runs, press Enter for the menu again — or type any command normally.

**Menu covers resources such as:**

| Pick… | Runs |
|-------|------|
| assessments list | `assessments list` |
| assessments get | lists assessments → pick one → `get` |
| assessments results / cases list | pick assessment → list |
| invites list / create / result / remind / cancel | list or pick assessment/invite |
| results list / get | pick assessment or invite |
| cases list / platform-list | org or platform catalog |
| pipelines list / get / enrollments | list or pick pipeline |
| webhooks list / get / deliveries / test | list or pick webhook |
| org get / stats / team, whoami, interviews, integrations | read commands |
| configure | runner-style setup |

Outside the shell, the same pickers work when you omit IDs:

```bash
praxicraft-assess assessments get          # ↑/↓ pick assessment
praxicraft-assess invites create           # pick assessment, then email/name form
praxicraft-assess invites cancel           # pick invite, then Yes/No
praxicraft-assess webhooks test            # pick webhook
```

Use `--non-interactive` (and pass args/flags) in CI. Use `--yes` to skip confirms.

## Pagination

List endpoints return one page by default (`page_size` 20 on the API).

```bash
# One page
praxicraft-assess assessments list

# Every page (follows next until null; uses page_size=100)
praxicraft-assess assessments list --all

# Manual page
praxicraft-assess assessments list --filter page_size=50
praxicraft-assess assessments list --filter cursor=cD0yMDI1 --filter page_size=50
```

Interactive pickers (omit ID) show the first page, then **↓ Load more…** to walk cursors until you pick or cancel.

## Common commands

```bash
# Organisation
praxicraft-assess org get
praxicraft-assess org stats

# Assessments
praxicraft-assess assessments list
praxicraft-assess assessments get senior-backend-screen
praxicraft-assess assessments get   # interactive pick

# Invites
praxicraft-assess invites create senior-backend-screen --email candidate@example.com --name "Jane Doe"
praxicraft-assess invites create    # interactive pick + form
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
| `--all` | On list commands: fetch every page (follow `next`) |
| `--yes` | Skip confirmation prompts |
| `--non-interactive` | Disable interactive prompts (CI-friendly) |

## Examples

More recipes (SDKs, curl, n8n, Zapier) live in [praxicraft-assess-examples](https://github.com/praxicraft-platform/praxicraft-assess-examples).

## License

MIT
