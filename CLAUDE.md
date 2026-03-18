# Grafana CLI - Agent Guide

## Quick Reference

- **Binary:** `grafana`
- **Config file:** `~/.config/grafana-cli/config.yaml`
- **Env vars:** `GRAFANA_URL`, `GRAFANA_TOKEN`, `GRAFANA_USERNAME`, `GRAFANA_PASSWORD`, `GRAFANA_ORG_ID`
- **Auth methods:** API token (service account token) or basic auth (username/password)
- **Config priority:** CLI flags > environment variables > profile config > defaults

## Setup

```bash
# Interactive login (prompts for URL, auth method, credentials, profile name)
grafana login

# Or set environment variables for non-interactive use
export GRAFANA_URL=https://grafana.example.com
export GRAFANA_TOKEN=glsa_xxxxxxxxxxxxxxxxxxxx
```

## Output Formats

All list/get commands support three output formats via `-o`:

- `-o table` (default) -- human-readable tabular output
- `-o json` -- JSON, ideal for programmatic parsing with jq
- `-o yaml` -- YAML, useful for config management

**For agents:** Always use `-o json` when you need to parse or process output programmatically.

## Common Workflows

### Find and inspect dashboards

```bash
# List all dashboards
grafana dashboard list -o json

# Search dashboards by title
grafana dashboard list -q "production" -o json

# Filter by tag
grafana dashboard list --tag monitoring -o json

# Filter by folder UID
grafana dashboard list --folder <folder-uid> -o json

# Get full dashboard details by UID
grafana dashboard get <uid> -o json

# Get dashboard versions history
grafana dashboard versions <uid> -o json
```

### Export and import dashboards (backup/restore)

```bash
# Export dashboard JSON to stdout
grafana dashboard export <uid>

# Export to a file
grafana dashboard export <uid> --output-file dashboard-backup.json

# Import from file
grafana dashboard import -f dashboard.json

# Import into a specific folder, overwriting if exists
grafana dashboard import -f dashboard.json --folder <folder-uid> --overwrite

# Import with a version commit message
grafana dashboard import -f dashboard.json --overwrite -m "Restored from backup"
```

### Manage datasources

```bash
# List all datasources
grafana datasource list -o json

# Get a specific datasource by UID or ID
grafana datasource get <uid-or-id> -o json

# Create a datasource from a JSON/YAML file
grafana datasource create -f datasource.json

# Update a datasource
grafana datasource update <uid> -f datasource.json

# Delete a datasource
grafana datasource delete <uid>
```

### Manage alerts

```bash
# List alert rules
grafana alert rule list -o json

# Get a specific alert rule
grafana alert rule get <uid> -o json

# Create an alert rule from file
grafana alert rule create -f rule.json

# List contact points
grafana alert contact-point list -o json

# Get notification policy tree
grafana alert policy get -o json

# List active silences
grafana alert silence list -o json

# Create a silence from file
grafana alert silence create -f silence.json

# List mute timings
grafana alert mute-timing list -o json

# List notification templates
grafana alert template list -o json
```

### Service accounts and tokens (for automation)

```bash
# List service accounts
grafana service-account list -o json

# Create a service account (from JSON file with name and role)
# Example JSON: {"name":"ci-bot","role":"Editor"}
grafana service-account create -f sa.json

# Get service account details
grafana service-account get <id> -o json

# List tokens for a service account
grafana service-account token list <sa-id> -o json

# Create a new token for a service account
grafana service-account token create <sa-id> -f token.json

# Delete a token
grafana service-account token delete <sa-id> <token-id>
```

### Annotations (mark deployments, incidents)

```bash
# List annotations (default limit 100)
grafana annotation list -o json

# Filter annotations by tags
grafana annotation list --tags deploy,release -o json

# Filter by dashboard ID and time range (epoch ms)
grafana annotation list --dashboard-id 42 --from 1705312800000 --to 1705399200000 -o json

# Filter by type (annotation or alert)
grafana annotation list --type alert -o json

# Increase result limit
grafana annotation list --limit 500 -o json

# Create an annotation (from JSON file)
# Example JSON: {"text":"Deployed v1.2.3","tags":["deploy","production"]}
grafana annotation create -f annotation.json

# Get a specific annotation
grafana annotation get <id> -o json

# List annotation tags
grafana annotation tags -o json
```

