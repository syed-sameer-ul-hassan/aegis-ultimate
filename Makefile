INSTALL_ROOT=/opt/aegis
ML_ENV=$(INSTALL_ROOT)/ml_env

.PHONY: all build install clean

all:
	@echo " -> Installing bpf2go tool..."
	go install github.com/cilium/ebpf/cmd/bpf2go@latest
	@echo " -> Generating BPF bindings..."
	cd cmd/aegisd && go generate
	@echo " -> Building Go Daemon..."
	cd cmd/aegisd && go build -o ../../bin/aegisd

install:
	@echo " -> Creating Directories..."
	mkdir -p $(INSTALL_ROOT)/ml
	mkdir -p /etc/aegis /var/lib/aegis /var/log/aegis
	
	@echo " -> Installing Binaries..."
	cp bin/aegisd /usr/local/bin/
	
	@echo " -> Setting up ML Sandbox..."
	python3 -m venv $(ML_ENV)
	$(ML_ENV)/bin/pip install -r ml/requirements.txt
	cp ml/scorer.py $(INSTALL_ROOT)/ml/
	
	@echo " -> Installing Systemd Service..."
	cp deploy/aegisd.service /etc/systemd/system/
	systemctl daemon-reload
	@echo "✅ Install Complete."

clean:
	rm -rf bin/ cmd/aegisd/bpf_*.go cmd/aegisd/bpf_*.o
