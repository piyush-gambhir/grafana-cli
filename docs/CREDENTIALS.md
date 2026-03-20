# Grafana CLI - Authentication & Credentials Guide

This guide covers every authentication method supported by the Grafana CLI, how to obtain credentials from different Grafana environments, and how to troubleshoot common issues.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Getting Your Credentials](#getting-your-credentials)
- [Minimum Required Permissions](#minimum-required-permissions)
- [Configuration](#configuration)
- [TLS / SSL Configuration](#tls--ssl-configuration)
- [Edge Cases & Troubleshooting](#edge-cases--troubleshooting)
- [Security Best Practices](#security-best-practices)

---

## Quick Start

### Option 1: Interactive Login (Recommended)

The fastest way to get started. Run `grafana login` and follow the prompts:

```bash
grafana login
```

You will be prompted for:

```
Grafana URL: https://grafana.example.com
Auth method (token/basic) [token]: token
API Token: glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Testing connection...
Connection successful! Org: Main Org. (ID: 1)
Profile name [default]: default
Profile "default" saved and set as current.
```

After login, credentials are saved to `~/.config/grafana-cli/config.yaml` and all subsequent commands will use them automatically.

### Option 2: Environment Variables

For non-interactive use, CI/CD pipelines, and scripting, set environment variables instead:

```bash
export GRAFANA_URL=https://grafana.example.com
export GRAFANA_TOKEN=glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Now run any command -- no login needed
grafana dashboard list
```

### Option 3: CLI Flags (One-Off Commands)

Pass credentials directly on any command for one-off use:

```bash
grafana dashboard list \
  --url https://grafana.example.com \
  --token glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Or with basic auth:

```bash
grafana dashboard list \
  --url https://grafana.example.com \
  --username admin \
  --password admin
```

---

## Getting Your Credentials

### Service Account Token (Recommended)

Service account tokens are the recommended authentication method. They are not tied to a specific user, can have scoped permissions, and work regardless of your identity provider configuration.

#### Self-Hosted Grafana

1. **Log in to Grafana** as an administrator (`http://your-grafana:3000`).

2. **Navigate to Service Accounts:**
   - Open the hamburger menu (top-left).
   - Go to **Administration** > **Users and access** > **Service Accounts**.
   - Or navigate directly to `http://your-grafana:3000/org/serviceaccounts`.

3. **Create a new service account:**
   - Click **"Add service account"**.
   - Enter a descriptive name (e.g., `grafana-cli`, `ci-deploy-bot`, `monitoring-readonly`).
   - Select a role:
     - **Viewer** -- read-only access to dashboards, datasources, folders, annotations.
     - **Editor** -- create/update dashboards, annotations, playlists, library elements.
     - **Admin** -- full access including user/org/team management and server admin.
   - Click **"Create"**.

4. **Generate a token:**
   - On the service account detail page, click **"Add service account token"**.
   - Give the token a name (e.g., `cli-token`).
   - Optionally set an expiration date (leave blank for no expiration).
   - Click **"Generate token"**.

5. **Copy the token immediately.** It starts with `glsa_` and is only shown once. Example:
   ```
   glsa_abc123def456ghi789jkl012mno345pqr678stu901_vwxyz12345
   ```

6. **Use the token:**
   ```bash
   # Interactive login
   grafana login
   # Enter URL, choose "token", paste the glsa_ token

   # Or via environment variable
   export GRAFANA_URL=http://your-grafana:3000
   export GRAFANA_TOKEN=glsa_abc123def456ghi789jkl012mno345pqr678stu901_vwxyz12345
   ```

#### Grafana Cloud

1. **Log in to Grafana Cloud** at [grafana.com](https://grafana.com) and open your stack.

2. **Open your Grafana instance** -- click the "Launch" button on your stack's Grafana tile. Your instance URL will be in the format:
   ```
   https://<your-stack-name>.grafana.net
   ```

3. **Create a service account** using the same steps as self-hosted above:
   - Navigate to **Administration** > **Users and access** > **Service Accounts**.
   - Create a service account and generate a token.

4. **Use the token with your cloud URL:**
   ```bash
   export GRAFANA_URL=https://mycompany.grafana.net
   export GRAFANA_TOKEN=glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

   grafana dashboard list
   ```

> **Note:** Grafana Cloud also supports Cloud API keys at the organization level (via grafana.com portal > API Keys). These are different from Grafana instance service account tokens. The CLI requires an instance-level service account token, not a Cloud API key.

#### Amazon Managed Grafana

Amazon Managed Grafana (AMG) uses a different token format and setup process.

1. **Open the AWS Console** and navigate to **Amazon Managed Grafana**.

2. **Select your workspace** and note the workspace URL:
   ```
   https://g-xxxxxxxxxx.grafana-workspace.<region>.amazonaws.com
   ```

3. **Create a Grafana API key** (AMG uses API keys, not service account tokens):
   - In the AWS Console, select your workspace.
   - Go to the **Authentication** tab.
   - Under **API keys**, click **Create API key**.
   - Set the key name, role (Viewer/Editor/Admin), and time-to-live.
   - Copy the generated key.

4. **Alternatively, use the Grafana UI** if your workspace is on Grafana 9+:
   - Open your workspace URL in a browser.
   - Navigate to **Administration** > **Users and access** > **Service Accounts**.
   - Create a service account and token as with self-hosted Grafana.

5. **Use the key with the CLI:**
   ```bash
   export GRAFANA_URL=https://g-xxxxxxxxxx.grafana-workspace.us-east-1.amazonaws.com
   export GRAFANA_TOKEN=<your-api-key-or-sa-token>

   grafana dashboard list
   ```

> **Note:** IAM-authenticated access (AWS SigV4 signing) is not currently supported by this CLI. You must use an API key or service account token.

#### Azure Managed Grafana

1. **Open the Azure Portal** and navigate to your Azure Managed Grafana resource.

2. **Note your instance URL:**
   ```
   https://<instance-name>-<random>.grafana.azure.com
   ```

3. **Create a service account** via the Grafana UI:
   - Open your Grafana instance URL in a browser.
   - Navigate to **Administration** > **Users and access** > **Service Accounts**.
   - Create a service account and generate a token.

4. **Use the token:**
   ```bash
   export GRAFANA_URL=https://myinstance-abc123.grafana.azure.com
   export GRAFANA_TOKEN=glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

   grafana dashboard list
   ```

### Basic Auth (Username/Password)

Basic authentication uses a Grafana username and password directly. This method is simpler but less secure than service account tokens.

**When to use:**
- Local development instances with default admin credentials.
- Quick testing against a development Grafana.
- When service accounts are not available (Grafana < 8.0).

**How to set up:**

```bash
# Interactive login
grafana login
# Choose "basic" for auth method
# Enter username and password

# Or via environment variables
export GRAFANA_URL=http://localhost:3000
export GRAFANA_USERNAME=admin
export GRAFANA_PASSWORD=admin

# Or via CLI flags
grafana dashboard list --url http://localhost:3000 --username admin --password admin
```

**Limitations:**
- If your Grafana instance uses an external identity provider (OAuth, SAML, LDAP) and has disabled built-in auth, basic auth with Grafana credentials will not work. Use a service account token instead.
- Basic auth credentials are stored in plain text in the config file.
- Cannot scope permissions -- the user's full permissions are used.
- If the user's password is changed or the account is disabled, the CLI stops working.

### API Keys (Legacy/Deprecated)

> **Deprecated:** API keys have been deprecated since Grafana 9.0 in favor of service account tokens. They will be removed in a future Grafana release.

If you are running Grafana 8.x or earlier, or have existing API keys:

1. Navigate to **Configuration** > **API keys** in the Grafana UI.
2. Click **"Add API key"**.
3. Set the name, role, and optional time-to-live.
4. Copy the generated key (starts with `eyJ`).

API keys are used exactly like service account tokens with the CLI:

```bash
export GRAFANA_URL=http://grafana.example.com
export GRAFANA_TOKEN=eyJrIjoiT0tTcG1pUlY2...

grafana dashboard list
```

**Migration:** You can convert existing API keys to service account tokens in the Grafana UI at **Administration** > **Users and access** > **Service Accounts** > **"Migrate API keys"**.

---

## Minimum Required Permissions

The table below shows the minimum Grafana role required for each CLI operation. Assigning a lower-privilege role than required will result in `403 Forbidden` errors.

### Dashboard Operations

| Operation | Minimum Role | Notes |
|---|---|---|
| `grafana dashboard list` | Viewer | Only returns dashboards the user can see |
| `grafana dashboard get <uid>` | Viewer | |
| `grafana dashboard export <uid>` | Viewer | |
| `grafana dashboard versions <uid>` | Viewer | |
| `grafana dashboard create` | Editor | |
| `grafana dashboard update` | Editor | Must have edit permission on the target dashboard |
| `grafana dashboard import` | Editor | |
| `grafana dashboard restore <uid> <ver>` | Editor | |
| `grafana dashboard delete <uid>` | Editor | Must have edit permission on the target dashboard |
| `grafana dashboard permissions get <uid>` | Viewer | |
| `grafana dashboard permissions update <uid>` | Admin | |

### Datasource Operations

| Operation | Minimum Role | Notes |
|---|---|---|
| `grafana datasource list` | Viewer | |
| `grafana datasource get <uid>` | Viewer | |
| `grafana datasource create` | Admin | Server admin or org admin |
| `grafana datasource update <id>` | Admin | |
| `grafana datasource delete <uid>` | Admin | |

### Folder Operations

| Operation | Minimum Role | Notes |
|---|---|---|
| `grafana folder list` | Viewer | Only returns folders the user can see |
| `grafana folder get <uid>` | Viewer | |
| `grafana folder create` | Editor | |
| `grafana folder update <uid>` | Editor | Must have edit permission on the folder |
| `grafana folder delete <uid>` | Editor | Must have edit permission on the folder |
| `grafana folder permissions get <uid>` | Admin | |
| `grafana folder permissions update <uid>` | Admin | |

### Alert Operations

| Operation | Minimum Role | Notes |
|---|---|---|
| `grafana alert rule list` | Viewer | |
| `grafana alert rule get <uid>` | Viewer | |
| `grafana alert rule create` | Editor | |
| `grafana alert rule update <uid>` | Editor | |
| `grafana alert rule delete <uid>` | Editor | |
| `grafana alert contact-point list` | Editor | |
| `grafana alert contact-point create` | Editor | |
| `grafana alert policy get` | Viewer | |
| `grafana alert policy update` | Editor | |
| `grafana alert silence list` | Viewer | |
| `grafana alert silence create` | Editor | |

### User & Organization Operations

| Operation | Minimum Role | Notes |
|---|---|---|
| `grafana user current` | Viewer | Any authenticated user |
| `grafana user list` | Admin | Server admin required |
| `grafana user get <id>` | Admin | Server admin required |
| `grafana user lookup` | Admin | Server admin required |
| `grafana user update <id>` | Admin | Server admin required |
| `grafana org list` | Admin | Server admin required |
| `grafana org get <id>` | Admin | Server admin required |
| `grafana org create` | Admin | Server admin required |
| `grafana org current` | Viewer | Any authenticated user |
| `grafana org switch <id>` | Viewer | User must be a member of the target org |
| `grafana org user list <org-id>` | Admin | Org admin or server admin |
| `grafana team list` | Admin | |
| `grafana team get <id>` | Admin | |
| `grafana team create` | Admin | |
| `grafana team member list <id>` | Admin | |

### Service Account Operations

| Operation | Minimum Role | Notes |
|---|---|---|
| `grafana service-account list` | Admin | Org admin or server admin |
| `grafana service-account get <id>` | Admin | |
| `grafana service-account create` | Admin | |
| `grafana service-account delete <id>` | Admin | |
| `grafana service-account token list <id>` | Admin | |
| `grafana service-account token create <id>` | Admin | |

### Other Operations

| Operation | Minimum Role | Notes |
|---|---|---|
| `grafana annotation list` | Viewer | |
| `grafana annotation create` | Editor | |
| `grafana annotation update <id>` | Editor | |
| `grafana annotation delete <id>` | Editor | |
| `grafana snapshot list` | Viewer | |
| `grafana snapshot create` | Editor | |
| `grafana playlist list` | Viewer | |
| `grafana playlist create` | Editor | |
| `grafana library-element list` | Viewer | |
| `grafana library-element create` | Editor | |
| `grafana admin settings` | Admin | Server admin required |
| `grafana admin stats` | Admin | Server admin required |
| `grafana admin reload <resource>` | Admin | Server admin required |
| `grafana preferences get` | Viewer | |
| `grafana preferences update` | Viewer | Any authenticated user (own prefs) |

---

## Configuration

### Config File

The CLI stores connection profiles in a YAML configuration file at:

```
~/.config/grafana-cli/config.yaml
```

If the `XDG_CONFIG_HOME` environment variable is set, the file is stored at:

```
$XDG_CONFIG_HOME/grafana-cli/config.yaml
```

Here is a complete example with all supported fields:

```yaml
# The currently active profile
current_profile: production

# Default settings
defaults:
  # Default output format for all commands: table, json, or yaml
  output: table

# Named connection profiles
profiles:
  # Production instance using a service account token
  production:
    url: https://grafana.prod.example.com
    token: glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    org_id: 1

  # Staging instance using basic auth
  staging:
    url: https://grafana.staging.example.com
    username: admin
    password: staging-password
    org_id: 2

  # Local development instance
  local:
    url: http://localhost:3000
    token: glsa_yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy

  # Read-only profile for safety
  prod-readonly:
    url: https://grafana.prod.example.com
    token: glsa_zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz
    read_only: true
```

**Field reference:**

| Field | Type | Description |
|---|---|---|
| `current_profile` | string | Name of the active profile |
| `defaults.output` | string | Default output format: `table`, `json`, or `yaml` |
| `profiles.<name>.url` | string | Grafana instance URL |
| `profiles.<name>.token` | string | Service account token or API key |
| `profiles.<name>.username` | string | Username for basic auth |
| `profiles.<name>.password` | string | Password for basic auth |
| `profiles.<name>.org_id` | integer | Grafana organization ID |
| `profiles.<name>.read_only` | boolean | If `true`, block write/delete operations |

> **File permissions:** The config file is written with mode `0600` (owner read/write only) to protect credentials.

### Environment Variables

All connection settings can be provided via environment variables. This is the recommended approach for CI/CD and automation.

| Environment Variable | Description | Example |
|---|---|---|
| `GRAFANA_URL` | Grafana instance URL | `https://grafana.example.com` |
| `GRAFANA_TOKEN` | Service account token or API key | `glsa_xxxxxxxxxxxxxxxxxxxx` |
| `GRAFANA_USERNAME` | Username for basic auth | `admin` |
| `GRAFANA_PASSWORD` | Password for basic auth | `admin` |
| `GRAFANA_ORG_ID` | Organization ID | `1` |
| `GRAFANA_READ_ONLY` | Enable read-only mode | `true` or `1` |
| `XDG_CONFIG_HOME` | Override config directory base path | `/home/user/.config` |

**Usage in a CI pipeline (GitHub Actions example):**

```yaml
jobs:
  deploy-dashboards:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Grafana CLI
        run: curl -sSfL https://raw.githubusercontent.com/piyush-gambhir/grafana-cli/main/install.sh | sh

      - name: Deploy dashboard
        env:
          GRAFANA_URL: ${{ secrets.GRAFANA_URL }}
          GRAFANA_TOKEN: ${{ secrets.GRAFANA_TOKEN }}
        run: |
          grafana dashboard create -f dashboards/production.json --overwrite -m "Deploy from CI"
```

**Usage in a shell script:**

```bash
#!/bin/bash
set -euo pipefail

export GRAFANA_URL="https://grafana.prod.example.com"
export GRAFANA_TOKEN="glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# Export all dashboards for backup
for uid in $(grafana dashboard list -o json | jq -r '.[].uid'); do
  grafana dashboard export "$uid" --output-file "backups/${uid}.json"
done
```

### Config Priority

Configuration is resolved using a layered approach. The first non-empty value wins:

```
1. CLI flags          (highest priority)
   --url, --token, --username, --password, --org-id

2. Environment variables
   GRAFANA_URL, GRAFANA_TOKEN, GRAFANA_USERNAME, GRAFANA_PASSWORD, GRAFANA_ORG_ID

3. Config file profile
   ~/.config/grafana-cli/config.yaml (active profile)

4. Built-in defaults  (lowest priority)
   output: table
```

This means you can:
- Set a base configuration via `grafana login` (saved in the config file).
- Override specific values per-session with environment variables.
- Override on a per-command basis with CLI flags.

**Example -- override the org for a single command:**

```bash
# Config file has org_id: 1 for the "production" profile.
# Run this one command against org 3 instead:
grafana dashboard list --org-id 3
```

**Example -- override the URL via environment for a script:**

```bash
# Profile points to prod, but temporarily target staging:
GRAFANA_URL=https://staging.example.com grafana dashboard list
```

### Multiple Profiles

Profiles let you manage credentials for multiple Grafana instances and switch between them.

**Create profiles via interactive login:**

```bash
# First profile (saved as "default")
grafana login
# Enter prod URL and token, name it "production"

# Second profile
grafana login
# Enter staging URL and token, name it "staging"

# Third profile
grafana login
# Enter local URL and credentials, name it "local"
```

**List all profiles:**

```bash
grafana config list-profiles
```

Output:

```
* production (https://grafana.prod.example.com, token)
  staging (https://grafana.staging.example.com, token)
  local (http://localhost:3000, basic)
```

The `*` marks the currently active profile.

**Switch profiles:**

```bash
# Set the active profile
grafana config use-profile staging

# Or use a profile for a single command without switching
grafana dashboard list --profile local
```

**Use different profiles in scripts:**

```bash
# Back up dashboards from production, deploy to staging
for uid in $(grafana dashboard list --profile production -o json | jq -r '.[].uid'); do
  grafana dashboard export "$uid" --profile production --output-file "/tmp/${uid}.json"
  grafana dashboard import -f "/tmp/${uid}.json" --profile staging --overwrite
done
```

---

## TLS / SSL Configuration

### Self-Signed Certificates

If your Grafana instance uses a self-signed certificate, you will see errors like:

```
Error: sending request: Get "https://grafana.internal:3000/api/org/":
  x509: certificate signed by unknown authority
```

**Workaround -- add the CA to your system trust store:**

On **macOS:**

```bash
# Add the CA certificate to the system keychain
sudo security add-trusted-cert -d -r trustRoot \
  -k /Library/Keychains/System.keychain /path/to/ca-cert.pem
```

On **Linux (Debian/Ubuntu):**

```bash
sudo cp /path/to/ca-cert.pem /usr/local/share/ca-certificates/grafana-ca.crt
sudo update-ca-certificates
```

On **Linux (RHEL/CentOS/Fedora):**

```bash
sudo cp /path/to/ca-cert.pem /etc/pki/ca-trust/source/anchors/grafana-ca.pem
sudo update-ca-trust
```

On **Windows:**

```powershell
Import-Certificate -FilePath "C:\path\to\ca-cert.pem" -CertStoreLocation Cert:\LocalMachine\Root
```

After adding the CA, the CLI will trust the certificate without any additional configuration.

**Alternative -- set `SSL_CERT_FILE` or `SSL_CERT_DIR`:**

Go programs (including this CLI) respect these environment variables:

```bash
export SSL_CERT_FILE=/path/to/ca-bundle.pem
grafana dashboard list
```

### Grafana Behind a Reverse Proxy

When Grafana runs behind a reverse proxy (nginx, Apache, Caddy, etc.), the CLI connects to the proxy URL. No special configuration is needed as long as the proxy forwards API requests correctly.

**nginx example:**

```nginx
server {
    listen 443 ssl;
    server_name grafana.example.com;

    ssl_certificate     /etc/ssl/certs/grafana.crt;
    ssl_certificate_key /etc/ssl/private/grafana.key;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Use the external proxy URL with the CLI:

```bash
export GRAFANA_URL=https://grafana.example.com
export GRAFANA_TOKEN=glsa_xxxxxxxxxxxxxxxxxxxx
grafana dashboard list
```

**Apache example:**

```apache
<VirtualHost *:443>
    ServerName grafana.example.com

    SSLEngine on
    SSLCertificateFile    /etc/ssl/certs/grafana.crt
    SSLCertificateKeyFile /etc/ssl/private/grafana.key

    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:3000/
    ProxyPassReverse / http://127.0.0.1:3000/
</VirtualHost>
```

**Grafana with a subpath:**

If Grafana is served from a subpath (e.g., `https://example.com/grafana/`), include the subpath in the URL:

```bash
export GRAFANA_URL=https://example.com/grafana
grafana dashboard list
```

Make sure Grafana's `server.root_url` setting matches:

```ini
[server]
root_url = https://example.com/grafana/
serve_from_sub_path = true
```

### HTTP Proxy Support

The CLI uses Go's standard HTTP client, which respects the `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY` environment variables:

```bash
export HTTPS_PROXY=http://proxy.corp.example.com:8080
export NO_PROXY=localhost,127.0.0.1

grafana dashboard list --url https://grafana.example.com --token glsa_xxxx
```

---

## Edge Cases & Troubleshooting

### Grafana with OAuth / SAML / LDAP

When Grafana is configured with an external identity provider (Google OAuth, Okta, SAML, LDAP, etc.):

- **Service account tokens always work.** Service accounts are internal to Grafana and bypass the external IdP entirely. This is the recommended approach.

- **Basic auth may not work** if Grafana has disabled built-in authentication. Check your Grafana configuration:
  ```ini
  # grafana.ini
  [auth]
  disable_login_form = true  # If true, basic auth will NOT work
  ```

- **OAuth tokens cannot be used** with this CLI. The CLI does not perform OAuth flows. Use service account tokens instead.

**Bottom line:** If your organization uses SSO, create a service account with a token. Do not try to use your SSO credentials with the CLI.

### Organization-Specific Access

Grafana supports multiple organizations. By default, API requests go to the user's current organization. To target a specific org:

**Using the `--org-id` flag:**

```bash
# List dashboards in org 3
grafana dashboard list --org-id 3

# Create a dashboard in org 2
grafana dashboard create -f dashboard.json --org-id 2
```

**Using the `GRAFANA_ORG_ID` environment variable:**

```bash
export GRAFANA_ORG_ID=3
grafana dashboard list
```

**Saving org ID in a profile:**

```bash
# The org_id is saved during login or can be added to config.yaml manually
grafana login
# Set org_id in your profile:
```

```yaml
profiles:
  org-3:
    url: https://grafana.example.com
    token: glsa_xxxxxxxxxxxxxxxxxxxx
    org_id: 3
```

**When you need org-specific tokens:**

Service accounts belong to a specific organization. A token created in org 1 has permissions only within org 1. To manage multiple orgs, create a service account in each org and set up separate profiles:

```yaml
profiles:
  prod-org1:
    url: https://grafana.example.com
    token: glsa_token_for_org1
    org_id: 1
  prod-org2:
    url: https://grafana.example.com
    token: glsa_token_for_org2
    org_id: 2
```

> **Exception:** A Grafana Server Admin service account can access all organizations by setting the `X-Grafana-Org-Id` header (which the `--org-id` flag does automatically).

### Grafana Behind a VPN

If your Grafana instance is only accessible through a VPN:

1. **Ensure VPN connectivity** before running CLI commands.
2. **Check DNS resolution** -- `nslookup grafana.internal.example.com` should resolve.
3. **Test basic connectivity:**
   ```bash
   curl -s -o /dev/null -w "%{http_code}" https://grafana.internal.example.com/api/health
   # Should return 200
   ```

**Common error messages when Grafana is unreachable:**

| Error Message | Meaning |
|---|---|
| `dial tcp: lookup grafana.example.com: no such host` | DNS cannot resolve the hostname. Check VPN connection or DNS settings. |
| `dial tcp 10.0.1.50:3000: i/o timeout` | Network path exists but connection timed out. Check VPN tunnel, firewall rules. |
| `dial tcp 10.0.1.50:3000: connect: connection refused` | Host is reachable but nothing is listening on that port. Verify Grafana is running and the port is correct. |
| `EOF` | Connection was established but closed unexpectedly. Could indicate a TLS issue or a load balancer dropping the connection. |

### Token Expiration

- **Service account tokens do not expire by default.** When you create a token without setting `secondsToLive`, it is valid indefinitely.

- **Administrators can enforce token expiration** via the Grafana configuration:
  ```ini
  # grafana.ini
  [service_accounts]
  token_expiration_day_limit = 90  # Tokens expire after 90 days max
  ```

- **API keys (legacy) can have a TTL** set at creation time. Once expired, they stop working.

**How to check if your token has expired:**

```bash
grafana user current
```

If you see a `401 Unauthorized` error, your token may be expired.

**How to rotate a token:**

1. Create a new token for the existing service account:
   ```bash
   grafana service-account token create <sa-id> -f token.json
   ```
   Where `token.json` contains:
   ```json
   {"name": "cli-token-2025-03", "secondsToLive": 7776000}
   ```

2. Update your config or environment variable with the new token.

3. Delete the old token:
   ```bash
   grafana service-account token delete <sa-id> <old-token-id> --confirm
   ```

### Common Errors

| Error | Cause | Fix |
|---|---|---|
| `grafana URL is required` | No URL configured via flag, env var, or profile | Set `GRAFANA_URL`, use `--url`, or run `grafana login` |
| `API error: status 401` | Token is invalid, expired, or malformed | Regenerate a new token and update your config |
| `API error: status 403` | Token is valid but lacks required permissions | Assign a higher role to the service account (Viewer -> Editor -> Admin) |
| `API error: status 404` | Endpoint not found; often wrong URL or Grafana version | Verify the URL includes any subpath; check Grafana version supports the API |
| `API error: status 412` | Precondition failed (e.g., dashboard version conflict) | Fetch the latest version and retry |
| `connection refused` | Grafana is not running, wrong port, or firewall block | Verify Grafana service is running: `systemctl status grafana-server` |
| `no such host` | DNS lookup failed | Check hostname spelling, VPN status, DNS configuration |
| `i/o timeout` | Network is unreachable or too slow | Check VPN, firewall rules, network routing |
| `x509: certificate signed by unknown authority` | Self-signed or untrusted TLS certificate | Add the CA to your trust store (see [TLS section](#self-signed-certificates)) |
| `x509: certificate is valid for X, not Y` | URL hostname does not match the TLS certificate | Use the hostname that matches the certificate's Subject Alternative Name |
| `invalid auth method: X` | Entered something other than "token" or "basic" during login | Re-run `grafana login` and enter `token` or `basic` |
| `profile "X" not found` | Referenced a profile that does not exist | Run `grafana config list-profiles` to see available profiles |
| `profile "X" already exists` | Tried to create a duplicate profile | The CLI overwrites existing profiles during `grafana login`, so this is typically only encountered programmatically |

### Read-Only Mode

The CLI supports a read-only mode that prevents accidental writes, updates, and deletes. This is useful for production profiles where you only need to query data.

**Enable read-only mode in a profile:**

```yaml
profiles:
  prod-readonly:
    url: https://grafana.prod.example.com
    token: glsa_xxxxxxxxxxxxxxxxxxxx
    read_only: true
```

**Enable via environment variable:**

```bash
export GRAFANA_READ_ONLY=true
grafana dashboard delete abc123  # Will be blocked
```

---

## Security Best Practices

### Use Service Account Tokens, Not API Keys

API keys are deprecated. Service account tokens provide better auditing (tied to a named service account), can be managed per-account, and will be supported long-term.

### Use Least-Privilege Roles

Assign the minimum role required for your use case:

- **Viewer** for monitoring scripts that only read dashboards.
- **Editor** for CI/CD that deploys dashboards and alert rules.
- **Admin** only for administrative automation (user management, datasource provisioning).

### Store Tokens Securely

**Do:**
- Use environment variables from a secrets manager in CI/CD (GitHub Secrets, AWS Secrets Manager, Vault).
- Use `grafana login` for interactive use (config file is stored with `0600` permissions).
- Use `.env` files that are excluded from version control via `.gitignore`.

**Do not:**
- Hard-code tokens in scripts or source code.
- Commit the config file (`~/.config/grafana-cli/config.yaml`) to version control.
- Pass tokens as command-line arguments in shared systems (they appear in process listings and shell history).

**Example with a `.env` file:**

```bash
# .env (add to .gitignore!)
GRAFANA_URL=https://grafana.prod.example.com
GRAFANA_TOKEN=glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

```bash
# Load and use
source .env
grafana dashboard list
```

### Rotate Tokens Periodically

Even though service account tokens do not expire by default, rotate them on a regular schedule:

1. Create a new token.
2. Update all systems using the old token.
3. Delete the old token.
4. Repeat every 90 days (or per your organization's policy).

### Use Read-Only Mode for Safety

When you only need to query data (especially in production), enable `read_only: true` in your profile or set `GRAFANA_READ_ONLY=true`. This prevents accidental mutations.

### Audit Service Account Usage

Periodically review your service accounts:

```bash
# List all service accounts
grafana service-account list -o json

# Check tokens for each service account
grafana service-account token list <sa-id> -o json
```

Remove unused service accounts and tokens to reduce your attack surface.

### Separate Profiles Per Environment

Maintain separate profiles for each environment (production, staging, development) to avoid accidentally running commands against the wrong instance:

```bash
# Always explicit about which environment you target
grafana dashboard list --profile staging
grafana dashboard list --profile production

# Set your default to the safest environment
grafana config use-profile local
```
