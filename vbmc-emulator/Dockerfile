# syntax=docker/dockerfile:1.9
FROM debian:12-slim
RUN <<EOF
apt-get update -y
apt-get install -y --no-install-recommends \
    python3-setuptools \
    python3-pip \
    python3-libvirt \
    python3-cryptography \
    python3-yaml
rm -rf /var/lib/apt/lists/*
EOF
COPY requirements.txt .
RUN pip3 install --break-system-packages --no-cache-dir -r requirements.txt
COPY . .
ENTRYPOINT ["./entrypoint.sh"]