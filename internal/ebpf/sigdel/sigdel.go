package sigdel

import "C"

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $GOARCH -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf sigdel_bpf.c -- -I/usr/include/bpf -I.

import (
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

type SigDelStruct struct {
	objs bpfObjects
	link link.Link
}

// Setup Sigdel structure
func Setup() (*SigDelStruct, error) {
	var err error
	fn := "signal_deliver"

	sigdel := new(SigDelStruct)

	if err = rlimit.RemoveMemlock(); err != nil {
		return nil, err
	}

	// Load BPF code
	if err = loadBpfObjects(&sigdel.objs, nil); err != nil {
		return nil, err
	}

	// Attach BPF code
	sigdel.link, err = link.Tracepoint("signal", fn, sigdel.objs.SigDeliver, nil)
	if err != nil {
		sigdel.objs.Close()
		return nil, err
	}
	return sigdel, nil
}

// Teardown Sigdel structure
func (sigdel *SigDelStruct) Teardown() {
	sigdel.objs.Close()
	sigdel.link.Close()
}
