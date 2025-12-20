#!/bin/bash
# Simulates attacks against localhost to test AEGIS-X
# REQUIRES: hping3 (sudo apt install hping3)

TARGET="127.0.0.1"

echo "========================================"
echo "🛡️  AEGIS-X PENETRATION TEST SUITE"
echo "========================================"

# Check if aegis is running
if ! pgrep -x "aegisd" > /dev/null; then
    echo "❌ ERROR: Aegis Daemon is NOT running."
    exit 1
fi

echo "[1/3] Testing SYN Flood Detection..."
# Send 5000 SYN packets as fast as possible
sudo hping3 -S -p 80 --flood -c 5000 $TARGET 2>/dev/null &
PID=$!
echo "      Attack running (PID $PID)... waiting 3s..."
sleep 3
kill $PID 2>/dev/null

# Check if we got banned
if sudo nft list set inet aegis blocklist | grep -q "$TARGET"; then
    echo "✅ SUCCESS: Localhost was banned by Firewall."
    # Cleanup self-ban so we can continue testing
    sudo nft delete element inet aegis blocklist { $TARGET }
else
    echo "❌ FAIL: No ban detected."
fi

echo "========================================"
echo "Test Complete."
