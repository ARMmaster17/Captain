[![Build Status](http://jenkins.firecor.me/buildStatus/icon?job=Captain+CI%2Fmain)](http://jenkins.firecor.me/job/Captain%20CI/job/main/)

![Logo](https://github.com/ARMmaster17/Captain/raw/main/static/Captain.png)

Captain is a container orchestration and streamlined PaaS provider for Proxmox-based datacenters. Captain can:
- Quickly create LXC containers based on a set of common defaults.
- Seamlessly scale instances up and down.
- Integrate with existing PowerDNS or phpIPAM installations.
- Provide health checks and automatic crash mitigation.
- Easy to use through the CLI, web interface, or GraphQL API.
- Deploy apps without ever having to think about the underlying infrastructure.
- Allow you to build your own "AWS-style" cloud by building services on top of Captain.

See the wiki on how to get started quickly. As of right now, Captain makes a few assumptions about your infrastructure (will be changed in the future):
- IP address management is done with an on-site phpIPAM installation.
- Local authoritative DNS services are provided by PowerDNS.
- You have at least one Proxmox server set up in cluster configuration.