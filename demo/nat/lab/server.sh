#!/bin/bash
apt update
apt install -y iproute2 netcat-openbsd

ip link set eth0 down

ip addr add 192.168.0.2/24 dev eth1
ip link set eth1 up

ip route add 192.168.0.1/32 dev eth1
nc -l 8080

sleep infinity