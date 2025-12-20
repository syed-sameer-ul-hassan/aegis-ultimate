package firewall

import (
	"fmt"
	"os/exec"
	"strings"
)

type NftablesEnforcer struct {
	Table string
	Set   string
}

func NewNftablesEnforcer() *NftablesEnforcer {
	n := &NftablesEnforcer{Table: "aegis", Set: "blocklist"}
	n.init()
	return n
}

func (n *NftablesEnforcer) init() {
    // Idempotent initialization
	exec.Command("nft", "add", "table", "inet", n.Table).Run()
	exec.Command("nft", "add", "set", "inet", n.Table, n.Set, "{ type ipv4_addr; flags interval; }").Run()
	exec.Command("nft", "add", "chain", "inet", n.Table, "filter_hook", "{ type filter hook input priority -100; }").Run()
    
    // Flush chain to prevent rule duplication on restart
    exec.Command("nft", "flush", "chain", "inet", n.Table, "filter_hook").Run()
	exec.Command("nft", "add", "rule", "inet", n.Table, "filter_hook", "ip", "saddr", "@blocklist", "drop").Run()
}

func (n *NftablesEnforcer) Block(ip string) error {
	out, err := exec.Command("nft", "add", "element", "inet", n.Table, n.Set, "{", ip, "}").CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "File exists") { return nil }
		return fmt.Errorf("nft block failed: %s", string(out))
	}
	return nil
}
