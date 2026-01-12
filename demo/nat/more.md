Là tu touches au cœur du système réseau Linux.
Je te fais la **carte mentale complète du “routeur Linux”** — celle que les devs kernel ont en tête.

---

# 🧠 Architecture réseau Linux — Big Picture

```
App
 │
Socket API
 │
TCP / UDP / ICMP
 │
Routing (FIB)
 │
Netfilter hooks ──▶ nftables / iptables
 │                   ├─ DNAT / SNAT
 │                   └─ Firewall
 │
Conntrack (state table)
 │
Device layer (veth, eth0, bridge)
 │
Driver / NIC
```

---

# 🧱 Rôle de chaque composant

| Composant                      | Rôle réel                              |
| ------------------------------ | -------------------------------------- |
| **Namespace réseau**           | Crée un kernel réseau isolé            |
| **Interface (eth0, veth, lo)** | Point d’entrée/sortie L2               |
| **Routing table (ip route)**   | Choisit l’interface de sortie          |
| **FIB / RIB**                  | Structure de données kernel du routage |
| **Netfilter hooks**            | Points d’interception des paquets      |
| **nftables / iptables**        | Règles programmées par l’admin         |
| **Conntrack**                  | Mémoire du trafic (stateful)           |
| **SNAT / DNAT / MASQUERADE**   | Réécriture d’adresses                  |
| **FORWARD / INPUT / OUTPUT**   | Filtrage L3/L4                         |
| **ARP / NDP**                  | Résolution MAC                         |
| **veth**                       | Câble Ethernet virtuel                 |
| **Bridge (br0)**               | Switch L2                              |
| **rp_filter**                  | Anti-spoofing / anti-asymétrie         |
| **TCP stack**                  | Fiabilité, ports, retransmissions      |
| **NIC / driver**               | Envoi réel sur le fil                  |

---

# 🧠 Le NAT vu par le kernel

Quand un paquet traverse le routeur :

1. **Routing lookup**
2. **Netfilter PREROUTING**
3. **DNAT éventuel**
4. **Routing relookup**
5. **FORWARD chain**
6. **POSTROUTING**
7. **SNAT / MASQUERADE**
8. **Émission**

Et le retour fait le chemin inverse grâce à conntrack.

---

# 🧩 Conclusion

Linux n’a pas de “fonction NAT”.

Il a une **chaîne d’organes spécialisés** qui ensemble forment une box Internet logicielle.
