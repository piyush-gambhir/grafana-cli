# Grafana CLI -- Complete Command Reference

Every command, subcommand, and flag in the `grafana` CLI.

## Global Flags

Available on every command:

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: `table` (default), `json`, `yaml` |
| `--profile` | | Configuration profile to use |
| `--url` | | Grafana server URL override |
| `--token` | | API token or service account token override |
| `--username` | | Username for basic auth override |
| `--password` | | Password for basic auth override |
| `--org-id` | | Organization ID override |

---

## Top-Level Commands

```
grafana login                  Interactive login; saves connection profile to ~/.config/grafana-cli/config.yaml
grafana version                Print CLI version, commit hash, and build date
grafana update                 Check for and install CLI updates
grafana completion <shell>     Generate shell completion (bash, zsh, fish, powershell)
```

---

## grafana config

Manage CLI configuration.

```
grafana config view                        Display current configuration
grafana config set <key> <value>           Set a config value (keys: defaults.output, current_profile)
grafana config use-profile <name>          Switch to a different profile
grafana config list-profiles               List all configured profiles
```

---

## grafana dashboard

Aliases: `dash`, `db`

### grafana dashboard list

Search and list dashboards.

```
grafana dashboard list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--query` | `-q` | | Search query string (matches title) |
| `--tag` | `-t` | | Filter by tag |
| `--folder` | | | Filter by folder UID |
| `--page` | | `1` | Page number |
| `--limit` | | `100` | Results per page |

### grafana dashboard get

```
grafana dashboard get <uid>
```

Retrieve a dashboard by UID. Returns the full dashboard model.

### grafana dashboard create

```
grafana dashboard create -f <file> [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--file` | `-f` | | Path to JSON or YAML file (use `-` for stdin) |
| `--folder` | | | Folder UID to place the dashboard in |
| `--overwrite` | | `false` | Overwrite existing dashboard with same UID |
| `--message` | `-m` | | Commit message for version history |

### grafana dashboard update

```
grafana dashboard update -f <file> [flags]
```

Same flags as `create`, except `--overwrite` defaults to `true`.

### grafana dashboard delete

```
grafana dashboard delete <uid> [--confirm]
```

| Flag | Description |
|------|-------------|
| `--confirm` | Skip interactive confirmation prompt |

### grafana dashboard export

```
grafana dashboard export <uid> [flags]
```

| Flag | Description |
|------|-------------|
| `--output-file` | Write to file instead of stdout |

### grafana dashboard import

```
grafana dashboard import -f <file> [flags]
```

Alias for `create`. Same flags as `create`.

### grafana dashboard versions

```
grafana dashboard versions <uid> [--page N] [--limit N]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--limit` | `100` | Results per page |

### grafana dashboard restore

```
grafana dashboard restore <uid> <version>
```

Restore a dashboard to a specific historical version number.

### grafana dashboard permissions get

```
grafana dashboard permissions get <uid>
```

### grafana dashboard permissions update

```
grafana dashboard permissions update <uid> -f <file>
```

---

## grafana datasource

Alias: `ds`

### grafana datasource list

```
grafana datasource list [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--type` | | Filter by datasource type (e.g. `prometheus`, `loki`, `elasticsearch`, `mysql`) |
| `--name` | `-n` | Filter by name (case-insensitive substring match) |

### grafana datasource get

```
grafana datasource get <uid>
```

### grafana datasource create

```
grafana datasource create -f <file>
```

### grafana datasource update

```
grafana datasource update <id> -f <file>
```

Note: uses numeric ID, not UID.

### grafana datasource delete

```
grafana datasource delete <uid> [--confirm]
```

---

## grafana folder

### grafana folder list

```
grafana folder list [--page N] [--limit N]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--limit` | `100` | Results per page |

### grafana folder get

```
grafana folder get <uid>
```

### grafana folder create

```
grafana folder create -f <file>
```

### grafana folder update

```
grafana folder update <uid> -f <file>
```

### grafana folder delete

```
grafana folder delete <uid> [--confirm]
```

### grafana folder permissions get

```
grafana folder permissions get <uid>
```

### grafana folder permissions update

```
grafana folder permissions update <uid> -f <file>
```

---

## grafana alert rule

### grafana alert rule list

```
grafana alert rule list [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--folder` | | Filter by folder UID |
| `--group` | | Filter by rule group name |
| `--limit` | `0` (all) | Maximum number of rules |
| `--page` | `1` | Page number (used with --limit) |

### grafana alert rule get

```
grafana alert rule get <uid>
```

### grafana alert rule create

```
grafana alert rule create -f <file>
```

### grafana alert rule update

