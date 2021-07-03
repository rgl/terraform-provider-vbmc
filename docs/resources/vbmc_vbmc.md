---
page_title: vbmc_vbmc Resource - terraform-provider-vbmc
subcategory:
description: |-
  Manage a libvirt_domain VirtualBMC (vbmc).
---

# vbmc_vbmc (Resource)

This manages a [libvirt_domain](https://github.com/dmacvicar/terraform-provider-libvirt) [VirtualBMC (vbmc)](https://github.com/openstack/virtualbmc).

## Example Usage

This is normally used as:

```terraform
resource "vbmc_vbmc" "example" {
  domain_name = libvirt_domain.example.name
  port = 6230
}

resource "libvirt_domain" "example" {
  name = "example"
  ...
}
```

After `terraform apply`, the vbmc will be running at `127.0.0.1:6230` with the default username and password. It can be used from [ipmitool](https://github.com/ipmitool/ipmitool) to control the `libvirt_domain` as, e.g.:

```bash
ipmitool -I lanplus -H 127.0.0.1 -p 6230 -U admin -P password chassis bootdev pxe
ipmitool -I lanplus -H 127.0.0.1 -p 6230 -U admin -P password chassis power reset
```

For a complete example see [rgl/terraform-provider-vbmc](https://github.com/rgl/terraform-provider-vbmc).

## Schema

### Required

- **domain_name** (String) The libvirt domain name. This should reference an existing `libvirt_domain` resource.
- **port** (Number) The vbmc port.

### Optional

- **address** (String) The vbmc address. Defaults to `127.0.0.1`.
- **username** (String) The vbmc username. Defaults to `admin`.
- **password** (String, Sensitive) The vbmc username password. Defaults to `password`.
