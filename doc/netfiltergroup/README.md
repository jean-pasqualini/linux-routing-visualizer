Netfilter est une infrastructure complète du noyau Linux, pas juste iptables.
Voici tout ce qui fait réellement partie de Netfilter.

⸻

1️⃣ Les hooks noyau (points d’interception)

C’est le cœur de Netfilter : les endroits où le noyau peut intercepter un paquet.

Hook	Rôle
NF_INET_PRE_ROUTING	Avant décision de routage
NF_INET_LOCAL_IN	Paquet destiné à la machine
NF_INET_FORWARD	Paquet routé à travers la machine
NF_INET_LOCAL_OUT	Paquet émis par la machine
NF_INET_POST_ROUTING	Après routage


⸻

2️⃣ Les tables Netfilter

Les tables organisent les règles par fonction.

Table	Fonction
filter	Firewall (DROP / ACCEPT)
nat	NAT, DNAT, SNAT, MASQUERADE
mangle	Modification des paquets
raw	Bypass conntrack
security	SELinux / LSM


⸻

3️⃣ Le moteur conntrack (suivi d’état)

Suivi de toutes les connexions réseau.
•	États : NEW, ESTABLISHED, RELATED, INVALID
•	NAT dépend entièrement de conntrack.
•	IDS / firewall stateful impossible sans lui.

⸻

4️⃣ Les modules noyau Netfilter

Module	Rôle
nf_conntrack	Suivi de connexion
nf_nat	NAT
nf_defrag_ipv4/6	Réassemblage
nf_tables	Backend nftables
xt_conntrack	Match conntrack
xt_recent	Détection brute-force
xt_hashlimit	Rate limiting
nfnetlink_queue	Paquets vers userspace
nf_log	Logging


⸻

5️⃣ Sous-systèmes spécialisés

Couche	Élément
Ethernet	ebtables
ARP	arptables
Bridge	br_netfilter
IPv4/IPv6	iptables, nftables


⸻

6️⃣ Filesystem & interfaces noyau

Interface	Rôle
/proc/net/nf_conntrack	Table des connexions
/proc/sys/net/netfilter/*	Paramètres kernel
Netlink NETLINK_NETFILTER	Communication userspace


⸻

7️⃣ Front-ends userspace

Outil	Statut
iptables	Legacy
ip6tables	Legacy IPv6
arptables	Legacy ARP
ebtables	Legacy Ethernet
nft	Moderne unifié


⸻

Phrase clé

Netfilter n’est pas un firewall.
C’est une plateforme complète de contrôle du trafic réseau dans le noyau Linux.