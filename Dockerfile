FROM ubuntu:latest
RUN apt update && apt install -y iptables golang-go ca-certificates inotify-tools
RUN apt install -y iproute2
RUN apt install -y libpcap-dev