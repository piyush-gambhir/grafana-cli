# Grafana compatibility and maintenance

Checked 2026-09-06. These are documentation and local-test baselines, not certification against every server deployment.

## Build requirements

- Go 1.27 or later; this checkout selects Go 1.27.3 using `toolchain` in `cli-go/go.mod`.
- Go 1.27 builds for macOS require macOS 13 or later. See the [Go 1.27 release notes](https://go.dev/doc/go1.27).
- Docs: Node.js 26.12.1 and pnpm 12.15.0. Install with `pnpm install --frozen-lockfile` in `web/`.
- YAML uses the maintained [YAML organization v3 implementation](https://github.com/yaml/go-yaml), preserving the v3 configuration API.

## Upstream API baseline

API reference baseline: [Grafana 13.2 HTTP API](https://grafana.com/docs/grafana/v13.2/developer-resources/api-reference/http-api/). Upstream release checked: [13.2.1](https://github.com/grafana/grafana/releases/tag/v13.2.1).

This CLI currently uses Grafana's legacy `/api` endpoints. Grafana documents these as deprecated in favor of the versioned `/apis` resource APIs. Updating the Go dependencies does not migrate the wire protocol. Dashboard, folder, and other legacy commands require those routes to remain available on your deployment. Full migration to `/apis` requires response adapters and integration testing; this checkout does not claim that migration is complete.

## Maintaining this baseline

Dependabot is configured to propose weekly Go, npm, and GitHub Actions updates. Review migrations before merging; CI verifies the Go tests, race checks, security/static analysis, command documentation, and docs build.

Run `go list -m -u all` from `cli-go/` and `pnpm outdated` from `web/` to check newer releases. Commit `go.mod` with `go.sum`, and `package.json` with `pnpm-lock.yaml`. Keep `fumadocs-core` and the `fumadocs-ui` alias to `@fumadocs/base-ui` on matching versions.

The website's static search uses Fumadocs' default search engine. After building, run `pnpm test:search` to verify that the exported index can be loaded and queried by the installed client.

## Node tooling verification (26 October 2026)

The private documentation site now uses Node 26.12.1, pnpm 12.15.0, and the
current matching Fumadocs core/UI and MDX releases. The existing `.nvmrc`, CI,
and package engines agree; TypeScript 7.0.2 works directly without a second
compiler. Newer main security configuration and bounded overrides are preserved.
The Go CLI's source, runtime contract, and binary release version are unchanged.

Frozen installation, peers, typecheck, Cloudflare static export, search-index
and client checks, and Markdown link checks pass. Wrangler's deployment dry-run
packages 276 static files. Full and production dependency audits report zero
known advisories; `pnpm outdated` reports no direct dependency updates. These
checks establish build and search compatibility, not a CLI runtime speedup or
a complete application security audit.
