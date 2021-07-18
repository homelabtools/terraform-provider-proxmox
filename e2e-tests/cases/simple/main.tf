resource "proxmox_virtual_environment_role" "example" {
  privileges = [
    "VM.Monitor",
  ]
  role_id = "test-role"
}

provider "proxmox" {
  virtual_environment {
    endpoint = "http://localhost:8000"
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
