# Changelog

The format of this file is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

- Radar: New formation field with relative path to preflight Ansible playbook.

### Changed

- Radar: Updated User Interface of WebGUI

## [v0.2.0] - 2021-06-01

### Added

- ATC: All configuration is now stored in a unified config file at `/etc/captain/config.yaml`
  that is generated on first run.
- ATC: A dummy provider driver is now available for testing purposes.
  Helps with simulating the Captain stack without modifying a live
  hypervisor environment.
- CaptainLib: A golang library is now included for each release of Captain for interfacing with
  ATC instances.
- Radar: A new Web GUI has been introduced that natively integrates with ATC using CaptainLib.
- Both ATC and Radar are now available as pre-built DEB files.
- ATC: Logging can now be filtered with the `config.loglevel` property in `config.yaml`.
- ATC: An integrated IPAM module is now available to provider drivers.
- ATC: Formations can now be provisioned with Ansible playbooks.

### Changed

- ATC: Image mappings can now be made on a per-driver basis in `config.yaml`.
- ATC: Moved ATC to the `ATC` directory to make room for other Captain modules.
- ATC: Log output is now shown for all HTTP requests, and not just errors.
- ALL: Makefiles have been revamped to be more unified across the project.
- ATC: Error stacks are now shown one per line instead of a big glob.

### Removed

- ATC: Plane defaults are no longer stored in `defaults.yaml`. This is now part of `config.yaml`.
- ATC: The HTTP listen port for the API can no longer be specified as a command-line parameter.
  This is now part of `config.yaml`.

### Fixed
- Radar: ATC API URL is now configurable in a new `config.yaml`.
- ATC: Fixed issue where DEB packages couldn't be built from source.

### Security

- ATC: The Proxmox LXC driver can now be configured to require TLS certificate checks when
  connecting to the Proxmox API.

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

[v0.2.0]: https://github.com/ARMmaster17/Captain/releases/tag/v0.2.0
[v0.1.0]: https://github.com/ARMmaster17/Captain/releases/tag/v0.1.0
[v0.0.1]: https://github.com/ARMmaster17/Captain/releases/tag/v0.0.1
