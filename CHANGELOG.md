# Changelog

## 0.1.1

- Runner-style `configure` (ASCII Praxicraft banner, product/docs/API-key links, `[press Enter for …]` prompts, `√` confirmations).
- Aligned table output for list commands (short ids, sensible columns).
- Split API list `--filter key=value` from global JMESPath `--query`.
- Richer `--help` (Long/Examples on root; Short on leaf commands); removed unused pagination/debug flags.
- Broader destructive confirms; clearer missing-key / unknown-profile errors and exit codes.
- Auto-release workflow: pushes that change CLI source bump changelog/version, tag, and publish binaries.

## 0.1.0

- Initial release: full Public API command surface, interactive REPL, config profiles, JSON/table/yaml output.
