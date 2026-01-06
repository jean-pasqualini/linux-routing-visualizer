LEVEL = medium

DOCKER_RUN = docker run -e TERM=xterm-256color --net=host --pid=host --privileged \
	-w /app \
	-v go-build-cache:/root/.cache/go-build \
	-v go-module-cache:/root/go/pkg/mod \
	-v $(CURDIR):/app \
	--rm -it linux-routing

DOCKER_RUN_NO_CGO = docker run -e TERM=xterm-256color -e CGO_ENABLED=0 --net=host --pid=host --privileged \
	-w /app \
	-v go-build-cache:/root/.cache/go-build \
	-v go-module-cache:/root/go/pkg/mod \
	-v $(CURDIR):/app \
	--rm -it linux-routing

DOCKER_RUN_DETACHED = docker run -d -e TERM=xterm-256color --net=host --privileged \
	-w /app \
	-v go-build-cache:/root/.cache/go-build \
	-v go-module-cache:/root/go/pkg/mod \
	-v $(CURDIR):/app \
	--rm -it linux-routing

setcap:
	sudo setcap cap_net_admin+ep /usr/sbin/xtables-nft-multi
build-docker:
	docker build -t linux-routing:latest .
run-docker-list:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go
run-docker-sniff:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go libpcap-dev -i eth0
run-docker-socket:
	$(DOCKER_RUN_NO_CGO) go run -tags $(LEVEL) main.go socket
run-docker-iptable:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go iptable
run-docker-link:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go link
run-docker-ipvs:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go ipvs
run-docker-sysctl:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go sysctl
run-docker-routing:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go routing
run-docker-bubble:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go bubble
run-docker-simulate:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go simulate
run-docker-tui:
	$(DOCKER_RUN) go run -tags $(LEVEL) main.go tui
dev-docker-tui:
	$(DOCKER_RUN_DETACHED) /app/kill.sh
	$(DOCKER_RUN) /app/dev.sh
enter-docker:
	docker run --rm -it --net=host --privileged linux-routing bash
clean-trace:
	sudo iptables -D OUTPUT 1 -t raw
complete:
	iptables -I PREROUTING 1 -p tcp --sport 8888 --dport 9090 -t raw -j TRACE
	iptables -I PREROUTING 2 -t raw -j TRACE
trace:
	sudo iptables -I OUTPUT 1 -m tcp -p tcp --dport 9090 -t raw -j TRACE
	sudo iptables -I PREROUTING 1 -m tcp -p tcp --dport 9090 -t raw -j TRACE
list-trace:
	sudo iptables -L OUTPUT --line-numbers -t raw
monitor:
	sudo nft monitor trace
listen:
	sudo nc -l -p 9090
kernel-log:
	sudo dmesg -w # can be journalctl -kf
install-conntrack:
	sudo apt install conntrack
delete-entry:
	sudo conntrack -D -p tcp --dport 9090
ip:
	hostname -I | awk '{ print $$1}'
from_all_ns:
	sysctl -w net.netfilter.nf_log_all_netns=1
list-rules:
	sudo nft list ruleset
nsenter:
	docker run -it --rm --privileged --pid=host justincormack/nsenter1
curl:
	curl 192.168.32.213:9090
list-tables:
	ip route show table local
	ip route show table main
	ip route show table 2
man:
	docker build -t linux-man -f Dockerfile.man .
	docker run --rm -it linux-man man iptables
ipvs-create:
	sudo ipvsadm -A -t 1.1.1.1:8080 -s rr