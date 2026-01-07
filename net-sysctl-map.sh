#!/bin/sh

print_group() {
  TITLE="$1"
  REGEX="$2"

  echo
  echo "══════════════════ $TITLE ══════════════════"
  sysctl -a 2>/dev/null | grep -E "$REGEX" | grep -v -E "bond0|tunl0|teql0|sit0|services1|lo|ip_vti0|ip6tnl0|ip6_vti0|eth0|erspan0|dummy0|gretap0|gre0|ip6gre0|docker0|all" || true
}

print_group ROUTING     '^net.ipv4.ip_forward|^net.ipv4.conf'
print_group TCP         '^net.ipv4.tcp_'
print_group CONNTRACK   '^net.netfilter.nf_conntrack'
print_group IPVS        '^net.ipv4.vs'
print_group NETFILTER   '^net.netfilter'
print_group BUFFER      '^net.core.rmem|^net.core.wmem|^net.core.netdev'
print_group DEBUG       'log_martians|printk|nf_log'
print_group MULTIPATH   'fib_multipath'
print_group BRIDGE      '^net.bridge'
print_group MPTCP       '^net.mptcp'
print_group IPV6        '^net.ipv6'