# About

[![build](https://github.com/rgl/terraform-provider-vbmc/actions/workflows/build.yml/badge.svg)](https://github.com/rgl/terraform-provider-vbmc/actions/workflows/build.yml) [![terraform provider](https://img.shields.io/badge/terraform%20provider-rgl%2Fvbmc-blue)](https://registry.terraform.io/providers/rgl/vbmc)

This manages a [libvirt_domain](https://github.com/dmacvicar/terraform-provider-libvirt) [VirtualBMC (vbmc)](https://github.com/openstack/virtualbmc) through the [`vbmc_vbmc` resource](https://github.com/rgl/terraform-provider-vbmc/blob/main/docs/resources/vbmc.md).

For a Redfish based provider see the [rgl/terraform-provider-sushy-vbmc](https://github.com/rgl/terraform-provider-sushy-vbmc) source repository.

## Usage (Ubuntu 20.04 host)

Install Terraform:

```bash
wget https://releases.hashicorp.com/terraform/1.0.1/terraform_1.0.1_linux_amd64.zip
unzip terraform_1.0.1_linux_amd64.zip
sudo install terraform /usr/local/bin
rm terraform terraform_*_linux_amd64.zip
```

Install VirtualBMC and ipmitool:

```bash
# see https://github.com/openstack/virtualbmc
# see https://docs.openstack.org/virtualbmc/latest/user/index.html
python3 -m pip install virtualbmc
sudo apt-get install -y ipmitool
```

Start the virtual bmc daemon:

```bash
vbmcd
```

Build the development version of this provider and install it:

**NB** This is only needed when you want to develop this plugin. If you just want to use it, let `terraform init` install it [from the terraform registry](https://registry.terraform.io/providers/rgl/vbmc).

```bash
make
```

Create the infrastructure:

```bash
terraform init
terraform plan -out=tfplan
terraform apply tfplan
```

**NB** if you have errors alike `Could not open '/var/lib/libvirt/images/terraform_vbmc_example_root.img': Permission denied'` you need to reconfigure libvirt by setting `security_driver = "none"` in `/etc/libvirt/qemu.conf` and restart libvirt with `sudo systemctl restart libvirtd`.

Show information about the libvirt/qemu guest:

```bash
virsh dumpxml terraform_vbmc_example
virsh qemu-agent-command terraform_vbmc_example '{"execute":"guest-info"}' --pretty
```

Show information about the VM vbmc:

```bash
vbmc list
vbmc show terraform_vbmc_example
```

Create the `vbmc_ipmitool` alias to make `ipmitool` simpler to use:

```bash
alias vbmc_ipmitool="\
ipmitool \
-I lanplus \
-U admin \
-P password \
-H "$(terraform output --raw vbmc_address)" \
-p "$(terraform output --raw vbmc_port)" \
"
```

Show the power status:

```bash
vbmc_ipmitool chassis status
vbmc_ipmitool chassis power status
```

Do a soft power off (ACPI shutdown):

**NB** A soft power off will be handled by the `qemu-ga` daemon and the `/var/log/syslog` file contains the lines `qemu-ga: info: guest-shutdown called, mode powerdown.` and `systemd: Stopped target Default.`.

```bash
vbmc_ipmitool chassis power soft # NB use "off" for a hard power off.
vbmc_ipmitool chassis power status
```

Set the machine boot device to PXE boot from the default network interface and power it on:

```bash
vbmc_ipmitool chassis bootdev pxe
vbmc_ipmitool chassis bootparam get 5 # get the current boot device.
vbmc_ipmitool chassis power on
vbmc_ipmitool chassis power status
```

Set the machine boot device to boot from the default disk and reset it:

```bash
vbmc_ipmitool chassis bootdev disk
vbmc_ipmitool chassis bootparam get 5 # get the current boot device.
vbmc_ipmitool chassis power reset # NB this is an hard-reset.
```

Destroy the infrastructure:

```bash
terraform destroy -target vbmc_vbmc.example # destroy just the vbmc.
terraform destroy -auto-approve # destroy everything.
```
