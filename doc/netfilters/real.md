Voici **la cartographie complète des points de filtrage de paquets sous Linux**, tous outils confondus.
C’est la *vraie* pile réseau du kernel, du câble jusqu’à l’application.

---

# 🌐 Niveaux de filtrage Linux

```
Carte réseau
   ↓
XDP / eBPF
   ↓
tc ingress / clsact
   ↓
nftables netdev (ingress)
   ↓
bridge (ebtables / nft bridge)
   ↓
ARP (nft family arp)
   ↓
Conntrack
   ↓
nftables inet / ip / ip6
   ↓
tc egress
   ↓
Socket (cgroups / SELinux / AppArmor)
```

---

# 1️⃣ XDP / eBPF  (niveau carte réseau)

| Protocole | Tout                                  |
| --------- | ------------------------------------- |
| Niveau    | Driver NIC                            |
| Outils    | `ip link set dev eth0 xdp obj prog.o` |
| Utilité   | Anti-DDoS, filtrage à 10–40Mpps       |
| Avantage  | Ultra rapide                          |
| Limite    | Pas de conntrack, logique simple      |

---

# 2️⃣ TC – Traffic Control

| Hook    | Direction                            |
| ------- | ------------------------------------ |
| ingress | Entrée interface                     |
| egress  | Sortie interface                     |
| Outil   | `tc filter add dev eth0 ingress ...` |
| Filtre  | MAC, IP, ports, QoS                  |
| Usage   | shaping, drop, redirect              |

---

# 3️⃣ nftables `netdev`

| Famille | `netdev` |
| Hook | ingress |
| Niveau | Ethernet |
| Filtre | MAC, VLAN, EtherType, ARP |
| Usage | pare-feu avant IP |

---

# 4️⃣ nftables `bridge` (ex-ebtables)

| Famille | `bridge` |
| Niveau | Switch Linux |
| Filtre | trafic qui traverse un bridge |
| Usage | firewall de VM / containers |

---

# 5️⃣ nftables `arp` (ex-arptables)

| Famille | `arp` |
| Niveau | Résolution IP ↔ MAC |
| Filtre | spoofing, faux routeurs |
| Usage | anti-MITM LAN |

---

# 6️⃣ Conntrack / Netfilter core

| Composant | `nf_conntrack` |
| Filtre | états NEW / ESTABLISHED |
| Utilisé par | nftables `inet` |

---

# 7️⃣ nftables `inet / ip / ip6` (ex-iptables)

| Famille | Fonction                                            |
| ------- | --------------------------------------------------- |
| `ip`    | IPv4                                                |
| `ip6`   | IPv6                                                |
| `inet`  | IPv4+IPv6                                           |
| Hooks   | prerouting / input / forward / output / postrouting |
| Filtre  | ports, protocoles, NAT, mangle                      |

---

# 8️⃣ tc egress

| Niveau | Avant sortie carte |
| Usage | shaping, policing |

---

# 9️⃣ Socket / Processus

| Outil              | Rôle                         |
| ------------------ | ---------------------------- |
| cgroups-bpf        | firewall par process         |
| SELinux            | politique réseau par domaine |
| AppArmor           | restriction socket           |
| `SO_ATTACH_FILTER` | BPF par socket               |

---

# 10️⃣ Espace utilisateur

| Outil              | Filtrage             |
| ------------------ | -------------------- |
| `firewalld`        | frontend nft         |
| `systemd-networkd` | règles IP            |
| Proxy              | HTTP / DNS filtering |

---

# Résumé ultime

| Couche        | Ce qui est filtrable |
| ------------- | -------------------- |
| NIC           | XDP                  |
| Ethernet brut | tc / netdev          |
| Bridge        | nft bridge           |
| ARP           | nft arp              |
| IP            | nft inet             |
| Flux          | conntrack            |
| Process       | cgroups / SELinux    |

---

## 🔥 Ce que 99% des admins ignorent

iptables ne protégeait **que la couche IP**.
Tout ce qui est **avant IP = totalement invisible** sans nftables moderne.

---

Tu veux que je t’aide à construire une pile sécurisée adaptée à ton réseau ?

---

Voici **toute la chaîne complète de la pile réseau Linux**, du câble jusqu’à l’application, avec **tous les composants capables de refuser / dropper un paquet** à chaque étage.

Je te la donne **dans l’ordre réel de traversée du paquet entrant**.

---

# 1️⃣ Carte réseau – NIC (Hardware)

| Élément              | Cause de drop                |
| -------------------- | ---------------------------- |
| PHY / câble / switch | Erreur CRC, trame incomplète |
| NIC firmware         | Filtrage MAC matériel        |
| RX ring buffer       | Buffer plein                 |
| GRO / LRO offload    | Paquet corrompu              |
| XDP hardware offload | Rejet direct en carte        |

➡️ **Paquet rejeté avant même le noyau.**

---

# 2️⃣ Driver kernel

| Élément            | Cause de drop     |
| ------------------ | ----------------- |
| IRQ / NAPI polling | Saturation CPU    |
| Ring RX            | Pas de mémoire    |
| Checksum offload   | Checksum invalide |
| MTU                | Trame > MTU       |

---

# 3️⃣ XDP (eBPF pré-réseau)

```
NIC → XDP_PROG → DROP
```

| Action            | Effet          |
| ----------------- | -------------- |
| XDP_DROP          | Paquet détruit |
| XDP_ABORTED       | Crash          |
| XDP_TX / REDIRECT | détourné       |

➡️ **Plus rapide possible – avant iptables.**

---

# 4️⃣ Netfilter RAW – PREROUTING