### User and org management

```bash
# List all users
grafana user list -o json

# Get user details by ID
grafana user get <id> -o json

# Look up user by login or email
grafana user lookup --login admin -o json

# Show current authenticated user
grafana user current -o json

# List organizations
grafana org list -o json

# Get current org
grafana org current -o json

# Switch to a different org
grafana org switch <org-id>

# List teams
grafana team list -o json

# Get team details
grafana team get <id> -o json

# List team members
grafana team member list <team-id> -o json
```

### Folder management

```bash
# List all folders
grafana folder list -o json

# Get folder details
grafana folder get <uid> -o json

# Create a folder from file
grafana folder create -f folder.json

# View folder permissions
grafana folder permissions <uid> -o json
```

### Configuration management

```bash
# View current configuration
grafana config view

# Set default output format
grafana config set defaults.output json

# List all profiles
grafana config list-profiles

# Switch to a different profile
grafana config use-profile staging
```

## Tips for Agents

- Always use `-o json` when you need to parse output programmatically.
- Dashboard UIDs are more stable than numeric IDs; prefer UIDs for cross-instance references.
- Use `--query` or `-q` flags on `dashboard list` to filter results server-side before parsing.
- For bulk operations: list with `-o json`, parse with jq, then loop over results.
- Many create/update commands require a `-f` flag pointing to a JSON or YAML file. Prepare the file first, then pass it.
- Pagination is available on list commands via `--page` and `--limit` flags.
- The `dashboard` command has aliases: `dash`, `db`. The `datasource` command has alias `ds`. The `service-account` command has alias `sa`. The `library-element` command has alias `le`. The `preferences` command has alias `prefs`.
- Use `grafana config view` to check current connection settings and confirm which profile is active.
- Annotation time values (`--from`, `--to`) are in epoch milliseconds, not seconds.

## Complete Command Reference

### Top-level commands

| Command | Description |
|---------|-------------|
| `grafana login` | Interactively log in and save a connection profile |
| `grafana version` | Print CLI version |
| `grafana update` | Check for and install CLI updates |
| `grafana completion` | Generate shell completion scripts |

### `grafana config` -- Manage CLI configuration

| Command | Description |
|---------|-------------|
| `grafana config view` | Display the current configuration |
| `grafana config set <key> <value>` | Set a configuration value (defaults.output, current_profile) |
| `grafana config use-profile <name>` | Switch to a different profile |
| `grafana config list-profiles` | List all configured profiles |

### `grafana dashboard` (aliases: `dash`, `db`) -- Manage dashboards

| Command | Description |
|---------|-------------|
| `grafana dashboard list` | Search and list dashboards (flags: -q, --tag, --folder, --page, --limit) |
| `grafana dashboard get <uid>` | Get a dashboard by UID |
| `grafana dashboard create -f <file>` | Create a dashboard from a JSON/YAML file |
| `grafana dashboard update <uid> -f <file>` | Update a dashboard from file |
| `grafana dashboard delete <uid>` | Delete a dashboard |
| `grafana dashboard export <uid>` | Export dashboard JSON (--output-file for file) |
| `grafana dashboard import -f <file>` | Import a dashboard (--folder, --overwrite, -m) |
| `grafana dashboard versions <uid>` | List dashboard version history |
| `grafana dashboard restore <uid>` | Restore a dashboard to a previous version |
| `grafana dashboard permissions <uid>` | View/manage dashboard permissions |

### `grafana datasource` (alias: `ds`) -- Manage datasources

| Command | Description |
|---------|-------------|
| `grafana datasource list` | List all datasources |
| `grafana datasource get <uid-or-id>` | Get a datasource by UID or ID |
| `grafana datasource create -f <file>` | Create a datasource from file |
| `grafana datasource update <uid> -f <file>` | Update a datasource |
| `grafana datasource delete <uid>` | Delete a datasource |

