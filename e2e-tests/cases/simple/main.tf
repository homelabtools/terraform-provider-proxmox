provider "proxmox" {
  virtual_environment {
    endpoint = "https://127.0.0.1:8006"
    username = "root@pam"
    password = "proxmox"
    insecure = true
  }
}

terraform {
  required_version = ">=1.0.0"
  required_providers {
    proxmox = {
      source = "registry.terraform.io/danitso/proxmox"
    }
  }
}

resource "proxmox_virtual_environment_role" "example" {
  privileges = [
    "VM.Monitor",
  ]
  role_id = "terraform-provider-proxmox-example"
}