```
grafana alert rule update <uid> -f <file>
```

### grafana alert rule delete

```
grafana alert rule delete <uid> [--confirm]
```

---

## grafana alert contact-point

Alias: `cp`

### grafana alert contact-point list

```
grafana alert contact-point list
```

### grafana alert contact-point get

```
grafana alert contact-point get <uid>
```

### grafana alert contact-point create

```
grafana alert contact-point create -f <file>
```

### grafana alert contact-point update

```
grafana alert contact-point update <uid> -f <file>
```

### grafana alert contact-point delete

```
grafana alert contact-point delete <uid> [--confirm]
```

---

## grafana alert policy

### grafana alert policy get

```
grafana alert policy get
```

Returns the full notification policy tree.

### grafana alert policy update

```
grafana alert policy update -f <file>
```

### grafana alert policy reset

```
grafana alert policy reset [--confirm]
```

Reset notification policies to the default.

---

## grafana alert mute-timing

Alias: `mt`

### grafana alert mute-timing list

```
grafana alert mute-timing list
```

### grafana alert mute-timing get

```
grafana alert mute-timing get <name>
```

### grafana alert mute-timing create

```
grafana alert mute-timing create -f <file>
```

### grafana alert mute-timing update

```
grafana alert mute-timing update <name> -f <file>
```

### grafana alert mute-timing delete

```
grafana alert mute-timing delete <name> [--confirm]
```

---

## grafana alert template

Alias: `tmpl`

### grafana alert template list

```
grafana alert template list
```

### grafana alert template get

```
grafana alert template get <name>
```

### grafana alert template update

```
grafana alert template update <name> -f <file>
```

Creates or updates a notification template.

### grafana alert template delete

```
grafana alert template delete <name> [--confirm]
```

---

## grafana alert silence

### grafana alert silence list

```
grafana alert silence list
```

### grafana alert silence get

```
grafana alert silence get <id>
```

### grafana alert silence create

```
grafana alert silence create -f <file>
```

### grafana alert silence delete

```
grafana alert silence delete <id> [--confirm]
```

Expires (deletes) a silence.

---

## grafana org

### grafana org list

```
grafana org list [--page N] [--limit N]
```

### grafana org get

```
grafana org get <id>
```

### grafana org create

```
grafana org create -f <file>
```

### grafana org update

```
grafana org update <id> -f <file>
```

### grafana org delete

```
grafana org delete <id> [--confirm]
```

### grafana org current

```
grafana org current
```

Show the current organization context.

### grafana org switch

```
grafana org switch <org-id>
```

### grafana org user list

```
grafana org user list <org-id> [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--role` | | Filter by role: `Viewer`, `Editor`, `Admin` |
| `--query` | `-q` | Search by login, email, or name |

### grafana org user add

```
grafana org user add <org-id> -f <file>
```

JSON body: `{"loginOrEmail":"user@example.com","role":"Editor"}`

### grafana org user update

```
grafana org user update <org-id> <user-id> -f <file>
```

JSON body: `{"role":"Admin"}`

### grafana org user remove

```
grafana org user remove <org-id> <user-id> [--confirm]
```

---

## grafana team

### grafana team list

```
grafana team list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--query` | `-q` | | Search query |
| `--page` | | `1` | Page number |
| `--limit` | | `100` | Results per page |

### grafana team get

```
grafana team get <id>
```

### grafana team create

```
grafana team create -f <file>
```

JSON body: `{"name":"Team Name","email":"team@example.com"}`

### grafana team update

```
grafana team update <id> -f <file>
```

### grafana team delete

```
grafana team delete <id> [--confirm]
```

### grafana team member list

```
grafana team member list <team-id>
```

### grafana team member add

```
grafana team member add <team-id> <user-id>
```

### grafana team member remove

```
grafana team member remove <team-id> <user-id> [--confirm]
```

### grafana team preferences get

```
grafana team preferences get <team-id>
```

### grafana team preferences update

```
grafana team preferences update <team-id> -f <file>
```

---

## grafana user

### grafana user list

```
grafana user list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--query` | `-q` | | Search query |
| `--page` | | `1` | Page number |
| `--limit` | | `100` | Results per page |

### grafana user get

```
grafana user get <id>
```

### grafana user lookup

```
grafana user lookup <login-or-email>
```

### grafana user update

```
grafana user update <id> -f <file>
```

### grafana user orgs

```
grafana user orgs <user-id>
```

### grafana user teams

```
grafana user teams <user-id>
```

### grafana user current

```
grafana user current
```

Alias: `grafana user whoami`

### grafana user star add

```
grafana user star add <dashboard-id>
```

