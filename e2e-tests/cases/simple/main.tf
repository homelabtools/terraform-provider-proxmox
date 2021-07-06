resource "proxmox_virtual_environment_role" "example" {
  privileges = [
    "VM.Monitor",
  ]
  role_id = "terraform-provider-proxmox-example"
}