# see https://github.com/hashicorp/terraform
terraform {
  required_version = "1.2.5"
  required_providers {
    # see https://registry.terraform.io/providers/hashicorp/random
    random = {
      source = "hashicorp/random"
      version = "3.3.2"
    }
    # see https://registry.terraform.io/providers/dmacvicar/libvirt
    # see https://github.com/dmacvicar/terraform-provider-libvirt
    libvirt = {
      source = "dmacvicar/libvirt"
      version = "0.6.14"
    }
    # see https://registry.terraform.io/providers/rgl/vbmc
    # see https://github.com/rgl/terraform-provider-vbmc
    vbmc = {
      source = "rgl/vbmc"
      version = "0.4.0"
    }
  }
}

provider "libvirt" {
  uri = "qemu:///system"
}

variable "prefix" {
  default = "terraform_vbmc_example"
}

output "vm_ip" {
  value = length(libvirt_domain.example.network_interface[0].addresses) > 0 ? libvirt_domain.example.network_interface[0].addresses[0] : ""
}

output "vbmc_port" {
  value = vbmc_vbmc.example.port
}

output "vbmc_address" {
  value = vbmc_vbmc.example.address
}

resource "vbmc_vbmc" "example" {
  domain_name = libvirt_domain.example.name
  port = 6230 # NB when port is unset, the port will be automatically allocated.
}

# see https://github.com/dmacvicar/terraform-provider-libvirt/blob/v0.6.14/website/docs/r/domain.html.markdown
resource "libvirt_domain" "example" {
  name = var.prefix
  cpu {
    mode = "host-passthrough"
  }
  vcpu = 2
  memory = 1024
  qemu_agent = true
  disk {
    volume_id = libvirt_volume.example_root.id
    scsi = true
  }
  network_interface {
    network_id = libvirt_network.example.id
    wait_for_lease = true
    addresses = ["10.17.3.2"]
  }
}

# this uses the vagrant ubuntu image imported from https://github.com/rgl/ubuntu-vagrant.
# see https://github.com/dmacvicar/terraform-provider-libvirt/blob/v0.6.14/website/docs/r/volume.html.markdown
resource "libvirt_volume" "example_root" {
  name = "${var.prefix}_root.img"
  base_volume_name = "ubuntu-20.04-amd64_vagrant_box_image_0_box.img"
  format = "qcow2"
}

# see https://github.com/dmacvicar/terraform-provider-libvirt/blob/v0.6.14/website/docs/r/network.markdown
resource "libvirt_network" "example" {
  name = var.prefix
  mode = "nat"
  domain = "example.test"
  addresses = ["10.17.3.0/24"]
  dhcp {
    enabled = false
  }
  dns {
    enabled = true
    local_only = false
  }
}