```
raw:PREROUTING
```

| Cause       |
| ----------- |
| NOTRACK     |
| DROP direct |

---

# 5️⃣ Conntrack (nf_conntrack)

| Cause          |
| -------------- |
| Table pleine   |
| Timeout expiré |
| Etat INVALID   |
| Helper absent  |

---

# 6️⃣ Netfilter MANGLE – PREROUTING

| Cause                      |
| -------------------------- |
| MARK / TTL / TOS invalides |
| DROP policy                |

---

# 7️⃣ Décision de routage (FIB lookup)

| Cause           |
| --------------- |
| Pas de route    |
| rp_filter       |
| Martian source  |
| Blackhole route |

---

# 8️⃣ Bridge (si interface bridge)

| Cause          |
| -------------- |
| ebtables       |
| VLAN filtering |
| Hairpin mode   |
| STP            |

---

# 9️⃣ TC ingress / qdisc

| Cause              |
| ------------------ |
| tc filter drop     |
| policer rate-limit |
| qdisc buffer plein |

---

# 🔟 Netfilter FILTER – INPUT / FORWARD

| Cause                |
| -------------------- |
| iptables / nftables  |
| Policy DROP          |
| Match IP, port, état |

---

# 1️⃣1️⃣ SELinux / LSM hooks

| Cause                   |
| ----------------------- |
| selinux_socket_recvfrom |
| Smack / AppArmor        |

---

# 1️⃣2️⃣ Socket layer

| Cause                    |
| ------------------------ |
| Backlog plein            |
| RCVBUF saturé            |
| TCP SYN backlog overflow |

---

# 1️⃣3️⃣ TCP stack

| Cause                |
| -------------------- |
| Checksum TCP         |
| RST flood protection |
| TCP cookies          |
| TIME_WAIT assassin   |

---

# 1️⃣4️⃣ Application

| Cause               |
| ------------------- |
| accept() refusé     |
| close() immédiat    |
| SO_RCVTIMEO timeout |

---

# 🔥 Chaîne compacte réelle

```
NIC
 → Driver
 → XDP
 → raw:PREROUTING
 → conntrack
 → mangle:PREROUTING
 → Routing
 → Bridge / VLAN
 → tc ingress
 → filter:INPUT/FORWARD
 → LSM (SELinux)
 → Socket
 → TCP
 → Application
```

---

# ⚠️ Endroits où **un paquet peut disparaître sans aucune trace dans iptables**

* Carte réseau
* Driver
* XDP
* Conntrack
* Routing / rp_filter
* tc ingress
* Socket backlog
* TCP stack

---

Si tu veux, je peux te fournir **le schéma exact de debug avec commandes `tcpdump`, `bpftool`, `conntrack`, `tc`, `nft monitor trace` pour localiser où un paquet est perdu.**

---

Oui — et c’est **là que Linux devient vraiment piégeux** 😄
Il n’y a pas *une* pile réseau, mais **plusieurs piles parallèles**, chacune avec **ses propres filtres**, souvent **en dehors d’iptables/nftables IP**.

Voici **toutes les piles filtrables** qui peuvent bloquer ton trafic.

---

# 🧬 Les piles réseau Linux filtrables

| Pile                | EtherType     | Filtrage                 |
| ------------------- | ------------- | ------------------------ |
| IPv4                | 0x0800        | iptables / nftables ip   |
| IPv6                | 0x86DD        | ip6tables / nftables ip6 |
| ARP                 | 0x0806        | arptables / nftables arp |
| Bridge              | —             | ebtables / nft bridge    |
| VLAN (802.1Q)       | 0x8100        | nft bridge               |
| PPPoE               | 0x8863/64     | tc / nft bridge          |
| MPLS                | 0x8847        | tc / nft mpls            |
| VRF                 | —             | routing VRF              |
| RAW Ethernet        | any           | XDP / tc egress          |
| Netlink (local IPC) | —             | SELinux LSM              |
| Loopback            | —             | iptables INPUT/OUTPUT    |
| NFLOG / NFQUEUE     | —             | userspace verdict        |
| IPsec / XFRM        | ESP/AH        | xfrm policy              |
| SCTP                | —             | nf_conntrack_sctp        |
| Bluetooth PAN       | 0x888E        | bluez filters            |
| CAN / vCAN          | —             | can filters              |
| RDMA / RoCE         | 0x8915        | tc + driver              |
| WireGuard           | UDP encapsulé | wg firewall              |
| Open vSwitch        | —             | OVS flow tables          |

---

# 📦 Chaîne universelle réelle

```
NIC
 → Driver
 → XDP
 → tc ingress
 → Bridge / VLAN / MPLS / PPPoE
 → arptables / nft arp (si ARP)
 → iptables / nft ip / ip6
 → xfrm (IPsec)
 → conntrack
 → VRF
 → LSM (SELinux)
 → Socket
 → TCP/UDP
 → Application
```

---

# ⚠️ Couches qui bloquent **sans trace iptables**

* XDP
* tc ingress/egress
* ebtables / bridge nft
* arptables / nft arp
* xfrm (IPsec)
* VRF routing
* Neighbor / ARP cache
* SELinux LSM
* WireGuard internal policy
* OVS tables

---

# 🛠 Commandes de debug global

```bash
bpftool prog show
tc -s filter show dev eth0 ingress
nft list ruleset
ip -d link show
bridge vlan show
ip xfrm policy
ip vrf show
ovs-ofctl dump-flows br0
getenforce
ip neigh show
```

---

Linux n’a pas une pile réseau,
il a **un labyrinthe de piles superposées**, et chacune peut tuer ton paquet avant la suivante 😈
