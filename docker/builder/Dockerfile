FROM ubuntu:22.04

# Install OS software
RUN apt-get update && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
        clang \
        curl \
        gcc \
        golang-go \
        git \
        gpg-agent \
        libelf-dev \
        llvm \
        make \
        && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

# Install libbpf and bpftool
ARG LIBBPF_VERSION=v1.2.0
ARG BPFTOOL_VERSION=v7.2.0
RUN cd /tmp && \
    mkdir /tmp/libbpf && \
    curl -Ls https://github.com/libbpf/libbpf/archive/refs/tags/${LIBBPF_VERSION}.tar.gz | tar zxvf - -C /tmp/libbpf --strip-components 1 && \
    cd /tmp/libbpf/src && \
    make && \
    make install && \
    rm -rf tmp/libbpf && \
    cd /tmp && \
    mkdir /tmp/bpftool && \
    git clone https://github.com/libbpf/bpftool.git --branch ${BPFTOOL_VERSION} --recurse-submodules --single-branch /tmp/bpftool && \
    cd /tmp/bpftool/src && \
    make && \
    make install && \
    rm -rf tmp/bpftool

COPY . /opt/scope-ebpf

WORKDIR /opt/scope-ebpf

RUN make all && \
    cp /opt/scope-ebpf/bin/scope-ebpf /bin/scope-ebpf
