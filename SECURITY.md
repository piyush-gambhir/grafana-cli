# Security notes

## CodeQL build assertion triage (2026-10-29)

`web/scripts/check-llms.mjs` checks whether built documentation contains the
expected absolute docs URL. Its `full.includes(...)` assertion is a content
smoke test, not a URL authorization or sanitization check. It does not fetch that
URL or authorize a host. The `js/incomplete-url-substring-sanitization` finding
is a false positive in this build assertion. Any future network allowlist must
parse URLs and compare their origin explicitly.
