# Changelog

The format of this file is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0] - 2021-05-02

The MVP release of Captain. Supports the ability to create and destroy containers as requested through a REST API.

### Added
- Created a changelog.
- Security policy defined in `SECURITY.md`.
- Automated versioning to compiled output.
- Installation script with GNU make.
- Builder jobs are now multi-threaded.
- Setting `CAPTAIN_DRY_RUN` to `TRUE` will allow users to simulate Captain without modifying the hypervisor state.

### Changed

- `defaults.yaml` is now stored at `/etc/captain/defaults.yaml`.

### Removed

- Use of `DATABASE_CONN` to override database driver type in CI environment.
- `db.conf` is no longer allowed for database configuration. Use `CAPTAIN_DB` instead.
- `:memory:` database type for SQLite3 driver no longer works and has been removed.

### Fixed
- LXC containers now use DHCP for IP address assignment rather than a hard-coded value.

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

[v0.1.0]: https://github.com/ARMmaster17/Captain/releases/tag/v0.1.0
[v0.0.1]: https://github.com/ARMmaster17/Captain/releases/tag/v0.0.1