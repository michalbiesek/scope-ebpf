# scope-ebpf

![example workflow](https://github.com/criblio/scope-ebpf/actions/workflows/build.yml/badge.svg)

## Contents

scope-ebpf is an eBPF loader.

## Contents
1. [Dependencies](#dependencies)
2. [Build](#build)
    * [Docker](#docker)
    * [Host](#host)
3. [Run](#run)
4. [AppScope integration](#appscope-integration)


## Dependencies

See the **[Dockerfile](Dockerfile)**  to get an idea what packages are required to build scope-epbf project.

You will need to install the following required packages in the system, to build `scope-ebpf` from sources on host:

* **build-essential**
* **clang**
* **golang**
* **llvm**
* **libbpf-dev**
* **linux-tools**

## Build

Pull a copy of the code with:

```bash
git clone https://github.com/criblio/scope-ebpf.git
cd scope-ebpf
```

### Docker 

To build the Docker image, run the following command:

```bash
make build-container
```

### Host 

To build scope-ebpf directly on the host machine, run the following command:

```bash
make all
```

## Run

To run the scope-ebpf loader, run the following command:

```bash
sudo ./bin/scope-ebpf
```

## AppScope integration

The scope-ebpf can be used with conjuction of [AppScope](https://appscope.dev/).
In the following example, `scope-ebpf` is responsible for loading the eBPF code that the `scope daemon` process will read from the eBPF maps.

```bash
sudo ./scope-ebpf &
sudo ./scope daemon 
```
