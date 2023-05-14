package oom

import "C"

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $GOARCH -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf oom_bpf.c -- -I/usr/include/bpf -I.

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

type OomType byte

const (
	OomGlobal OomType = 0
	OomCgroup OomType = 1
)

func (o OomType) String() string {
	switch o {
	case OomGlobal:
		return "global"
	case OomCgroup:
		return "cgroup"
	}
	return "unknown"
}

// oomEvent is corresponding to oom_data_t structure
// IMPORTANT:
// The follow structure is used in binary.Read so please consider
// padding. See details in https://github.com/cilium/ebpf/issues/821
type oomEvent struct {
	CgroupMemLimit uint64
	Com            [16]byte
	Pid            uint32
	OomInfo        OomType
}

// Returns string value of prometheus metrics
func (oe *oomEvent) String() string {
	return fmt.Sprintf("pid=\"%d\" name=\"%s\", oomtype=\"%s\", cgrouppageLimit=\"%d\"", oe.Pid, oe.Com, oe.OomInfo.String(), oe.CgroupMemLimit)
}

func SetupReadOOM(oomEventChan chan string) error {
	fn := "oom_kill_process"
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		return err
	}
	objs := bpfObjects{}

	// Load BPF code
	if err := loadBpfObjects(&objs, nil); err != nil {
		return err
	}
	defer objs.Close()

	// Attach BPF code
	link, err := link.Kprobe(fn, objs.KprobeOomKillProcess, nil)
	if err != nil {
		objs.Close()
		return err
	}
	defer link.Close()

	rd, err := ringbuf.NewReader(objs.OomEvents)
	if err != nil {
		return err
	}
	defer rd.Close()

	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				return err
			}
			continue
		}

		var event oomEvent
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			continue
		}

		oomEventChan <- event.String()
	}
}
