---
name: grafana
description: "Expert guide for using the grafana CLI to manage Grafana instances. Use this skill whenever the user mentions Grafana dashboards, datasources, alerts, alert rules, contact points, silences, mute timings, notification policies, notification templates, organizations, teams, users, service accounts, API tokens, annotations, snapshots, playlists, library elements, correlations, folders, preferences, or Grafana admin operations. Also trigger when the user wants to query Grafana, search dashboards, export or import dashboards, manage Grafana permissions, set up monitoring, configure alerting, create or manage service account tokens, mark deployments with annotations, automate any Grafana operations from the command line, back up or restore dashboards, manage Grafana org users, switch Grafana orgs, reload provisioning, or check Grafana server stats. This skill provides the exact CLI commands, flags, and workflows needed to accomplish any Grafana management task."
---

# Grafana CLI Skill

## Prerequisites and Setup

```bash
# Check if grafana CLI is installed
grafana version

# If not installed:
curl -sSfL https://raw.githubusercontent.com/piyush-gambhir/grafana-cli/main/install.sh | sh

# Interactive login (prompts for URL, auth method, credentials, profile name)
grafana login

# Or set environment variables for non-interactive use:
export GRAFANA_URL=https://grafana.example.com
export GRAFANA_TOKEN=glsa_xxxxxxxxxxxxxxxxxxxx

# Verify connection
grafana config view
grafana user current -o json
```

### Auth Priority

Configuration is resolved in this order (first match wins):
1. CLI flags (`--url`, `--token`, `--username`, `--password`, `--org-id`)
2. Environment variables (`GRAFANA_URL`, `GRAFANA_TOKEN`, `GRAFANA_USERNAME`, `GRAFANA_PASSWORD`, `GRAFANA_ORG_ID`)
3. Config file profile (`~/.config/grafana-cli/config.yaml`)

### Multiple Profiles

```bash
grafana login                             # saves as default profile
grafana config use-profile prod           # switch profiles
grafana dashboard list --profile staging  # use a specific profile for one command
grafana config list-profiles              # list all profiles
```

## Core Principles for Agents

- **ALWAYS use `-o json`** when you need to parse output programmatically. Pipe through `jq` to extract specific fields.
- **Dashboard UIDs are more stable than numeric IDs** -- prefer UIDs for cross-instance references.
- **Many create/update commands require `-f <file>`** pointing to a JSON or YAML file. Prepare the file first, then pass it. You can also pipe via stdin with `-f -`.
- **For bulk operations:** list with `-o json`, parse with `jq`, then loop over results.
- **For destructive operations** (delete, reset), the CLI requires `--confirm` to skip the interactive confirmation prompt. Always use `--confirm` in scripts.
- **Pagination** is available on list commands via `--page` and `--limit`.
- **Annotation times** (`--from`, `--to`) are in epoch **milliseconds**, not seconds.
- **Always verify connection** before running operations: `grafana config view`.
- **Command aliases:** `dashboard` = `dash` / `db`, `datasource` = `ds`, `service-account` = `sa`, `library-element` = `le`, `contact-point` = `cp`, `mute-timing` = `mt`, `alert template` = `alert tmpl`, `preferences` = `prefs`.

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

# Get full dashboard model by UID
grafana dashboard get <uid> -o json

# Extract just dashboard UIDs from search results
grafana dashboard list -q "production" -o json | jq -r '.[].uid'

# Count panels in a dashboard
grafana dashboard export <uid> | jq '.panels | length'
```

### Export and import dashboards (backup/restore)

```bash
# Export dashboard JSON to stdout
grafana dashboard export <uid>

# Export to a file
grafana dashboard export <uid> --output-file dashboard-backup.json

# Import from file
grafana dashboard import -f dashboard.json

# Import into a specific folder, overwriting if it already exists
grafana dashboard import -f dashboard.json --folder <folder-uid> --overwrite

# Import with a version message
grafana dashboard import -f dashboard.json --overwrite -m "Restored from backup"

# Bulk export all dashboards
grafana dashboard list -o json | jq -r '.[].uid' | while read uid; do
  grafana dashboard export "$uid" --output-file "${uid}.json"
done
```

### Dashboard version history and restore

```bash
# List version history
grafana dashboard versions <uid> -o json

# Restore to a previous version
grafana dashboard restore <uid> <version-number>
```

### Dashboard permissions

```bash
# View permissions
grafana dashboard permissions get <uid> -o json

# Update permissions from file
grafana dashboard permissions update <uid> -f perms.json
```

### Manage datasources

```bash
# List all datasources
grafana datasource list -o json

# List datasources with summary (name, type, uid)
grafana datasource list -o json | jq '.[] | {name, type, uid}'

# Filter by type
grafana datasource list --type prometheus -o json

# Filter by name
grafana datasource list --name "prod" -o json

# Get a specific datasource by UID
grafana datasource get <uid> -o json

# Create a datasource from file
grafana datasource create -f datasource.json

# Update a datasource (uses numeric ID)
grafana datasource update <id> -f updated-ds.json

