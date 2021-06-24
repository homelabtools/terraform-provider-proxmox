# proxmox-box

This is a build for creating a Vagrant box that runs Proxmox VE on Debian.
It was created for end-to-end tests for the Proxmox Terraform provider project,
[danitso/terraform-provider-proxmox](https://github.com/danitso/terraform-provider-proxmox).

It is meant to be used for testing and development, but can also serve
as a source of assistance for anyone wanting to produce production-ready
Proxmox VE VM images. It is probably suitable enough for home lab use.

The `install-proxmox.sh` script contains everything needed to perform a
fully unattended install. The resulting configuration is not production-ready
and has not been tested for any purposes other than development.