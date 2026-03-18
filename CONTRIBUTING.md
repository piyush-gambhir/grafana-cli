# Contributing to Grafana CLI

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

- Go 1.22 or later
- Make
- Git

### Clone and Build

```bash
git clone https://github.com/piyush-gambhir/grafana-cli.git
cd grafana-cli
make build
```

### Run Locally

```bash
./bin/grafana --help
./bin/grafana version
```

### Run Tests

```bash
make test
```

### Lint

```bash
make lint    # requires golangci-lint
make vet     # go vet
make fmt     # gofmt
```

## Project Structure

```
.
├── main.go                 # Entry point
├── cmd/                    # Cobra command definitions
│   ├── root.go             # Root command, global flags
│   ├── login.go            # Auth commands
│   ├── dashboard/          # Dashboard CRUD commands
│   │   ├── list.go
│   │   ├── get.go
│   │   ├── create.go
│   │   └── ...
│   ├── datasource/         # Datasource commands
│   ├── folder/             # Folder commands
│   ├── alert/              # Unified Alerting (rules, contact points, policies, silences, mute timings, templates)
│   ├── org/                # Organization and org user commands
│   ├── team/               # Team, member, and preferences commands
│   ├── user/               # User commands
│   ├── serviceaccount/     # Service account and token commands
│   ├── annotation/         # Annotation and tag commands
│   ├── snapshot/           # Dashboard snapshot commands
│   ├── playlist/           # Playlist commands
│   ├── libraryelement/     # Library panel/variable commands
│   ├── correlation/        # Datasource correlation commands
│   ├── admin/              # Server admin commands
│   ├── preferences/        # User preferences commands
│   └── config/             # CLI config management
├── internal/
│   ├── client/             # HTTP API client
│   │   ├── client.go       # Base client (auth, headers, errors)
│   │   ├── dashboards.go   # Dashboard API methods
│   │   ├── datasources.go  # Datasource API methods
│   │   ├── alerting.go     # Unified Alerting API methods
│   │   ├── annotations.go  # Annotation API methods
│   │   ├── folders.go      # Folder API methods
│   │   ├── orgs.go         # Organization API methods
│   │   └── ...
│   ├── cmdutil/            # Shared command utilities (Factory, flag helpers)
│   ├── config/             # Config file and auth resolution
│   ├── output/             # JSON/YAML/Table formatters
│   ├── build/              # Build version info
│   └── update/             # Self-update logic
├── Makefile
├── .goreleaser.yaml
└── .github/workflows/
    ├── ci.yml              # Build + test on every push/PR
    └── release.yml         # GoReleaser on tag push
```

## Adding a New Command

1. **Add the API method** in `internal/client/<resource>.go`:
   ```go
   func (c *Client) ListWidgets(params ...) ([]Widget, error) {
       // HTTP call to the Grafana API
   }
   ```

2. **Create the command** in `cmd/<resource>/list.go`:
   ```go
   func NewListCmd(f *cmdutil.Factory) *cobra.Command {
       // Define flags, run function, help text with examples
   }
   ```

3. **Register** the command in the parent command's `New*Cmd()` function (e.g., `cmd/<resource>/<resource>.go`).

4. **Add a test** in the corresponding `_test.go` file using `httptest.NewServer`.

5. **Update documentation**:
   - Add a `Long` description with examples to the command
   - Update `README.md` with the new command
   - Update `CLAUDE.md` if it's a commonly-used command
   - Update the skill's `references/commands.md`

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use meaningful variable names
- Every command must have:
  - `Short` description (one line)
  - `Long` description with usage examples
  - Proper flag definitions with descriptions
- Use `-o json` output in all examples for agent-friendliness
- Table output should have meaningful column headers

## Commit Messages

Follow conventional commits:
```
feat: add widget list command
fix: correct pagination in dashboard search
docs: update README with new alert commands
test: add tests for credential CRUD
chore: update dependencies
```

## Pull Requests

1. Fork the repo and create a feature branch
2. Make your changes with tests
3. Run `make test` and `make vet` to ensure everything passes
4. Commit with a clear message
5. Open a PR against `main`

## Releasing

Releases are automated via GoReleaser. To create a release:

```bash
git tag v0.2.0
git push origin v0.2.0
```

This triggers GitHub Actions to:
1. Build binaries for all platforms
2. Create a GitHub Release with assets
3. Generate a changelog

## Reporting Issues

- Use GitHub Issues
- Include: CLI version (`grafana version`), OS/arch, command that failed, error output
- For feature requests, describe the use case

## License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.
