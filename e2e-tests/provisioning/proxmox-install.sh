#!/usr/bin/env bash
set -exuo pipefail

PROXMOX_MIRROR="http://download.proxmox.com/debian"

main() {
    apt-get -y update
    apt-get -y upgrade
    apt-get -y install gnupg2
    hostnamectl set-hostname proxmox-e2etests.test.com --static
    echo "127.0.0.1 proxmox-e2etests.test.com proxmox-e2etests" | tee -a /etc/hosts
    wget -qO - "$PROXMOX_MIRROR/proxmox-ve-release-6.x.gpg" | apt-key add -
    echo "deb $PROXMOX_MIRROR/pve buster pve-no-subscription" | tee /etc/apt/sources.list.d/pve-install-repo.list
    echo "deb $PROXMOX_MIRROR/ceph-nautilus buster main" | tee /etc/apt/sources.list.d/ceph.list
    apt-get -y update
    apt-get -y install proxmox-ve postfix open-iscsi
    #reboot
}

main