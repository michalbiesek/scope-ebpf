DOCKER_BUILDER_IMAGE_NAME ?= criblo/scope-ebpf-builder
DOCKER_BUILDER_TAG ?= local
DOCKER_IMAGE_NAME ?= criblio/scope-epbf
DOCKER_IMAGE_TAG ?= latest
EBPF_LOADER := scope-ebpf
BPFTOOL ?= bpftool
CLANG ?= clang
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)
BTF_VMLINUX ?= /sys/kernel/btf/vmlinux
EBPF_DIR := internal/ebpf

ARCH := $(shell uname -m)
GOARCH := $(subst aarch64,arm64,$(subst x86_64,amd64,$(ARCH)))
BPF_ARCH := $(subst aarch64,arm64,$(subst x86_64,x86,$(ARCH)))

GO ?= $(shell which go 2>&1)
ifeq (,$(GO))
$(error "error: \`go\` not in PATH; install or set GO to it's path")
endif

# Define a variable to store the list of Go files
GO_FILES := $(shell find . -name "*.go" ! -name "*bpfel*.go" -type f)

all: build
build: scope-ebpf

docker-builder:
	docker build -t $(DOCKER_BUILDER_IMAGE_NAME):$(DOCKER_BUILDER_TAG) --file docker/builder/Dockerfile .

image: docker-builder
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) --file docker/base/Dockerfile .

clean:
	$(RM) bin/${EBPF_LOADER}
	$(RM) internal/ebpf/vmlinux.h
	@$(foreach entry,$(wildcard $(EBPF_DIR)/*), \
		if [ -d "$(entry)" ]; then \
			$(RM) $(entry)/bpf_bpfel_${BPF_ARCH}.o; \
			$(RM) $(entry)/bpf_bpfel_${BPF_ARCH}.go; \
		fi; \
	)

scope-ebpf: generate
	$(GO) build -ldflags="-extldflags=-static" -o bin/${EBPF_LOADER} ./cmd/scope-ebpf
	chmod +x bin/${EBPF_LOADER}

fmt:
	@for file in $(GO_FILES); do \
		$(GO) fmt $$file; \
	done

generate: export BPF_CLANG := $(CLANG)
generate: export BPF_CFLAGS := $(CFLAGS)
generate: vmlinux
	@$(foreach entry,$(wildcard $(EBPF_DIR)/*), \
		if [ -d "$(entry)" ]; then \
			$(GO) generate $(entry)/$(notdir $(entry)).go; \
		fi; \
	)

help:
	@echo "Available targets:"
	@echo "  all             - Default target, builds the scope-ebpf binary"
	@echo "  build           - Builds the scope-ebpf binary"
	@echo "  docker-builder  - Builds the scope-ebpf docker image builder"
	@echo "  image           - Builds the scope-ebpf docker image"
	@echo "  clean           - Cleans up build artifacts"
	@echo "  scope-ebpf      - Builds the scope-ebpf binary"
	@echo "  fmt             - Formats Go source files"
	@echo "  generate        - Generates Go code for ebpf programs"
	@echo "  vet             - Runs Go vet on source files"
	@echo "  vmlinux         - Generates vmlinux.h header file"

vet:
	@for file in $(GO_FILES); do \
		$(GO) vet $$file; \
	done

vmlinux:
	$(BPFTOOL) btf dump file $(BTF_VMLINUX) format c > internal/ebpf/vmlinux.h

.PHONY: all build clean docker-builder fmt generate help image scope-ebpf vet vmlinux
