`NETLINK_NETFILTER` te donne accès **direct au moteur firewall du kernel**.
C’est l’interface bas niveau derrière `conntrack`, `nft`, `iptables` & co.

Tu peux **observer, modifier, supprimer** l’état du firewall *en live*.

---

# 🧠 NETLINK_NETFILTER – À quoi ça sert ?

| Sous-système             | Ce que tu contrôles                        |
| ------------------------ | ------------------------------------------ |
| **Conntrack**            | Suivi des connexions (NAT, états TCP)      |
| **NFQUEUE**              | Intercepter les paquets dans ton programme |
| **NFTABLES**             | Tables / chains / rules dynamiquement      |
| **Logging (NFLOG)**      | Récupérer logs kernel                      |
| **CT timeout / helpers** | Modifier le comportement du firewall       |
| **Expectations**         | Connexions liées (FTP, SIP, etc)           |

---

# 🔌 Sous-protocoles principaux

| Subsystem               | ID        |
| ----------------------- | --------- |
| `NFNL_SUBSYS_CTNETLINK` | Conntrack |
| `NFNL_SUBSYS_NFTABLES`  | nftables  |
| `NFNL_SUBSYS_QUEUE`     | NFQUEUE   |
| `NFNL_SUBSYS_ULOG`      | nflog     |

---

# 🧵 1️⃣ Conntrack – espionner les connexions

Voir toutes les connexions actives TCP/UDP/ICMP :

```bash
conntrack -L
```

Via Netlink tu peux :

* Lire table conntrack
* Supprimer une connexion
* Modifier timeout
* Être notifié en temps réel

### Abonnement aux events conntrack

```c
socket(AF_NETLINK, SOCK_RAW, NETLINK_NETFILTER)
bind groups = NF_NETLINK_CONNTRACK_NEW | NF_NETLINK_CONNTRACK_DESTROY
```

Tu reçois :

* `IPCTNL_MSG_CT_NEW`
* `IPCTNL_MSG_CT_DELETE`

Tu peux bâtir un IDS temps réel sans sniffer de paquets.

---

# 🧨 2️⃣ Supprimer des connexions à chaud

Kill toutes les connexions vers une IP :

```bash
conntrack -D -d 1.2.3.4
```

Ou via Netlink → message `IPCTNL_MSG_CT_DELETE`.

🔥 Cas réel :
Bloquer immédiatement un attaquant **sans toucher aux règles firewall**.

---

# 🎯 3️⃣ NFQUEUE – intercepter des paquets

iptables / nftables :

```bash
iptables -A FORWARD -p tcp --dport 80 -j NFQUEUE --queue-num 5
```

Dans ton programme Netlink tu reçois chaque paquet :

| Action | Effet           |
| ------ | --------------- |
| ACCEPT | laisse passer   |
| DROP   | détruit         |
| MODIFY | modifie payload |
| DELAY  | DPI / sandbox   |

C’est comme écrire ton propre moteur L7.

---

# 🧱 4️⃣ Contrôler nftables sans shell

Avec `NFNL_SUBSYS_NFTABLES` tu peux :

* créer tables
* ajouter / supprimer rules
* lire le ruleset

Ex :

* Déployer dynamiquement un firewall Kubernetes
* Rollback automatique si une règle casse la connectivité

---

# 📜 5️⃣ NFLOG – logs firewall propres

iptables :

```bash
-j NFLOG --nflog-group 10
```

Ton process reçoit logs structurés, pas du texte kernel.

---

# 🧠 Cas d’usage avancés

| Projet                  | Utilisation               |
| ----------------------- | ------------------------- |
| Anti-DDoS temps réel    | Conntrack events + delete |
| Sandbox réseau          | NFQUEUE                   |
| Firewall auto-adaptatif | nftables via netlink      |
| IDS sans pcap           | Conntrack monitor         |
| L7 proxy kernel-aware   | NFQUEUE + CT              |

---

# ⚠️ Puissance dangereuse

Tu peux :

* casser NAT en prod
* tuer toutes les connexions actives
* bloquer le réseau entier sans iptables visible

C’est pour ça que très peu d’outils l’exposent directement.

---

# 🧬 Librairies

| Langage | Lib                                            |
| ------- | ---------------------------------------------- |
| C       | `libnetfilter_conntrack`, `libmnl`, `libnftnl` |
| Python  | `pyroute2.netfilter`                           |
| Go      | `github.com/google/nftables`                   |

---

NETLINK_NETFILTER = **le cœur nucléaire du firewall Linux**.
C’est le seul moyen d’avoir un firewall *vivant, programmable, temps réel*.