### `grafana folder` -- Manage folders

| Command | Description |
|---------|-------------|
| `grafana folder list` | List all folders |
| `grafana folder get <uid>` | Get folder details |
| `grafana folder create -f <file>` | Create a folder |
| `grafana folder update <uid> -f <file>` | Update a folder |
| `grafana folder delete <uid>` | Delete a folder |
| `grafana folder permissions <uid>` | View/manage folder permissions |

### `grafana alert` -- Manage alerting resources

#### `grafana alert rule` -- Alert rules

| Command | Description |
|---------|-------------|
| `grafana alert rule list` | List alert rules |
| `grafana alert rule get <uid>` | Get an alert rule |
| `grafana alert rule create -f <file>` | Create an alert rule |
| `grafana alert rule update <uid> -f <file>` | Update an alert rule |
| `grafana alert rule delete <uid>` | Delete an alert rule |

#### `grafana alert contact-point` (alias: `cp`) -- Contact points

| Command | Description |
|---------|-------------|
| `grafana alert contact-point list` | List contact points |
| `grafana alert contact-point get <uid>` | Get a contact point |
| `grafana alert contact-point create -f <file>` | Create a contact point |
| `grafana alert contact-point update <uid> -f <file>` | Update a contact point |
| `grafana alert contact-point delete <uid>` | Delete a contact point |

#### `grafana alert policy` -- Notification policies

| Command | Description |
|---------|-------------|
| `grafana alert policy get` | Get the notification policy tree |
| `grafana alert policy update -f <file>` | Update the notification policy tree |
| `grafana alert policy reset` | Reset policies to default |

#### `grafana alert mute-timing` (alias: `mt`) -- Mute timings

| Command | Description |
|---------|-------------|
| `grafana alert mute-timing list` | List mute timings |
| `grafana alert mute-timing get <name>` | Get a mute timing |
| `grafana alert mute-timing create -f <file>` | Create a mute timing |
| `grafana alert mute-timing update <name> -f <file>` | Update a mute timing |
| `grafana alert mute-timing delete <name>` | Delete a mute timing |

#### `grafana alert template` (alias: `tmpl`) -- Notification templates

| Command | Description |
|---------|-------------|
| `grafana alert template list` | List notification templates |
| `grafana alert template get <name>` | Get a template |
| `grafana alert template update <name> -f <file>` | Update a template |
| `grafana alert template delete <name>` | Delete a template |

#### `grafana alert silence` -- Alert silences

| Command | Description |
|---------|-------------|
| `grafana alert silence list` | List silences |
| `grafana alert silence get <id>` | Get a silence |
| `grafana alert silence create -f <file>` | Create a silence |
| `grafana alert silence delete <id>` | Delete (expire) a silence |

### `grafana org` -- Manage organizations

| Command | Description |
|---------|-------------|
| `grafana org list` | List all organizations |
| `grafana org get <id>` | Get organization details |
| `grafana org create -f <file>` | Create an organization |
| `grafana org update <id> -f <file>` | Update an organization |
| `grafana org delete <id>` | Delete an organization |
| `grafana org current` | Show the current organization |
| `grafana org switch <id>` | Switch to a different organization |
| `grafana org user list <org-id>` | List users in an organization |
| `grafana org user add <org-id> -f <file>` | Add a user to an organization |
| `grafana org user update <org-id> <user-id>` | Update user role in an organization |
| `grafana org user remove <org-id> <user-id>` | Remove a user from an organization |

### `grafana team` -- Manage teams

| Command | Description |
|---------|-------------|
| `grafana team list` | List all teams |
| `grafana team get <id>` | Get team details |
| `grafana team create -f <file>` | Create a team |
| `grafana team update <id> -f <file>` | Update a team |
| `grafana team delete <id>` | Delete a team |
| `grafana team preferences <id>` | View team preferences |
| `grafana team member list <team-id>` | List team members |
| `grafana team member add <team-id> -f <file>` | Add a member to a team |
| `grafana team member remove <team-id> <user-id>` | Remove a member from a team |

### `grafana user` -- Manage users

