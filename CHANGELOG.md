# Changelog

## 0.1.4

- Interactive **command center**: numbered workspaces (1–9) → lettered actions; Praxicraft wordmark.
- Framed UI everywhere: pickers, confirms (Y/N), invite wizard, configure, errors, REPL done strip.
- Centered Praxicraft ASCII wordmark; single-object responses (e.g. invite create) print as a clear key/value table with full `invite_token` / `take_url`.
- Table mode unwraps nested list envelopes (`integrations`, `results`, …) instead of dumping raw JSON; `meta` scalars print below the table.
- Scenarios: empty lists, missing API key tips, overwrite profile confirm, abort/cancel, post-configure client refresh.

## 0.1.3

- `--all` on list commands follows cursor/`page` pagination and merges `results` (default `page_size=100`).
- Interactive pickers walk pages with **↓ Load more…** until you select or cancel.
- Manual paging still works: `--filter cursor=…` / `--filter page_size=…`.

## 0.1.2

### Interactive mode

- Running `praxicraft-assess` (no args) or `interactive` opens a full arrow-key resource menu (↑/↓, Enter): assessments, invites, results, cases, pipelines, webhooks, org, interviews, integrations, configure.
- Menu opens immediately on shell start; press Enter or type `menu` anytime to reopen it.
- Omit IDs/slugs on commands to pick from live API lists (assessments, invites, pipelines, webhooks).
- Invite create uses an interactive form (email, name, send email?).
- Destructive actions use arrow-key Yes/No confirms.
- `configure` stays runner-style (`[press Enter for …]`, `√` lines).
- `--non-interactive` / flags still work for CI and scripts.

### Docs

- README: Interactive mode section with resource menu table and omit-ID examples.

## 0.1.1

- Runner-style `configure` (ASCII Praxicraft banner, product/docs/API-key links, `[press Enter for …]` prompts, `√` confirmations).
- Aligned table output for list commands (short ids, sensible columns).
- Split API list `--filter key=value` from global JMESPath `--query`.
- Richer `--help` (Long/Examples on root; Short on leaf commands); removed unused pagination/debug flags.
- Broader destructive confirms; clearer missing-key / unknown-profile errors and exit codes.
- Auto-release workflow: pushes that change CLI source bump changelog/version, tag, and publish binaries.

## 0.1.0

- Initial release: full Public API command surface, interactive REPL, config profiles, JSON/table/yaml output.
