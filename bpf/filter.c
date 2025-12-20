// +build ignore
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>

// --- BPF MAPS ---

// 1. BLOCKLIST (O(1) Drop)
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 200000);
    __type(key, __u32);
    __type(value, __u8);
} blocklist SEC(".maps");

// 2. RING BUFFER (Streaming Events to Go)
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24); // 16MB Buffer
} events SEC(".maps");

// Event Structure (Must match Go struct)
struct event_t {
    __u32 src_ip;
    __u32 reason_code; 
    __u16 dst_port;
    __u8  protocol;
} __attribute__((packed));

SEC("xdp")
int aegis_main(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) return XDP_PASS;
    if (eth->h_proto != __constant_htons(ETH_P_IP)) return XDP_PASS;

    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end) return XDP_PASS;

    __u32 src_ip = ip->saddr;

    // 1. ENFORCEMENT: Check Blocklist
    if (bpf_map_lookup_elem(&blocklist, &src_ip)) return XDP_DROP;

    // 2. DETECTION: Telemetry
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)(ip + 1);
        if ((void *)(tcp + 1) > data_end) return XDP_PASS;

        // SYN Flood Logic (SYN=1, ACK=0)
        if (tcp->syn && !tcp->ack) {
             struct event_t *e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
             if (e) {
                 e->src_ip = src_ip;
                 e->reason_code = 1; // 1 = SYN_FLOOD_SUSPICION
                 e->dst_port = tcp->dest;
                 e->protocol = IPPROTO_TCP;
                 bpf_ringbuf_submit(e, 0);
             }
        }
    }
    return XDP_PASS;
}
char _license[] SEC("license") = "GPL";