| Command | Description |
|---------|-------------|
| `grafana user list` | List all users |
| `grafana user get <id>` | Get user details |
| `grafana user lookup` | Look up user by login or email (--login, --email) |
| `grafana user update <id> -f <file>` | Update a user |
| `grafana user orgs <id>` | List organizations a user belongs to |
| `grafana user teams <id>` | List teams a user belongs to |
| `grafana user current` | Show the current authenticated user |
| `grafana user star` | Star/unstar dashboards |

### `grafana service-account` (alias: `sa`) -- Manage service accounts

| Command | Description |
|---------|-------------|
| `grafana service-account list` | List service accounts |
| `grafana service-account get <id>` | Get service account details |
| `grafana service-account create -f <file>` | Create a service account |
| `grafana service-account update <id> -f <file>` | Update a service account |
| `grafana service-account delete <id>` | Delete a service account |
| `grafana service-account token list <sa-id>` | List tokens for a service account |
| `grafana service-account token create <sa-id> -f <file>` | Create a token |
| `grafana service-account token delete <sa-id> <token-id>` | Delete a token |

### `grafana annotation` -- Manage annotations

| Command | Description |
|---------|-------------|
| `grafana annotation list` | List annotations (--dashboard-id, --panel-id, --from, --to, --tags, --limit, --type) |
| `grafana annotation get <id>` | Get an annotation |
| `grafana annotation create -f <file>` | Create an annotation |
| `grafana annotation update <id> -f <file>` | Update an annotation |
| `grafana annotation delete <id>` | Delete an annotation |
| `grafana annotation tags` | List annotation tags |

### `grafana snapshot` -- Manage dashboard snapshots

| Command | Description |
|---------|-------------|
| `grafana snapshot list` | List snapshots |
| `grafana snapshot get <key>` | Get a snapshot |
| `grafana snapshot create -f <file>` | Create a snapshot |
| `grafana snapshot delete <key>` | Delete a snapshot |

### `grafana playlist` -- Manage playlists

| Command | Description |
|---------|-------------|
| `grafana playlist list` | List playlists |
| `grafana playlist get <uid>` | Get a playlist |
| `grafana playlist create -f <file>` | Create a playlist |
| `grafana playlist update <uid> -f <file>` | Update a playlist |
| `grafana playlist delete <uid>` | Delete a playlist |

### `grafana library-element` (alias: `le`) -- Manage library elements

| Command | Description |
|---------|-------------|
| `grafana library-element list` | List library elements |
| `grafana library-element get <uid>` | Get a library element |
| `grafana library-element create -f <file>` | Create a library element |
| `grafana library-element update <uid> -f <file>` | Update a library element |
| `grafana library-element delete <uid>` | Delete a library element |
| `grafana library-element connections <uid>` | List dashboards connected to this element |

### `grafana correlation` -- Manage datasource correlations

| Command | Description |
|---------|-------------|
| `grafana correlation list` | List correlations |
| `grafana correlation get <uid>` | Get a correlation |
| `grafana correlation create -f <file>` | Create a correlation |
| `grafana correlation update <uid> -f <file>` | Update a correlation |
| `grafana correlation delete <uid>` | Delete a correlation |

### `grafana preferences` (alias: `prefs`) -- Manage user preferences

| Command | Description |
|---------|-------------|
| `grafana preferences get` | View current user preferences |
| `grafana preferences update -f <file>` | Update user preferences |

### `grafana admin` -- Server administration (requires admin permissions)

| Command | Description |
|---------|-------------|
| `grafana admin settings` | View server settings |
| `grafana admin stats` | View server usage statistics |
| `grafana admin reload` | Reload server provisioning configurations |

## Global Flags

| Flag | Description |
|------|-------------|
| `-o, --output <format>` | Output format: table (default), json, yaml |
| `--profile <name>` | Configuration profile to use |
| `--url <url>` | Grafana server URL override |
| `--token <token>` | API token override |
| `--username <user>` | Username for basic auth override |
| `--password <pass>` | Password for basic auth override |
| `--org-id <id>` | Organization ID override |
