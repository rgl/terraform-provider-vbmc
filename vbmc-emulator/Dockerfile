# syntax=docker/dockerfile:1.25
FROM debian:13-slim
RUN <<EOF
apt-get update -y
apt-get install -y --no-install-recommends \
    python3-cryptography \
    python3-libvirt \
    python3-pip \
    python3-setuptools \
    python3-venv \
    python3-yaml
rm -rf /var/lib/apt/lists/*
python3 -m venv --system-site-packages /opt/venv
EOF
ENV PATH="/opt/venv/bin:$PATH"
COPY requirements.txt .
RUN python -m pip --disable-pip-version-check --no-cache-dir install -r requirements.txt
COPY . .
ENTRYPOINT ["./entrypoint.sh"]