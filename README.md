# Grafana CLI

A command-line interface for managing Grafana instances -- dashboards, datasources, folders, alerts, orgs, teams, users, service accounts, and more.

Designed for both human operators and coding agents (LLMs). All commands support `--help` for detailed usage, and `-o json` / `-o yaml` for machine-readable output.

[![Go Version](https://img.shields.io/github/go-mod/go-version/piyush-gambhir/grafana-cli)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/piyush-gambhir/grafana-cli)](https://github.com/piyush-gambhir/grafana-cli/releases)
[![License](https://img.shields.io/github/license/piyush-gambhir/grafana-cli)](LICENSE)
[![CI](https://github.com/piyush-gambhir/grafana-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/piyush-gambhir/grafana-cli/actions/workflows/ci.yml)

## Features

- Full API coverage — every Grafana API endpoint accessible from the command line
- Multiple output formats — table, JSON, YAML (`-o json`)
- Profile management — multiple instances with `--profile`
- Auto-update — checks for new versions, `grafana update` to self-update
- Agent-friendly — comprehensive help text, structured output for LLM coding agents
- Cross-platform — macOS, Linux, Windows (amd64 and arm64)

## Installation

```bash
# curl (recommended)
curl -sSfL https://raw.githubusercontent.com/piyush-gambhir/grafana-cli/main/install.sh | sh

# Go
go install github.com/piyush-gambhir/grafana-cli@latest

# From source
git clone https://github.com/piyush-gambhir/grafana-cli.git
cd grafana-cli && make install
```

## Quick Start

```bash
# Install
curl -sSfL https://raw.githubusercontent.com/piyush-gambhir/grafana-cli/main/install.sh | sh

# Authenticate
grafana login

# Start using
grafana dashboard list
grafana dashboard get <uid> -o json
```

## Authentication

```bash
# Interactive login (saves profile to ~/.config/grafana-cli/config.yaml)
grafana login

# Environment variables
export GRAFANA_URL=https://grafana.example.com
export GRAFANA_TOKEN=glsa_xxxx

# Multiple profiles
grafana login                        # saves as default profile
grafana config use-profile prod      # switch profiles
grafana dashboard list --profile staging
```

### Auth Priority

Configuration is resolved in this order (first match wins):

1. CLI flags (`--url`, `--token`, `--username`, `--password`, `--org-id`)
2. Environment variables (`GRAFANA_URL`, `GRAFANA_TOKEN`, `GRAFANA_USERNAME`, `GRAFANA_PASSWORD`, `GRAFANA_ORG_ID`)
3. Config file profile (`~/.config/grafana-cli/config.yaml`)

## Output Formats

All list and get commands support three output formats:

```bash
grafana dashboard list              # table (default, human-readable)
grafana dashboard list -o json      # JSON (machine-readable)
grafana dashboard list -o yaml      # YAML
```

## Global Flags

These flags are available on every command:

| Flag | Description |
|------|-------------|
| `--output`, `-o` | Output format: `table`, `json`, `yaml` |
| `--profile` | Configuration profile to use |
| `--url` | Grafana server URL |
| `--token` | API token or service account token |
| `--username` | Username for basic auth |
| `--password` | Password for basic auth |
| `--org-id` | Organization ID |

## Commands

### Dashboard

Manage Grafana dashboards. Alias: `dash`, `db`.

#### `grafana dashboard list`

Search and list dashboards with optional filters.

```bash
grafana dashboard list                          # list all
grafana dashboard list -q "production"          # search by title
grafana dashboard list --tag monitoring         # filter by tag
grafana dashboard list --folder abc123          # filter by folder UID
grafana dashboard list --page 2 --limit 50      # paginate
grafana dashboard list -o json                  # JSON output
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--query`, `-q` | | Search query string |
| `--tag`, `-t` | | Filter by tag |
| `--folder` | | Filter by folder UID |
| `--page` | `1` | Page number |
| `--limit` | `100` | Results per page |

#### `grafana dashboard get <uid>`

Retrieve a dashboard by UID. Returns the full dashboard model in JSON/YAML mode.

```bash
grafana dashboard get abc123
grafana dashboard get abc123 -o json
```

#### `grafana dashboard create`

Create a new dashboard from a JSON or YAML file.

```bash
grafana dashboard create -f dashboard.json
grafana dashboard create -f dashboard.json --folder folderUid123
grafana dashboard create -f dashboard.json --overwrite
grafana dashboard create -f dashboard.json -m "Initial version"
cat dashboard.json | grafana dashboard create -f -
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--file`, `-f` | | Path to JSON or YAML file (use `-` for stdin) |
| `--folder` | | Folder UID to place the dashboard in |
| `--overwrite` | `false` | Overwrite existing dashboard with same UID |
| `--message`, `-m` | | Commit message for version history |

#### `grafana dashboard update`

Update an existing dashboard from a JSON or YAML file.

```bash
grafana dashboard update -f dashboard.json
grafana dashboard update -f dashboard.json --folder newFolderUid
grafana dashboard update -f dashboard.json -m "Add new panels"
```

**Flags:** Same as `create`, except `--overwrite` defaults to `true`.

#### `grafana dashboard delete <uid>`

Permanently delete a dashboard by UID.

```bash
grafana dashboard delete abc123
grafana dashboard delete abc123 --confirm      # skip confirmation
```

#### `grafana dashboard export <uid>`

Export the full dashboard JSON to stdout or a file.

```bash
grafana dashboard export abc123
grafana dashboard export abc123 --output-file backup.json
grafana dashboard export abc123 | jq '.panels | length'
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--output-file` | Write to file instead of stdout |

#### `grafana dashboard import`

Import a dashboard from a JSON or YAML file (alias for `create`).

```bash
grafana dashboard import -f exported-dashboard.json
grafana dashboard import -f dashboard.json --folder folderUid --overwrite
```

#### `grafana dashboard versions <uid>`

List all versions of a dashboard.

```bash
grafana dashboard versions abc123
grafana dashboard versions abc123 --page 1 --limit 10
```

#### `grafana dashboard restore <uid> <version>`

Restore a dashboard to a specific historical version.

```bash
grafana dashboard restore abc123 3
```

#### `grafana dashboard permissions get <uid>`

Get current permissions for a dashboard.

```bash
grafana dashboard permissions get abc123
grafana dashboard permissions get abc123 -o json
```

#### `grafana dashboard permissions update <uid>`

Update dashboard permissions from a file.

```bash
grafana dashboard permissions update abc123 -f perms.json
```

---

### Datasource

Manage Grafana datasources. Alias: `ds`.

#### `grafana datasource list`

List all datasources with optional filters.

```bash
grafana datasource list                         # list all
grafana datasource list --type prometheus        # filter by type
grafana datasource list --name "prod"            # search by name
grafana datasource list --type loki --name "staging"
grafana datasource list -o json
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--type` | Filter by datasource type (e.g. `prometheus`, `elasticsearch`, `loki`, `mysql`) |
| `--name`, `-n` | Filter by name (case-insensitive substring match) |

#### `grafana datasource get <uid>`

Retrieve a datasource by UID.

```bash
grafana datasource get P1234
grafana datasource get P1234 -o json
```

#### `grafana datasource create`

Create a datasource from a JSON or YAML file.

```bash
grafana datasource create -f prometheus.json
```

#### `grafana datasource update <id>`

Update a datasource by numeric ID.

```bash
grafana datasource update 5 -f updated-ds.json
```

#### `grafana datasource delete <uid>`

Delete a datasource by UID.

```bash
grafana datasource delete P1234
grafana datasource delete P1234 --confirm
```

---

### Folder

Manage Grafana folders.

#### `grafana folder list`

List all folders.

```bash
grafana folder list
grafana folder list --page 1 --limit 50
grafana folder list -o json
```

#### `grafana folder get <uid>`

```bash
grafana folder get folderUid123
```

#### `grafana folder create`

```bash
grafana folder create -f folder.json
```

#### `grafana folder update <uid>`

```bash
grafana folder update folderUid123 -f updated-folder.json
```

#### `grafana folder delete <uid>`

```bash
grafana folder delete folderUid123
grafana folder delete folderUid123 --confirm
```

#### `grafana folder permissions get <uid>`

```bash
grafana folder permissions get folderUid123
```

#### `grafana folder permissions update <uid>`

```bash
grafana folder permissions update folderUid123 -f perms.json
```

---

### Alert

Manage Grafana Unified Alerting resources: rules, contact points, policies, mute timings, templates, and silences.

#### Alert Rule

##### `grafana alert rule list`

List alert rules with optional filters.

```bash
grafana alert rule list                         # list all
grafana alert rule list --folder abc123          # filter by folder UID
grafana alert rule list --group "High CPU"       # filter by rule group
grafana alert rule list --limit 10 --page 2      # paginate
grafana alert rule list -o json
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--folder` | | Filter by folder UID |
| `--group` | | Filter by rule group name |
| `--limit` | `0` (all) | Maximum number of rules |
| `--page` | `1` | Page number (used with --limit) |

##### `grafana alert rule get <uid>`

```bash
grafana alert rule get ruleUid123
grafana alert rule get ruleUid123 -o json
```

##### `grafana alert rule create`

```bash
grafana alert rule create -f rule.json
```

##### `grafana alert rule update <uid>`

```bash
grafana alert rule update ruleUid123 -f updated-rule.json
```

##### `grafana alert rule delete <uid>`

```bash
grafana alert rule delete ruleUid123
grafana alert rule delete ruleUid123 --confirm
```

#### Contact Point

Alias: `cp`.

##### `grafana alert contact-point list`

```bash
grafana alert contact-point list
grafana alert contact-point list -o json
```

##### `grafana alert contact-point get <uid>`

```bash
grafana alert contact-point get cpUid123
```

##### `grafana alert contact-point create`

```bash
grafana alert contact-point create -f contact-point.json
```

##### `grafana alert contact-point update <uid>`

```bash
grafana alert contact-point update cpUid123 -f updated-cp.json
```

##### `grafana alert contact-point delete <uid>`

```bash
grafana alert contact-point delete cpUid123 --confirm
```

#### Notification Policy

##### `grafana alert policy get`

```bash
grafana alert policy get
grafana alert policy get -o json
```

##### `grafana alert policy update`

```bash
grafana alert policy update -f policy.json
```

##### `grafana alert policy reset`

```bash
grafana alert policy reset
grafana alert policy reset --confirm
```

#### Mute Timing

Alias: `mt`.

##### `grafana alert mute-timing list`

```bash
grafana alert mute-timing list
```

##### `grafana alert mute-timing get <name>`

```bash
grafana alert mute-timing get "weekends"
```

##### `grafana alert mute-timing create`

```bash
grafana alert mute-timing create -f mute-timing.json
```

##### `grafana alert mute-timing update <name>`

```bash
grafana alert mute-timing update "weekends" -f updated-mt.json
```

##### `grafana alert mute-timing delete <name>`

```bash
grafana alert mute-timing delete "weekends" --confirm
```

#### Template

Alias: `tmpl`.

##### `grafana alert template list`

```bash
grafana alert template list
```

##### `grafana alert template get <name>`

```bash
grafana alert template get "my-template"
```

##### `grafana alert template update <name>`

Creates or updates a notification template.

```bash
grafana alert template update "my-template" -f template.json
```

##### `grafana alert template delete <name>`

```bash
grafana alert template delete "my-template" --confirm
```

#### Silence

##### `grafana alert silence list`

```bash
grafana alert silence list
grafana alert silence list -o json
```

##### `grafana alert silence get <id>`

```bash
grafana alert silence get silenceId123
```

##### `grafana alert silence create`

```bash
grafana alert silence create -f silence.json
```

##### `grafana alert silence delete <id>`

```bash
grafana alert silence delete silenceId123 --confirm
```

---

### Organization

Manage Grafana organizations.

#### `grafana org list`

```bash
grafana org list
grafana org list --page 1 --limit 50
```

#### `grafana org get <id>`

```bash
grafana org get 1
```

#### `grafana org create`

```bash
grafana org create -f org.json
# Example: echo '{"name":"My Org"}' | grafana org create -f -
```

#### `grafana org update <id>`

```bash
grafana org update 2 -f org.json
```

#### `grafana org delete <id>`

```bash
grafana org delete 2 --confirm
```

#### `grafana org current`

```bash
grafana org current
```

#### `grafana org switch <org-id>`

```bash
grafana org switch 2
```

#### Organization Users

##### `grafana org user list <org-id>`

List users in an organization with optional filters.

```bash
grafana org user list 1                          # list all users in org 1
grafana org user list 1 --role Admin             # filter by role
grafana org user list 1 --query "john"           # search by name/email/login
grafana org user list 1 --role Editor -q "dev"   # combine filters
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--role` | Filter by role: `Viewer`, `Editor`, `Admin` |
| `--query`, `-q` | Search by login, email, or name |

##### `grafana org user add <org-id>`

```bash
grafana org user add 1 -f user.json
# JSON: {"loginOrEmail":"admin@example.com","role":"Editor"}
```

##### `grafana org user update <org-id> <user-id>`

```bash
grafana org user update 1 5 -f role.json
# JSON: {"role":"Admin"}
```

##### `grafana org user remove <org-id> <user-id>`

```bash
grafana org user remove 1 5 --confirm
```

---

### Team

Manage Grafana teams, members, and preferences.

#### `grafana team list`

```bash
grafana team list
grafana team list -q "backend"
grafana team list --page 1 --limit 20
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--query`, `-q` | | Search query |
| `--page` | `1` | Page number |
| `--limit` | `100` | Results per page |

#### `grafana team get <id>`

```bash
grafana team get 5
```

#### `grafana team create`

```bash
grafana team create -f team.json
# JSON: {"name":"Backend Team","email":"backend@example.com"}
```

#### `grafana team update <id>`

```bash
grafana team update 5 -f team.json
```

#### `grafana team delete <id>`

```bash
grafana team delete 5 --confirm
```

#### Team Members

##### `grafana team member list <team-id>`

```bash
grafana team member list 5
```

##### `grafana team member add <team-id> <user-id>`

```bash
grafana team member add 5 10
```

##### `grafana team member remove <team-id> <user-id>`

```bash
grafana team member remove 5 10 --confirm
```

#### Team Preferences

##### `grafana team preferences get <team-id>`

```bash
grafana team preferences get 5
```

##### `grafana team preferences update <team-id>`

```bash
grafana team preferences update 5 -f prefs.json
```

---

### User

Manage Grafana users (most commands require server admin).

#### `grafana user list`

```bash
grafana user list
grafana user list -q "john"
grafana user list --page 2 --limit 50
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--query`, `-q` | | Search query |
| `--page` | `1` | Page number |
| `--limit` | `100` | Results per page |

#### `grafana user get <id>`

```bash
grafana user get 5
```

#### `grafana user lookup <login-or-email>`

```bash
grafana user lookup admin
grafana user lookup admin@example.com
```

#### `grafana user update <id>`

```bash
grafana user update 5 -f user.json
```

#### `grafana user orgs <user-id>`

```bash
grafana user orgs 5
```

#### `grafana user teams <user-id>`

```bash
grafana user teams 5
```

#### `grafana user current`

Show the currently authenticated user. Alias: `whoami`.

```bash
grafana user current
grafana user whoami
```

#### `grafana user star add <dashboard-id>`

```bash
grafana user star add 42
```

#### `grafana user star remove <dashboard-id>`

```bash
grafana user star remove 42
```

---

### Service Account

Manage service accounts and their API tokens. Alias: `sa`.

#### `grafana service-account list`

```bash
grafana service-account list
grafana service-account list -q "ci-bot"
grafana service-account list --page 1 --limit 20
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--query`, `-q` | | Search query |
| `--page` | `1` | Page number |
| `--limit` | `100` | Results per page |

#### `grafana service-account get <id>`

```bash
grafana service-account get 10
```

#### `grafana service-account create`

```bash
grafana service-account create -f sa.json
# JSON: {"name":"ci-bot","role":"Editor"}
```

#### `grafana service-account update <id>`

```bash
grafana service-account update 10 -f sa.json
```

#### `grafana service-account delete <id>`

```bash
grafana service-account delete 10 --confirm
```

#### Service Account Tokens

##### `grafana service-account token list <sa-id>`

```bash
grafana service-account token list 10
```

##### `grafana service-account token create <sa-id>`

The token key is only shown once. Save it immediately.

```bash
grafana service-account token create 10 -f token.json
# JSON: {"name":"deploy-token","secondsToLive":86400}
```

##### `grafana service-account token delete <sa-id> <token-id>`

```bash
grafana service-account token delete 10 3 --confirm
```

---

### Annotation

Manage annotations and annotation tags.

#### `grafana annotation list`

List annotations with optional filters.

```bash
grafana annotation list                                         # default limit 100
grafana annotation list --dashboard-id 42                        # by dashboard
grafana annotation list --dashboard-id 42 --panel-id 3           # by panel
grafana annotation list --from 1609459200000 --to 1609545600000  # time range
grafana annotation list --tags deploy,release                    # by tags
grafana annotation list --type alert                             # by type
grafana annotation list --limit 500                              # increase limit
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--dashboard-id` | `0` | Filter by dashboard ID |
| `--panel-id` | `0` | Filter by panel ID |
| `--from` | `0` | Start time (epoch ms) |
| `--to` | `0` | End time (epoch ms) |
| `--tags` | | Filter by tags (comma-separated) |
| `--type` | | Filter by type: `annotation` or `alert` |
| `--limit` | `100` | Maximum number of results |

#### `grafana annotation get <id>`

```bash
grafana annotation get 42
```

#### `grafana annotation create`

```bash
grafana annotation create -f annotation.json
# JSON: {"text":"Deployed v1.2.3","tags":["deploy"]}
```

#### `grafana annotation update <id>`

```bash
grafana annotation update 42 -f annotation.json
```

#### `grafana annotation delete <id>`

```bash
grafana annotation delete 42 --confirm
```

#### `grafana annotation tags`

List all unique annotation tags with usage counts.

```bash
grafana annotation tags
```

---

### Snapshot

Manage dashboard snapshots.

#### `grafana snapshot list`

```bash
grafana snapshot list
grafana snapshot list --limit 10
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `0` (all) | Maximum number of snapshots |

#### `grafana snapshot get <key>`

```bash
grafana snapshot get abc123key
```

#### `grafana snapshot create`

```bash
grafana snapshot create -f snapshot.json
```

#### `grafana snapshot delete <key>`

```bash
grafana snapshot delete abc123key --confirm
```

---

### Playlist

Manage dashboard playlists.

#### `grafana playlist list`

```bash
grafana playlist list
grafana playlist list --query "production"
grafana playlist list --limit 10
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--query`, `-q` | | Search query |
| `--limit` | `0` (all) | Maximum number of results |

#### `grafana playlist get <uid>`

```bash
grafana playlist get playlistUid
```

#### `grafana playlist create`

```bash
grafana playlist create -f playlist.json
```

#### `grafana playlist update <uid>`

```bash
grafana playlist update playlistUid -f playlist.json
```

#### `grafana playlist delete <uid>`

```bash
grafana playlist delete playlistUid --confirm
```

---

### Library Element

Manage reusable library panels and variables. Alias: `le`.

#### `grafana library-element list`

```bash
grafana library-element list                     # list all
grafana library-element list --kind 1            # panels only
grafana library-element list --kind 2            # variables only
grafana library-element list --search "CPU"      # search by name
grafana library-element list --folder "General"  # filter by folder
grafana library-element list --page 2 --limit 20
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--search`, `-q` | | Search string for element name |
| `--kind` | `0` (all) | Kind: `1` = panel, `2` = variable |
| `--folder` | | Filter by folder name |
| `--page` | `1` | Page number |
| `--limit` | `100` | Results per page |

#### `grafana library-element get <uid>`

```bash
grafana library-element get leUid123
```

#### `grafana library-element create`

```bash
grafana library-element create -f panel.json
```

#### `grafana library-element update <uid>`

```bash
grafana library-element update leUid123 -f updated-panel.json
```

#### `grafana library-element delete <uid>`

```bash
grafana library-element delete leUid123 --confirm
```

#### `grafana library-element connections <uid>`

List dashboards connected to a library element.

```bash
grafana library-element connections leUid123
```

---

### Correlation

Manage datasource correlations.

#### `grafana correlation list`

```bash
grafana correlation list
```

#### `grafana correlation get <source-uid> <correlation-uid>`

```bash
grafana correlation get sourceUid corrUid
```

#### `grafana correlation create <source-uid>`

```bash
grafana correlation create sourceUid -f correlation.json
```

#### `grafana correlation update <source-uid> <correlation-uid>`

```bash
grafana correlation update sourceUid corrUid -f correlation.json
```

#### `grafana correlation delete <source-uid> <correlation-uid>`

```bash
grafana correlation delete sourceUid corrUid --confirm
```

---

### Admin

Server administration commands (require admin permissions).

#### `grafana admin settings`

Display all Grafana server settings.

```bash
grafana admin settings
grafana admin settings -o json
```

#### `grafana admin stats`

Display server usage statistics.

```bash
grafana admin stats
grafana admin stats -o json
```

#### `grafana admin reload <resource>`

Reload provisioned resources. Supported: `dashboards`, `datasources`, `plugins`, `access-control`, `alerting`.

```bash
grafana admin reload dashboards
grafana admin reload datasources
grafana admin reload alerting
```

---

### Preferences

Manage user preferences. Alias: `prefs`.

#### `grafana preferences get`

```bash
grafana preferences get
grafana preferences get -o json
```

#### `grafana preferences update`

```bash
grafana preferences update -f prefs.json
# JSON: {"theme":"dark","timezone":"utc","weekStart":"monday"}
```

---

### Config

Manage CLI configuration.

#### `grafana config view`

Display the current configuration file.

```bash
grafana config view
```

#### `grafana config set <key> <value>`

Set a configuration value. Supported keys: `defaults.output`, `current_profile`.

```bash
grafana config set defaults.output json
grafana config set current_profile prod
```

#### `grafana config use-profile <name>`

Switch to a different profile.

```bash
grafana config use-profile prod
```

#### `grafana config list-profiles`

List all configured profiles.

```bash
grafana config list-profiles
```

---

### Other Commands

#### `grafana login`

Interactively log in to a Grafana instance and save credentials as a profile.

```bash
grafana login
```

#### `grafana version`

Print version, commit, and build date.

```bash
grafana version
```

#### `grafana completion <shell>`

Generate shell completion scripts for bash, zsh, fish, or PowerShell.

```bash
grafana completion bash
grafana completion zsh
```

#### `grafana update`

Check for and install CLI updates.

```bash
grafana update
```

## File Input Format

Commands that accept `--file/-f` support:
- JSON files (`.json`)
- YAML files (`.yaml`, `.yml`)
- Stdin (use `-f -` and pipe input)

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Agent Skills

This CLI ships with an agent skill for coding agents (Claude, Cursor, Copilot, etc.):

```bash
npx skills add piyush-gambhir/grafana-cli@grafana
```

Once installed, coding agents automatically know how to use this CLI effectively.

## License

MIT
