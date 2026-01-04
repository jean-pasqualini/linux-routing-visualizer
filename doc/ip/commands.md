Toutes ces valeurs correspondent aux **objets Netlink du kernel Linux** que la commande `ip` sait manipuler.
Chaque mot est un *sous-protocole netlink* différent.

En clair : `ip` n’est qu’un client générique pour **toute la pile réseau interne du kernel**.

---

## Carte mentale

| Domaine               | Objet                                               |
| --------------------- | --------------------------------------------------- |
| Interfaces            | `link`, `tuntap`, `tap`, `vrf`, `macsec`            |
| Adresses              | `address`, `addrlabel`, `maddress`, `ntable`        |
| Voisinage / ARP / ND  | `neighbor`, `neighbour`, `ntbl`                     |
| Routage               | `route`, `rule`, `nexthop`, `mroute`, `mrule`, `sr` |
| Tunnels               | `tunnel`, `fou`, `ila`, `l2tp`, `xfrm`              |
| Multipath / HA        | `mptcp`, `tcpmetrics`                               |
| Virtualisation réseau | `netns`, `netconf`, `vrf`                           |
| Monitoring            | `monitor`                                           |
| Autres noyau          | `amt`, `ioam`, `token`, `help`                      |

---

## Description des principaux

| Objet                    | Sert à quoi                                   |
| ------------------------ | --------------------------------------------- |
| `link`                   | Interfaces réseau (eth0, veth, bridge, vlan…) |
| `address`                | IP attachées aux interfaces                   |
| `addrlabel`              | Labels IPv6 (policy routing IPv6)             |
| `neighbor` / `neighbour` | Table ARP / NDP                               |
| `route`                  | Table de routage principale                   |
| `rule`                   | Policy routing (ip rule)                      |
| `nexthop`                | Objets next-hop partagés (ECMP avancé)        |
| `mroute` / `mrule`       | Routage multicast                             |
| `sr`                     | Segment Routing IPv6                          |
| `tunnel`                 | Tunnels IPIP / GRE / SIT                      |
| `fou`                    | Foo-over-UDP (tunnels encapsulés UDP)         |
| `ila`                    | Identifier Locator Addressing (IPv6 mobilité) |
| `l2tp`                   | Tunnels L2TP                                  |
| `xfrm`                   | IPsec (Security Associations)                 |
| `vrf`                    | Virtual Routing & Forwarding                  |
| `macsec`                 | Sécurité Ethernet (802.1AE)                   |
| `tuntap` / `tap`         | Interfaces virtuelles TUN/TAP                 |
| `mptcp`                  | Multipath TCP                                 |
| `tcpmetrics`             | Cache RTT / cwnd TCP                          |
| `netns`                  | Network namespaces                            |
| `netconf`                | Paramètres IPv4/IPv6 par interface            |
| `monitor`                | Écoute événements netlink en live             |
| `maddress`               | Adresses multicast                            |
| `ntable` / `ntbl`        | Tables neighbor avancées                      |
| `token`                  | Génération IPv6 tokens                        |
| `amt`                    | Automatic Multicast Tunneling                 |
| `ioam`                   | In-situ OAM (télémétrie réseau)               |

---

## Ce que tu utilises 95% du temps

```bash
ip link
ip addr
ip route
ip rule
ip neigh
ip netns
```

Tout le reste ce sont des **API kernel ultra-spécialisées** que `ip` expose mais que même beaucoup d’admins ne touchent jamais.

---

## En résumé

`ip` est un **client universel Netlink**.
Chaque sous-commande correspond directement à un **objet kernel réseau interne**.
Tu vois littéralement l’API privée du kernel exposée en CLI.
