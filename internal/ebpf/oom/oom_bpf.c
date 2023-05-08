//go:build ignore

#include "../vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

#ifndef TASK_COMM_LEN
#define TASK_COMM_LEN 16
#endif

#define OOM_GLOBAL 0
#define OOM_CGROUP 1

char __license[] SEC("license") = "GPL";

struct {
	__uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
	__uint(key_size, sizeof(u32));
	__uint(value_size, sizeof(u32));
} oom_events SEC(".maps");


struct oom_data_t {
    u64 cgroupMemoryMax;                    // Cgroup Maximum memory limit (value is presented in pages unit)
    unsigned char comm[TASK_COMM_LEN];      // Name of terminated process
    u32 pid;                                // PID of terminated process
    u8 oomType;                             // Type of OOM (global or cgroup related)
};

SEC("kprobe/oom_kill_process")
int kprobe__oom_kill_process(struct pt_regs *ctx)
{
    struct oom_control *oc = (struct oom_control*) PT_REGS_PARM1(ctx);
    struct oom_data_t oom_data = {};
    u64 max = 0;

    struct task_struct *chosen;
    if (bpf_probe_read_kernel(&chosen, sizeof(chosen), &oc->chosen) != 0) {
        bpf_printk("ERROR:oom_kill_process:bpf_probe_read_kernel chosen read fails\n");
        return 0;
    }

    u32 pid;
    if (bpf_probe_read_kernel(&pid, sizeof(pid), &chosen->pid) != 0) {
        bpf_printk("ERROR:oom_kill_process:bpf_probe_read_kernel pid read fails\n");
        return 0;
    }

    unsigned char chosencomm[TASK_COMM_LEN];
    if (bpf_probe_read_kernel(&chosencomm, sizeof(chosencomm), &chosen->comm) != 0) {
        bpf_printk("ERROR:oom_kill_process:bpf_probe_read_kernel comm read fails\n");
        return 0;
    }

    struct mem_cgroup *memcg;
    if (bpf_probe_read_kernel(&memcg, sizeof(memcg), &oc->memcg) != 0) {
        bpf_printk("ERROR:oom_kill_process:bpf_probe_read_kernel memcg read fails\n");
        return 0;
    }

    if (memcg) {
        struct page_counter memory;
        if (bpf_probe_read_kernel(&memory, sizeof(memory), &memcg->memory) != 0) {
            bpf_printk("ERROR:oom_kill_process:bpf_probe_read_kernel memory read fails\n");
            return 0;
        }

        if (bpf_probe_read_kernel(&max, sizeof(max), &memory.max) != 0) {
            bpf_printk("ERROR:oom_kill_process:bpf_probe_read_kernel memory max read fails\n");
            return 0;
        }
    }

    oom_data.cgroupMemoryMax = max;
    for (int i = 0; i < TASK_COMM_LEN; ++i ){
        oom_data.comm[i] = chosencomm[i];
    }
    oom_data.pid = pid;
    oom_data.oomType = memcg ? OOM_CGROUP : OOM_GLOBAL;


    if (bpf_perf_event_output(ctx, &oom_events, BPF_F_CURRENT_CPU,
							  &oom_data, sizeof(oom_data)) != 0) {
		bpf_printk("ERROR:oom:bpf_perf_event_output\n");
	}

    return 0;
}
