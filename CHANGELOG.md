# Changelog

The format of this file is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Added `CHANGELOG.md`.
- Added security policy in `SECURITY.md`.

### Security
- Defined security policy to provide security patches for latest patch version of the latest minor release.

## [v0.0.1] - 2021-04-25

### Added
- Plane creation through LXC provider driver.
- State synchronization with Postgres or Sqlite3 database.
- REST API for receiving commands from external tools.
- Swagger-compatible documentation for REST API.
- Multi-level logging output to `stdout` and `stderr`.
- Project logo in `README.md`.
- Unit testing and code coverage in CI pipeline.

[unreleased]: https://github.com/ARMmaster17/Captain/compare/v0.0.1...HEAD
[v0.0.1]: https://github.com/ARMmaster17/Captain/releases/tag/v0.0.1