# Delete a datasource
grafana datasource delete <uid> --confirm
```

### Manage folders

```bash
# List all folders
grafana folder list -o json

# Get folder details
grafana folder get <uid> -o json

# Create a folder
echo '{"title":"My Folder"}' | grafana folder create -f -

# Update a folder
grafana folder update <uid> -f folder.json

# Delete a folder
grafana folder delete <uid> --confirm

# View folder permissions
grafana folder permissions get <uid> -o json

# Update folder permissions
grafana folder permissions update <uid> -f perms.json
```

### Alert rules

```bash
# List all alert rules
grafana alert rule list -o json

# Filter alert rules by folder
grafana alert rule list --folder <folder-uid> -o json

# Filter by rule group
grafana alert rule list --group "High CPU" -o json

# Get a specific rule
grafana alert rule get <uid> -o json

# Create an alert rule from file
grafana alert rule create -f rule.json

# Update an alert rule
grafana alert rule update <uid> -f updated-rule.json

# Delete an alert rule
grafana alert rule delete <uid> --confirm
```

### Contact points

```bash
# List all contact points
grafana alert contact-point list -o json

# Get a specific contact point
grafana alert contact-point get <uid> -o json

# Create a contact point
grafana alert contact-point create -f contact-point.json

# Update a contact point
grafana alert contact-point update <uid> -f updated-cp.json

# Delete a contact point
grafana alert contact-point delete <uid> --confirm
```

### Notification policies

```bash
# Get the notification policy tree
grafana alert policy get -o json

# Update notification policies
grafana alert policy update -f policy.json

# Reset to default policies
grafana alert policy reset --confirm
```

### Alert silences (maintenance windows)

```bash
# List active silences
grafana alert silence list -o json

# Create a silence from a JSON file
cat > silence.json <<'SILENCEOF'
{
  "matchers": [{"name": "alertname", "value": ".*", "isRegex": true}],
  "startsAt": "2024-01-15T00:00:00Z",
  "endsAt": "2024-01-15T06:00:00Z",
  "comment": "Scheduled maintenance window"
}
SILENCEOF
grafana alert silence create -f silence.json

# Create a silence from stdin
echo '{"matchers":[{"name":"severity","value":"warning","isRegex":false}],"startsAt":"2024-01-15T00:00:00Z","endsAt":"2024-01-15T06:00:00Z","comment":"Suppress warnings"}' | grafana alert silence create -f -

# Delete (expire) a silence
grafana alert silence delete <id> --confirm
```

### Mute timings

```bash
# List mute timings
grafana alert mute-timing list -o json

# Get a specific mute timing
grafana alert mute-timing get "weekends" -o json

# Create a mute timing
grafana alert mute-timing create -f mute-timing.json

# Update a mute timing
grafana alert mute-timing update "weekends" -f updated-mt.json

# Delete a mute timing
grafana alert mute-timing delete "weekends" --confirm
```

### Notification templates

```bash
# List notification templates
grafana alert template list -o json

# Get a specific template
grafana alert template get "my-template" -o json

# Create or update a template
grafana alert template update "my-template" -f template.json

# Delete a template
grafana alert template delete "my-template" --confirm
```

### Service accounts and tokens (for CI/CD automation)

```bash
# List service accounts
grafana service-account list -o json

# Create a service account
echo '{"name":"ci-bot","role":"Editor"}' | grafana service-account create -f -

# Get the service account ID from the creation output
SA_ID=$(echo '{"name":"ci-bot","role":"Editor"}' | grafana service-account create -f - -o json | jq -r '.id')

# Create a token for the service account
echo '{"name":"deploy-token","secondsToLive":86400}' | grafana service-account token create "$SA_ID" -f - -o json | jq -r '.key'
# IMPORTANT: The token key is only shown once. Save it immediately.

# List tokens for a service account
grafana service-account token list <sa-id> -o json

# Delete a token
grafana service-account token delete <sa-id> <token-id> --confirm

# Delete a service account
grafana service-account delete <sa-id> --confirm
```

### Annotations (mark deployments, incidents)

```bash
# List annotations (default limit 100)
grafana annotation list -o json

# Filter by tags
grafana annotation list --tags deploy,release -o json

# Filter by dashboard and time range (epoch milliseconds)
grafana annotation list --dashboard-id 42 --from 1705312800000 --to 1705399200000 -o json

# Filter by type (annotation or alert)
grafana annotation list --type alert -o json

# Increase result limit
grafana annotation list --limit 500 -o json

# Create an annotation to mark a deployment
echo '{"text":"Deployed v2.1.0","tags":["deploy","production"]}' | grafana annotation create -f -

# Get a specific annotation
grafana annotation get <id> -o json

# Update an annotation
grafana annotation update <id> -f annotation.json

# Delete an annotation
grafana annotation delete <id> --confirm

# List all unique annotation tags
grafana annotation tags -o json
```

### User and org management

```bash
# Show current authenticated user
grafana user current -o json

# List all users
grafana user list -o json

# Search users
grafana user list -q "john" -o json

# Get user details
grafana user get <id> -o json

# Look up user by login or email
grafana user lookup admin
grafana user lookup admin@example.com

