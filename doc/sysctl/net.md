sysctl te permet de piloter le comportement interne du kernel réseau Linux.
C’est LE couteau suisse pour debug, perf, sécurité et routage.

Je te fais la carte complète utile pour ton projet linux-routing.

⸻

🌐 1️⃣ IP / ROUTING (net.ipv4.*)

Clé	Effet
net.ipv4.ip_forward	Active le routage L3
net.ipv4.conf.all.forwarding	Forward sur toutes les interfaces
net.ipv4.conf.eth0.forwarding	Forward par interface
net.ipv4.conf.all.rp_filter	Anti-spoofing
net.ipv4.conf.all.accept_redirects	Accepte ICMP redirect
net.ipv4.conf.all.send_redirects	Envoie ICMP redirect
net.ipv4.conf.all.proxy_arp	Active Proxy-ARP
net.ipv4.conf.all.accept_source_route	Source routing
net.ipv4.conf.all.log_martians	Log paquets suspects
net.ipv4.route.flush	Flush cache de routes


⸻

🧱 2️⃣ CONNTRACK / NETFILTER

Clé	Effet
net.netfilter.nf_conntrack_max	Taille table conntrack
net.netfilter.nf_conntrack_tcp_timeout_established	Timeout ESTABLISHED
net.netfilter.nf_conntrack_generic_timeout	Timeout générique
net.netfilter.nf_conntrack_buckets	Hash buckets
net.netfilter.nf_log_all_netns	Log iptables dans tous les netns
net.netfilter.nf_conntrack_checksum	Checksum validation


⸻

🔥 3️⃣ TCP STACK

Clé	Effet
net.ipv4.tcp_syncookies	Protection SYN flood
net.ipv4.tcp_tw_reuse	Réutilisation TIME_WAIT
net.ipv4.tcp_fin_timeout	Timeout FIN_WAIT
net.ipv4.tcp_max_syn_backlog	Queue SYN
net.ipv4.tcp_keepalive_time	Keepalive delay
net.ipv4.tcp_sack	Selective ACK
net.ipv4.tcp_timestamps	TCP timestamps
net.ipv4.tcp_congestion_control	Algo CC


⸻

📦 4️⃣ BUFFERS & PERFORMANCE

Clé	Effet
net.core.rmem_max	RX buffer max
net.core.wmem_max	TX buffer max
net.core.netdev_max_backlog	Queue NIC
net.ipv4.tcp_rmem	TCP RX min/def/max
net.ipv4.tcp_wmem	TCP TX min/def/max
net.ipv4.ip_local_port_range	Range ports éphémères


⸻

🕵️ 5️⃣ DEBUG / TRACE

Clé	Effet
net.netfilter.nf_log_all_netns	Trace inter-netns
net.ipv4.conf.all.log_martians	Log paquets invalides
kernel.printk	Verbosité kernel
net.ipv4.tcp_abort_on_overflow	Reset immédiat si backlog plein


⸻

🧪 6️⃣ MULTI-PATH / POLICY ROUTING

Clé	Effet
net.ipv4.fib_multipath_hash_policy	Hash ECMP
net.ipv4.fib_multipath_use_neigh	ECMP basé ARP
net.ipv4.conf.all.promote_secondaries	IP secondaire


⸻

🛡️ 7️⃣ SÉCURITÉ

Clé	Effet
net.ipv4.conf.all.accept_redirects	ICMP redirect
net.ipv4.conf.all.secure_redirects	Redirect sécurisé
net.ipv4.conf.all.send_redirects	Emission redirect
net.ipv4.conf.all.rp_filter	Anti-spoofing


⸻

⚡ Exemples pratiques

Activer routage L3 :

sysctl -w net.ipv4.ip_forward=1

Désactiver ICMP redirect :

sysctl -w net.ipv4.conf.all.accept_redirects=0
sysctl -w net.ipv4.conf.all.send_redirects=0

Augmenter conntrack :

sysctl -w net.netfilter.nf_conntrack_max=1048576

Debug netfilter multi-netns :

sysctl -w net.netfilter.nf_log_all_netns=1


⸻

Si tu veux, je peux te générer un profil sysctl spécial lab réseau / debug kernel / iptables / ipvs prêt à appliquer.