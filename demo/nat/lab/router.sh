#!/bin/bash
apt update
apt install -y nftables conntrack iproute2

ip -br link
ip link set eth0 down

ip addr add 10.0.0.1/24 dev eth1
ip addr add 192.168.0.1/24 dev eth2
ip link set eth1 up
ip link set eth2 up
sysctl -w net.ipv4.ip_forward=1

nft add table ip nat
nft add chain ip nat postrouting { type nat hook postrouting priority 100 \; }
nft add rule ip nat postrouting oifname "eth2" ip saddr 10.0.0.0/24 masquerade

nft add table ip filter
nft add chain ip filter forward { type filter hook forward priority 0 \; policy drop \; }
nft add rule ip filter forward iifname "eth1" oifname "eth2" accept
nft add rule ip filter forward ct state established,related accept

sleep infinity