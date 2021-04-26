![Logo](https://github.com/ARMmaster17/Captain/raw/main/static/Captain.png)
# Captain
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/ARMmaster17/Captain?sort=semver)
[![Go](https://github.com/ARMmaster17/Captain/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ARMmaster17/Captain/actions/workflows/go.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/ade54503d0d7daec431f/maintainability)](https://codeclimate.com/github/ARMmaster17/Captain/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/ade54503d0d7daec431f/test_coverage)](https://codeclimate.com/github/ARMmaster17/Captain/test_coverage)

Captain is a container orchestration and streamlined PaaS provider for Proxmox-based datacenters. Captain can:
- Quickly create LXC containers based on a set of common defaults.
- Seamlessly scale instances up and down.
- Manage DNS with PowerDNS.
- Provide health checks and automatic crash mitigation.
- Easy to use through the CLI, web interface, or GraphQL API.
- Deploy apps without ever having to think about the underlying infrastructure.
- Allow you to build your own "AWS-style" cloud by building services on top of Captain.

## Getting Started
1. Install Proxmox on at least one server (you may also use a hypervisor such as VirtualBox in place of a physical server).
2. Create an LXC container (or VM) in Proxmox using any (reasonable) Linux distro. Install the golang compiler and `make`.
3. Run the following commands:
```shell
git clone https://github.com/ARMmaster17/Captain
cd Captain
go get .
go build
```
4. Optionally, you may create a file at `/etc/captain/db.conf`. Enter a valid PostgreSQL URL (e.g. `postgres://...`) or a valid file path for Sqlite3 to use (e.g. `test.db`).
5. Set the following environment variables to match your Proxmox configuration.
   
   | Name | Value |
   |---|---|
   | `CAPTAIN_DB` | Optionally, this variable can be specified in place of creating `db.conf`. |
   | `CAPTAIN_PROXMOX_USER` | The fully qualified username of a user with privileges to create and destroy VMs/containers. (ex. `root@pam`) |
   | `CAPTAIN_PROXMOX_PASSWORD` | Password to specified proxmox user. |
   | `CAPTAIN_PROXMOX_URL` | Full path to Proxmox host. (ex. `https://192.168.1.2:8006/`) |
   | `CAPTAIN_PRIVATE_KEY` | Absolute filepath to your private key so Captain can provision new planes. |

6. Edit `defaults.yaml` in the captain directory to match the configuration of your network and Proxmox cluster setup. This is also where you provide your public SSH key.
7. Start Captain by running `./captain`.
8. Optionally you may create a SystemD service file so Captain runs on system startup.

If leave `db.conf` blank and don't set `CAPTAIN_DB`, by default Captain will run with an in-memory database that will be cleared on each restart. This is fine for testing, but to run an actual cluster it is recommended to use a Sqlite3 file or PostgreSQL database.

## Managing a Captain Cluster
First a bit of terminology. The highest level in Captain is called an *airspace*. An airspace is an isolated group of instances. For example, one airspace can hold all production instances of an app, and each developer gets their own airspace for testing purposes.

Each airspace has many *flights*. A flight is a complete app, which may include a reverse proxy, a database, and a web server. Each of those services is a *formation*. A formation is a collection of planes (usually containers or VMs) that can be seamlessly scaled up and down.

At the preset moment, the only way to provision planes is to manually edit the Sqlite3 or PostgreSQL database. In the near future it will be possible to provision airspaces, flights, and formations using a REST API, a CLI, or a GraphQL interface. Planes are usually not managed directly by the API, and instead are automatically managed by Captain depending on how the formation is configured.

# Contributing

To contribute, fork this repository, make your changes, and send back a PR. If you're not sure where to start, the issues page lists everything that needs to be done for the current release.
