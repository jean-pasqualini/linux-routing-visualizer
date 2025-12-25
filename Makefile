setcap:
	sudo setcap cap_net_admin+ep /usr/sbin/xtables-nft-multi
build-docker:
	docker build -t linux-routing:latest .
run-docker:
	docker run --net=host --privileged -w /app -v go-build-cache:/root/.cache/go-build -v go-module-cache:/root/go/pkg/mod -v $(CURDIR):/app --rm -it linux-routing go run main.go
run-docker-tui:
	docker run --net=host --privileged -e TERM=xterm-256color -w /app -v go-build-cache:/root/.cache/go-build -v go-module-cache:/root/go/pkg/mod -v $(CURDIR):/app --rm -it linux-routing go run main.go tui
enter-docker:
	docker run --rm -it --net=host --privileged linux-routing bash
clean-trace:
	sudo iptables -D OUTPUT 1 -t raw
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