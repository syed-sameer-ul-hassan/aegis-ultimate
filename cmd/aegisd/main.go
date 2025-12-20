package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
    "github.com/aegis-ultimate/internal/firewall"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("🛡️  AEGIS-X ULTIMATE v8.0 Starting...")

	// 1. Unlimit Memory for eBPF
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Memlock removal failed:", err)
	}

	// 2. Load eBPF Objects
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("Loading BPF objects failed: %v", err)
	}
	defer objs.Close()

	// 3. Initialize Firewall
	fw := firewall.NewNftablesEnforcer()
	
	// 4. Attach RingBuffer Reader
	rd, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		log.Fatalf("Ringbuf init failed: %v", err)
	}
	defer rd.Close()

	log.Println("✅ Kernel Pipeline Attached. Monitoring...")

	// 5. Event Loop
	go func() {
		for {
			record, err := rd.Read()
			if err != nil { continue }

			var e struct {
				SrcIP      uint32
				ReasonCode uint32
				DstPort    uint16
				Protocol   uint8
			}
			if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &e); err == nil {
                if e.ReasonCode == 1 {
                    ip := intToIP(e.SrcIP)
                    log.Printf("⚠️  Threat Detected (SYN Flood) from %s", ip)
                    if err := fw.Block(ip); err != nil {
                        log.Printf("Block failed: %v", err)
                    } else {
                        log.Printf("⛔ BLOCKED: %s", ip)
                    }
                }
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func intToIP(nn uint32) string {
	ip := make(net.IP, 4)
	binary.LittleEndian.PutUint32(ip, nn)
	return ip.String()
}
