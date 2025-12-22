setcap:
	sudo setcap cap_net_admin+ep /usr/sbin/xtables-nft-multi
build-docker:
	docker build -t linux-routing:latest .
run-docker:
	docker run --net=host --privileged -w /app -v go-build-cache:/root/.cache/go-build -v go-module-cache:/root/go/pkg/mod -v $(CURDIR):/app --rm -it linux-routing go run main.go
enter-docker:
	docker run --rm -it --net=host --privileged linux-routing bash