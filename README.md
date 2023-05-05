# scope-ebpf
eBPF capabilities potentially used with AppScope

![example workflow](https://github.com/criblio/scope-ebpf/actions/workflows/build.yml/badge.svg)


## Build

Pull a copy of the code with:

```bash
git clone https://github.com/criblio/scope-ebpf.git
cd scope-ebpf
```

Build directly on the host machine

```bash
make all
```

Build docker image:

```bash
make build-container
```

## Run

Run the scope-ebpf loader:

```bash
sudo ./bin/scope-ebpf
```
