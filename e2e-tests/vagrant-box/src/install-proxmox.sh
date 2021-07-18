#!/usr/bin/env bash
set -exuo pipefail

PROXMOX_MIRROR="${PROXMOX_MIRROR:-http://download.proxmox.com/debian}"
HOSTNAME="${HOSTNAME:-proxmox-e2etests}"
FQDN="${FQDN:-$HOSTNAME.internal}"
IFACE="${IFACE:-eth0}"

main() {
    # Proxmox expects the hostname/FQDN to map to the IP that it's listening on.
    ip="$(hostname -I)"
    tee /etc/hosts << HERE
127.0.0.1 localhost.internal localhost
$ip $FQDN $HOSTNAME
HERE

    hostnamectl set-hostname "$FQDN" --static

    apt install -y gnupg2 debian-keyring debian-archive-keyring apt-transport-https

    # Fetch Proxmox key
    wget -qO - "$PROXMOX_MIRROR/proxmox-ve-release-6.x.gpg" | apt-key add -
    echo "deb $PROXMOX_MIRROR/pve buster pve-no-subscription" | tee /etc/apt/sources.list.d/pve-install-repo.list
    echo "deb $PROXMOX_MIRROR/ceph-nautilus buster main" | tee /etc/apt/sources.list.d/ceph.list

    # Fetch Caddy key
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | apt-key add -
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list

    apt -y update
    apt -y install caddy

    export DEBIAN_FRONTEND=noninteractive
    apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade
    apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" install proxmox-ve postfix open-iscsi

    rm /etc/apt/sources.list.d/pve-enterprise.list
    apt remove os-prober  # As recommended by Proxmox wiki

    printf 'proxmox\nproxmox\n' | passwd root
}

main
