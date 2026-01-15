Netfilter n’est pas “iptables”.
**iptables n’est qu’une interface utilisateur.**
Netfilter est un **sous-système noyau complet** composé de plusieurs briques.

Voici l’architecture réelle.

---

# 🧱 1️⃣ Netfilter core (noyau Linux)

C’est le moteur bas niveau qui fournit les **hooks réseau** :

| Hook          | Rôle                          |
| ------------- | ----------------------------- |
| `PREROUTING`  | avant décision de routage     |
| `INPUT`       | paquets destinés à la machine |
| `FORWARD`     | paquets routés                |
| `OUTPUT`      | paquets émis localement       |
| `POSTROUTING` | juste avant sortie NIC        |

Modules :

```
nf_hook_slow
nf_defrag_ipv4 / ipv6
nfnetlink
```

---

# 🧱 2️⃣ Conntrack – suivi de connexion

Sous-système clé :

```
nf_conntrack
nf_conntrack_ipv4
nf_conntrack_ipv6
xt_conntrack
```

Fonctions :

* garde l’état des flux
* fournit : NEW, ESTABLISHED, RELATED, INVALID
* indispensable pour NAT, firewall stateful, FTP helpers, etc.

---

# 🧱 3️⃣ NAT engine

Modules :

```
nf_nat
nf_nat_ipv4
nf_nat_ipv6
xt_nat
xt_MASQUERADE
```

Rôles :

* DNAT
* SNAT
* MASQUERADE
* REDIRECT

---

# 🧱 4️⃣ Match / Target extensions (xtables)

Ce sont les briques utilisées par iptables.

Exemples :

| Type      | Module         |
| --------- | -------------- |
| TCP       | `xt_tcpudp`    |
| state     | `xt_conntrack` |
| limit     | `xt_limit`     |
| hashlimit | `xt_hashlimit` |
| owner     | `xt_owner`     |
| LOG       | `xt_LOG`       |
| DROP      | `xt_standard`  |

---

# 🧱 5️⃣ Tables Netfilter

| Table      | Rôle             |
| ---------- | ---------------- |
| `raw`      | bypass conntrack |
| `mangle`   | modifier paquets |
| `nat`      | NAT              |
| `filter`   | firewall         |
| `security` | SELinux          |

---

# 🧱 6️⃣ Interface utilisateur

Ce n’est PAS Netfilter, juste une API :

| Outil       | Rôle            |
| ----------- | --------------- |
| `iptables`  | IPv4 legacy     |
| `ip6tables` | IPv6 legacy     |
| `ebtables`  | bridge          |
| `arptables` | ARP             |
| `nft`       | backend moderne |

---

# 🧱 7️⃣ nftables (nouvelle génération)

Modules kernel :

```
nf_tables
nf_tables_ipv4
nf_tables_ipv6
```

Remplace progressivement :

```
iptables + ebtables + arptables
```

---

# 🧠 Vue d’ensemble

```
┌────────────┐
│ Userspace  │  nft / iptables
└──────┬─────┘
       ↓
┌──────────────────────┐
│ Netfilter Core       │  hooks, routing, dispatch
└─────────┬────────────┘
          ↓
┌──────────────────────┐
│ Conntrack / NAT      │  nf_conntrack, nf_nat
└─────────┬────────────┘
          ↓
┌──────────────────────┐
│ xtables / nftables   │  match + targets
└──────────────────────┘
```

---

# 🏁 En une phrase

Netfilter est un **framework noyau modulaire** composé de :

* hooks réseau
* suivi de connexions
* moteur NAT
* tables de filtrage
* modules match/target
* interfaces utilisateurs (iptables / nft)

iptables n’est que la **télécommande**. Netfilter est le moteur du firewall Linux.
