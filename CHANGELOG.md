# Changelog

All notable changes to this project are documented here. This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] — 2026-05-03

### Security

- Resolved Dependabot / `pnpm audit` findings for the local **MCP Inspector** stack under `scripts/inspector`: upgraded `@modelcontextprotocol/inspector` to **0.21.x** and applied **`pnpm.overrides`** so transitive dependencies (for example `minimatch`, `path-to-regexp`, `qs`, `body-parser`, `ajv`, `brace-expansion`, `diff`) resolve to patched versions.
- Go module graph checked with **govulncheck** (no vulnerabilities reported for this release at scan time).

### Changed

- **Dependabot** now also tracks the **npm** manifest in `scripts/inspector` (weekly), alongside `gomod` at the repository root.
- **`scripts/inspector/run.sh`**: install with **pnpm** (aligned with `pnpm-lock.yaml`), resolve the inspector directory via the script path, and invoke **`mcp-inspector`** via `pnpm exec`.

## [0.2.9] — 2026-05-02

### Added

- TLS for the Nomad HTTP client via `NOMAD_CACERT`, `NOMAD_SKIP_VERIFY`, and `NOMAD_TLS_SERVER_NAME`.
- Documentation for MCP transports: **stdio** (default), **sse**, and **streamable-http** (including Inspector URLs for local HTTP).
- Stronger allocation-related tooling and namespace handling across Nomad MCP tools.
- `context.Context` threaded through Nomad client calls; canonical `Authorization: Bearer` handling and clearer HTTP error paths.
- Unit tests and refactors around effective namespace resolution in prompts.

### Changed

- Build pipeline uses a patched Go toolchain (`go.mod` version alignment).
- Dependency: `github.com/mark3labs/mcp-go` updated from v0.43.2 through **v0.49.0** (series of bumps).
- Build scripts and auxiliary configuration refreshed.

### Contributors

Thanks for community pull requests merged since [v0.2.8](https://github.com/kocierik/mcp-nomad/releases/tag/v0.2.8):

- [@zw5](https://github.com/zw5) — [PR #39](https://github.com/kocierik/mcp-nomad/pull/39) (patched Go toolchain builds).
- Patrick B. ([@Theragus](https://github.com/Theragus)) — [PR #36](https://github.com/kocierik/mcp-nomad/pull/36) (Nomad TLS environment variables).

Sequential `github.com/mark3labs/mcp-go` upgrades after v0.2.8 were proposed and landed via Dependabot pull requests merged by maintainers/GitHub Actions: [#28](https://github.com/kocierik/mcp-nomad/pull/28), [#29](https://github.com/kocierik/mcp-nomad/pull/29), [#30](https://github.com/kocierik/mcp-nomad/pull/30), [#31](https://github.com/kocierik/mcp-nomad/pull/31), [#33](https://github.com/kocierik/mcp-nomad/pull/33), [#34](https://github.com/kocierik/mcp-nomad/pull/34), [#37](https://github.com/kocierik/mcp-nomad/pull/37), [#38](https://github.com/kocierik/mcp-nomad/pull/38).

## [0.2.8] — 2026-02-16

### Added

- Expanded testing infrastructure and CI/CD workflow coverage.

### Changed

- Dependency: `github.com/mark3labs/mcp-go` updated across the **v0.42.x → v0.43.x** range (see git history for individual bumps).

### Contributors

Automated bumps from Dependabot, merged via pull requests [#24](https://github.com/kocierik/mcp-nomad/pull/24), [#25](https://github.com/kocierik/mcp-nomad/pull/25), [#26](https://github.com/kocierik/mcp-nomad/pull/26).

## [0.2.7] — 2025-10-28

### Fixed

- Handle null timestamp fields from the Nomad API when parsing responses.

### Changed

- Dependency: `github.com/mark3labs/mcp-go` updated (v0.41.x → v0.42.0).

### Contributors

Dependabot proposals merged as [#20](https://github.com/kocierik/mcp-nomad/pull/20), [#21](https://github.com/kocierik/mcp-nomad/pull/21). Timestamp parsing fix tracked as [#23](https://github.com/kocierik/mcp-nomad/pull/23).

## [0.2.6] — 2025-10-02

Maintenance and dependency updates; see [git compare](https://github.com/kocierik/mcp-nomad/compare/v0.2.5...v0.2.6) for commit-level detail.

### Contributors

- Nick Wales ([@nickwales](https://github.com/nickwales)) — [PR #17](https://github.com/kocierik/mcp-nomad/pull/17) (`httpServer` bound to `0.0.0.0`).
- Austin Culter ([@ChefAustin](https://github.com/ChefAustin)) — [PR #8](https://github.com/kocierik/mcp-nomad/pull/8) — *fix: namespace* (merged 2025-06-20 on GitHub). Addresses namespace-related Nomad API usage in `client` and MCP tooling. Included from [v0.2.5](https://github.com/kocierik/mcp-nomad/releases/tag/v0.2.5).

Dependabot merges in this timeframe include [#12](https://github.com/kocierik/mcp-nomad/pull/12) through [#18](https://github.com/kocierik/mcp-nomad/pull/18) plus [#15](https://github.com/kocierik/mcp-nomad/pull/15).

---

[Unreleased]: https://github.com/kocierik/mcp-nomad/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/kocierik/mcp-nomad/compare/v0.2.9...v0.3.0
[0.2.9]: https://github.com/kocierik/mcp-nomad/compare/v0.2.8...v0.2.9
[0.2.8]: https://github.com/kocierik/mcp-nomad/compare/v0.2.7...v0.2.8
[0.2.7]: https://github.com/kocierik/mcp-nomad/compare/v0.2.6...v0.2.7
[0.2.6]: https://github.com/kocierik/mcp-nomad/compare/v0.2.5...v0.2.6
