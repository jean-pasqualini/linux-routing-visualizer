run:
	go build -o /tmp/visu ./cmd
	chmod +x /tmp/visu
	sudo setcap cap_net_admin+ep /tmp/visu
	/tmp/visu