#!/bin/bash
apt update
apt install -y iproute2 netcat-openbsd

ip -br link
ip link set eth0 down

ip addr add 10.0.0.2/24 dev eth1
ip link set eth1 up

ip route add 192.168.0.2/32 via 10.0.0.1

echo "I'm your worst nightmare" | nc 192.168.0.2 8080

sleep infinity
