package oom

import "C"

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $GOARCH -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf oom_bpf.c -- -I/usr/include/bpf -I.

import (
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

type OomStruct struct {
	objs bpfObjects
	link link.Link
}

// Setup Oom structure
func Setup() (*OomStruct, error) {
	var err error
	fn := "oom_kill_process"

	oomInst := new(OomStruct)

	// Allow the current process to lock memory for eBPF resources.
	if err = rlimit.RemoveMemlock(); err != nil {
		return nil, err
	}

	// Load BPF code
	if err = loadBpfObjects(&oomInst.objs, nil); err != nil {
		return nil, err
	}

	// Attach BPF code
	oomInst.link, err = link.Kprobe(fn, oomInst.objs.KprobeOomKillProcess, nil)
	if err != nil {
		oomInst.objs.Close()
		return nil, err
	}
	return oomInst, nil

}

// Teardown Oom structure
func (oom *OomStruct) Teardown() {
	oom.objs.Close()
	oom.link.Close()
}
