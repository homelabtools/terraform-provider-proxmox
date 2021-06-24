#!/usr/bin/env bash
set -exuo pipefail

MOUNT_POINT="${MOUNT_POINT:-/mnt/vb}"
ISO_URL="${ISO_URL:-"https://download.virtualbox.org/virtualbox/${VBOX_VERSION}/VBoxGuestAdditions_${VBOX_VERSION}.iso"}"

main() {
    echo "Installing VBoxGuestAdditions for VBox version '$VBOX_VERSION'"
    apt-get -y install build-essential dkms pve-headers
    iso_path="/tmp/vbga.iso"
    curl -sSLo "$iso_path" "$ISO_URL"
    mkdir -p "$MOUNT_POINT"
    mount -o loop "$iso_path" "$MOUNT_POINT"
    "$MOUNT_POINT/VBoxLinuxAdditions.run"
    umount "$MOUNT_POINT"
    rmdir "$MOUNT_POINT"
    rm "$iso_path"
}

main