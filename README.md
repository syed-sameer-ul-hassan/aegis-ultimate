# 🛡️ AEGIS-X ULTIMATE (v8.0)



## 👤 Author & Credits
**Architect:** Noor
**Role:** Principal Blue Team Engineer

---

## ☕ Support the Development
If you like this tool and want to support me for more tools like this, **buy me a coffee bag** so I can write for all of you!

**Bitcoin (BTC) Address:**
`bc1qkp3q2vwl6lav3sczu4x5hqfks6xfpkpyagnvrk`

---

## 🏗️ Architecture
1.  **Kernel Plane (eBPF):** RingBuffer event streaming & O(1) blocking at driver level.
2.  **Control Plane (Go):** Multi-factor threat scoring engine.
3.  **Enforcement (Nftables):** Atomic firewall set updates.
4.  **Intelligence (ML):** Python-based anomaly detection sidecar (Isolation Forest).

## 🚀 Installation & Usage

### Prerequisites
* Linux Kernel 5.8+
* Go 1.21+, Clang, Nftables, Python3

### 1. Build
```bash
cd aegis-ultimate
# Installs dependencies and compiles the project
make all

```

### 2. Install

```bash
# Sets up Systemd, ML VirtualEnv, and Configs
sudo make install

```

### 3. Run

```bash
sudo systemctl start aegisd
sudo systemctl status aegisd

```
