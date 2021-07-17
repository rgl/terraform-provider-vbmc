FROM ubuntu:20.04
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
        python3-setuptools \
        python3-pip \
        python3-libvirt \
        python3-cryptography \
        python3-yaml
COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt
COPY . .
ENTRYPOINT ["./entrypoint.sh"]