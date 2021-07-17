#!/bin/bash
set -euo pipefail

install -d ~/.vbmc/$VBMC_EMULATOR_DOMAIN_NAME
cat >~/.vbmc/$VBMC_EMULATOR_DOMAIN_NAME/config <<EOF
[VirtualBMC]
username = $VBMC_EMULATOR_USERNAME
password = $VBMC_EMULATOR_PASSWORD
address = 0.0.0.0
port = 6230
domain_name = $VBMC_EMULATOR_DOMAIN_NAME
libvirt_uri = qemu:///system
active = True
EOF

exec vbmcd --foreground "$@"
