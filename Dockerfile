FROM ubuntu:latest
RUN apt update && apt install -y iptables golang-go ca-certificates