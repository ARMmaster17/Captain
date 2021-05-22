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

# Getting Started

1. Install Proxmox on at least one server (you may also use a hypervisor such as VirtualBox in place of a physical server).
2. Create an LXC container (or VM) in Proxmox using any Debian or Ubuntu image.
3. Download the DEB files from GitHub Actions or from the Releases page.
4. Run the following commands:

```shell
sudo apt install ./captain-atc_*_amd64.deb
sudo apt install ./captain-radar_*_amd64.deb # Only needed if you want the web GUI.
nano /etc/captain/atc/config.yaml # Edit to match your proxmox cluster configuration.
sudo systemctl enable captain-atc
sudo systemctl start captain-atc
sudo systemctl enable captain-radar
sudo systemctl start captain-radar
```

You should be able to access the API on `<IP>:5000` and the web GUI on `<IP>:5001`.

| Name | Value |
   |---|---|
| `CAPTAIN_DB` | Path to a SQLite3 file or a Postgres connection string starting with `postgres://...`. |
| `CAPTAIN_PROXMOX_USER` | The fully qualified username of a user with privileges to create and destroy VMs/containers. (ex. `root@pam`) |
| `CAPTAIN_PROXMOX_PASSWORD` | Password to specified proxmox user. |
| `CAPTAIN_PRIVATE_KEY` | Absolute filepath to your private key so Captain can provision new planes. |

## Managing a Captain Cluster
First a bit of terminology. The highest level in Captain is called an *airspace*. An airspace is an isolated group of instances. For example, one airspace can hold all production instances of an app, and each developer gets their own airspace for testing purposes.

Each airspace has many *flights*. A flight is a complete app, which may include a reverse proxy, a database, and a web server. Each of those services is a *formation*. A formation is a collection of planes (usually containers or VMs) that can be seamlessly scaled up and down.

To modify the state database to trigger builds in Proxmox, you may
use a tool like Curl, or the CLI tool (migrating to this repo soon).

## Building From Source

Running `make build` inside the ATC or Radar sub-directories will build executables that can be run in-place.
If you would like to build your own DEB files, run `make deb` in each project directory. Then you may install them
following the steps above.

# Contributing

To contribute, fork this repository, make your changes, and send back a PR. If you're not sure where to start, the issues page lists everything that needs to be done for the current release.
