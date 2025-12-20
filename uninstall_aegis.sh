#!/bin/bash
# Full Cleanup Script

if [[ $EUID -ne 0 ]]; then
   echo "Run as root."
   exit 1
fi

echo "[*] Stopping Service..."
systemctl stop aegisd
systemctl disable aegisd

echo "[*] Removing Binary and Scripts..."
rm -f /usr/local/bin/aegisd
rm -f /usr/local/bin/aegis-ctl
rm -f /etc/systemd/system/aegisd.service

echo "[*] Cleaning Source Code..."
# Tries to find where you installed it. 
# WARNING: Adjust path if you installed somewhere else.
if [ -d "$HOME/Desktop/aegis-ultimate" ]; then
    rm -rf "$HOME/Desktop/aegis-ultimate"
fi
rm -rf /opt/aegis
rm -rf /var/lib/aegis
rm -rf /etc/aegis

echo "[*] Flushing Firewall Rules..."
nft delete table inet aegis 2>/dev/null

echo "[*] Reloading Systemd..."
systemctl daemon-reload

echo "✅ AEGIS-X has been removed."