### grafana user star remove

```
grafana user star remove <dashboard-id>
```

---

## grafana service-account

Alias: `sa`

### grafana service-account list

```
grafana service-account list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--query` | `-q` | | Search query |
| `--page` | | `1` | Page number |
| `--limit` | | `100` | Results per page |

### grafana service-account get

```
grafana service-account get <id>
```

### grafana service-account create

```
grafana service-account create -f <file>
```

JSON body: `{"name":"sa-name","role":"Editor"}`

### grafana service-account update

```
grafana service-account update <id> -f <file>
```

### grafana service-account delete

```
grafana service-account delete <id> [--confirm]
```

### grafana service-account token list

```
grafana service-account token list <sa-id>
```

### grafana service-account token create

```
grafana service-account token create <sa-id> -f <file>
```

JSON body: `{"name":"token-name","secondsToLive":86400}`

The token key is only shown once in the response. Save it immediately.

### grafana service-account token delete

```
grafana service-account token delete <sa-id> <token-id> [--confirm]
```

---

## grafana annotation

### grafana annotation list

```
grafana annotation list [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--dashboard-id` | `0` | Filter by dashboard ID |
| `--panel-id` | `0` | Filter by panel ID |
| `--from` | `0` | Start time in epoch milliseconds |
| `--to` | `0` | End time in epoch milliseconds |
| `--tags` | | Filter by tags (comma-separated) |
| `--type` | | Filter by type: `annotation` or `alert` |
| `--limit` | `100` | Maximum number of results |

### grafana annotation get

```
grafana annotation get <id>
```

### grafana annotation create

```
grafana annotation create -f <file>
```

JSON body: `{"text":"message","tags":["tag1","tag2"]}`

### grafana annotation update

```
grafana annotation update <id> -f <file>
```

### grafana annotation delete

```
grafana annotation delete <id> [--confirm]
```

### grafana annotation tags

```
grafana annotation tags
```

List all unique annotation tags with usage counts.

---

## grafana snapshot

### grafana snapshot list

```
grafana snapshot list [--limit N]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `0` (all) | Maximum number of snapshots |

### grafana snapshot get

```
grafana snapshot get <key>
```

### grafana snapshot create

```
grafana snapshot create -f <file>
```

### grafana snapshot delete

```
grafana snapshot delete <key> [--confirm]
```

---

## grafana playlist

### grafana playlist list

```
grafana playlist list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--query` | `-q` | | Search query |
| `--limit` | | `0` (all) | Maximum number of results |

### grafana playlist get

```
grafana playlist get <uid>
```

### grafana playlist create

```
grafana playlist create -f <file>
```

### grafana playlist update

```
grafana playlist update <uid> -f <file>
```

### grafana playlist delete

```
grafana playlist delete <uid> [--confirm]
```

---

## grafana library-element

Alias: `le`

### grafana library-element list

```
grafana library-element list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--search` | `-q` | | Search string for element name |
| `--kind` | | `0` (all) | Kind: `1` = panel, `2` = variable |
| `--folder` | | | Filter by folder name |
| `--page` | | `1` | Page number |
| `--limit` | | `100` | Results per page |

### grafana library-element get

```
grafana library-element get <uid>
```

### grafana library-element create

```
grafana library-element create -f <file>
```

### grafana library-element update

```
grafana library-element update <uid> -f <file>
```

### grafana library-element delete

```
grafana library-element delete <uid> [--confirm]
```

### grafana library-element connections

```
grafana library-element connections <uid>
```

List dashboards connected to this library element.

---

## grafana correlation

### grafana correlation list

```
grafana correlation list
```

### grafana correlation get

```
grafana correlation get <source-uid> <correlation-uid>
```

### grafana correlation create

```
grafana correlation create <source-uid> -f <file>
```

### grafana correlation update

```
grafana correlation update <source-uid> <correlation-uid> -f <file>
```

### grafana correlation delete

```
grafana correlation delete <source-uid> <correlation-uid> [--confirm]
```

---

## grafana preferences

Alias: `prefs`

### grafana preferences get

```
grafana preferences get
```

### grafana preferences update

```
grafana preferences update -f <file>
```

JSON body: `{"theme":"dark","timezone":"utc","weekStart":"monday"}`

---

## grafana admin

Server administration (requires admin permissions).

### grafana admin settings

```
grafana admin settings
```

Display all Grafana server settings.

### grafana admin stats

```
grafana admin stats
```

Display server usage statistics.

### grafana admin reload

```
grafana admin reload <resource>
```

Supported resources: `dashboards`, `datasources`, `plugins`, `access-control`, `alerting`.