# List a user's orgs
grafana user orgs <user-id> -o json

# List a user's teams
grafana user teams <user-id> -o json

# Update a user
grafana user update <id> -f user.json

# Show current org
grafana org current -o json

# List all organizations
grafana org list -o json

# Switch to a different org
grafana org switch <org-id>

# List users in an organization
grafana org user list <org-id> -o json

# Filter org users by role
grafana org user list <org-id> --role Admin -o json

# Search org users
grafana org user list <org-id> -q "dev" -o json

# Add a user to an org
echo '{"loginOrEmail":"newuser@example.com","role":"Viewer"}' | grafana org user add <org-id> -f -

# Update a user's role in an org
echo '{"role":"Editor"}' | grafana org user update <org-id> <user-id> -f -

# Remove a user from an org
grafana org user remove <org-id> <user-id> --confirm
```

### Team management

```bash
# List teams
grafana team list -o json

# Search teams
grafana team list -q "backend" -o json

# Get team details
grafana team get <id> -o json

# Create a team
echo '{"name":"Backend Team","email":"backend@example.com"}' | grafana team create -f -

# Update a team
grafana team update <id> -f team.json

# Delete a team
grafana team delete <id> --confirm

# List team members
grafana team member list <team-id> -o json

# Add a member to a team
grafana team member add <team-id> <user-id>

# Remove a member from a team
grafana team member remove <team-id> <user-id> --confirm

# View team preferences
grafana team preferences get <team-id> -o json

# Update team preferences
grafana team preferences update <team-id> -f prefs.json
```

### Snapshots

```bash
# List snapshots
grafana snapshot list -o json

# Get a snapshot by key
grafana snapshot get <key> -o json

# Create a snapshot
grafana snapshot create -f snapshot.json

# Delete a snapshot
grafana snapshot delete <key> --confirm
```

### Playlists

```bash
# List playlists
grafana playlist list -o json

# Search playlists
grafana playlist list -q "production" -o json

# Get a playlist
grafana playlist get <uid> -o json

# Create a playlist
grafana playlist create -f playlist.json

# Update a playlist
grafana playlist update <uid> -f playlist.json

# Delete a playlist
grafana playlist delete <uid> --confirm
```

### Library elements (reusable panels and variables)

```bash
# List all library elements
grafana library-element list -o json

# List panels only (kind=1)
grafana library-element list --kind 1 -o json

# List variables only (kind=2)
grafana library-element list --kind 2 -o json

# Search by name
grafana library-element list --search "CPU" -o json

# Filter by folder
grafana library-element list --folder "General" -o json

# Get a library element
grafana library-element get <uid> -o json

# Create a library element
grafana library-element create -f panel.json

# Update a library element
grafana library-element update <uid> -f updated-panel.json

# Delete a library element
grafana library-element delete <uid> --confirm

# List dashboards using a library element
grafana library-element connections <uid> -o json
```

### Correlations

```bash
# List all correlations
grafana correlation list -o json

# Get a specific correlation
grafana correlation get <source-uid> <correlation-uid> -o json

# Create a correlation
grafana correlation create <source-uid> -f correlation.json

# Update a correlation
grafana correlation update <source-uid> <correlation-uid> -f correlation.json

# Delete a correlation
grafana correlation delete <source-uid> <correlation-uid> --confirm
```

### User preferences

```bash
# View current preferences
grafana preferences get -o json

# Update preferences
echo '{"theme":"dark","timezone":"utc","weekStart":"monday"}' | grafana preferences update -f -
```

### Server administration (requires admin permissions)

```bash
# View server settings
grafana admin settings -o json

# View server usage statistics
grafana admin stats -o json

# Reload provisioning configurations
grafana admin reload dashboards
grafana admin reload datasources
grafana admin reload plugins
grafana admin reload access-control
grafana admin reload alerting
```

### Star/unstar dashboards

```bash
grafana user star add <dashboard-id>
grafana user star remove <dashboard-id>
```

## Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `401 Unauthorized` | Invalid or expired token | Check token with `grafana config view`; re-run `grafana login` |
| `403 Forbidden` | Token lacks required permissions | Use a token with Admin role for admin operations |
| `404 Not Found` | Wrong UID, wrong org context, or resource does not exist | Verify UID; check org with `grafana org current` |
| Connection refused / timeout | Wrong URL or server down | Verify URL with `grafana config view`; use `--url` to override |
| TLS certificate errors | Self-signed certificates | Some environments may need `--insecure` flag (if supported) or proper CA configuration |
| `409 Conflict` | Resource already exists | Use `--overwrite` on dashboard import/create |
| Empty results | Wrong org context or filters too narrow | Check `grafana org current`; broaden search filters |

## Output Formats

All list and get commands support three output formats via `-o`:

- `-o table` (default) -- human-readable tabular output
- `-o json` -- JSON, ideal for programmatic parsing with `jq`
- `-o yaml` -- YAML, useful for config management

## References

For the complete list of every command with all flags, see [`references/commands.md`](references/commands.md) in this skill directory.

For full documentation, read the [README.md](../README.md) at the root of the grafana-cli repository.
