Très bonne question — et c’est exactement ce qu’on apprend **après s’être fait casser la gueule par IPVS** 😄

Je te fais la vraie cartographie du kernel netlink en mode *“où tu peux aller les yeux fermés, où tu dois mettre un casque”*.

---

# 🟢 Zones solides / fiables

Ces API ont une **ABI stable**, utilisées par iproute2, systemd, nftables, ss, NetworkManager, etc.

| Sous-système                      | Famille                          | Fiabilité     |
| --------------------------------- | -------------------------------- | ------------- |
| Interfaces / routes               | `NETLINK_ROUTE` (rtnetlink)      | 🟢 béton      |
| TCP / UDP sockets                 | `NETLINK_INET_DIAG`, `SOCK_DIAG` | 🟢 béton      |
| Conntrack / NAT                   | `NETLINK_NETFILTER`              | 🟢 béton      |
| nftables                          | `NETLINK_NETFILTER`              | 🟢 béton      |
| tc / qdisc / cls                  | `NETLINK_ROUTE`                  | 🟢 béton      |
| XDP / BPF maps                    | `NETLINK_ROUTE` / `bpf()`        | 🟢 béton      |
| Generic netlink familles modernes | `genetlink`                      | 🟢 en général |

👉 Si `ip`, `ss`, `tc`, `nft` utilisent ce canal → c’est fiable.

---

# 🟡 Zones moyennes (ça marche mais attention)

| Sous-système          | Famille         | Problèmes                  |
| --------------------- | --------------- | -------------------------- |
| WireGuard             | generic netlink | versioning jeune           |
| Team / bonding avancé | generic netlink | options kernel-dépendantes |
| nl80211 (WiFi)        | generic netlink | instable selon drivers     |
| RDMA / infiniband     | generic netlink | ABI évolutive              |

👉 Ça fonctionne, mais :

* attributs optionnels
* champs qui apparaissent / disparaissent
* parsing défensif obligatoire

---

# 🔴 Zones pourries / à éviter

| Sous-système                | Famille             | Pourquoi                        |
| --------------------------- | ------------------- | ------------------------------- |
| IPVS                        | generic netlink     | blobs de structs C, ABI cassée  |
| bridge legacy               | rtnetlink           | comportements kernel-dépendants |
| netfilter legacy (iptables) | `NETLINK_NETFILTER` | format historique tordu         |
| devlink ancien              | generic netlink     | ABI pas figée                   |
| drivers vendor              | generic netlink     | aucune garantie                 |

---

# 🧠 Règle d’or

> Si l’API existe depuis +10 ans et est utilisée par `iproute2` → safe.
> Si elle est pilotée par **un seul outil historique obscur** → danger.

---

# 🛠 Ce que tu dois faire côté code

### Pour les zones 🟢

* parser strict
* structs fixes
* pas besoin de fallback version

### Pour les zones 🟡

* parser tolérant
* ignorer attributs inconnus
* vérifier tailles dynamiques

### Pour les zones 🔴

* soit tu utilises une lib dédiée (libipvs)
* soit tu acceptes que ton code soit kernel-dépendant
* soit tu évites complètement

---

## TL;DR

Netlink est **excellent**…
mais seulement là où le kernel a décidé d’avoir une vraie API.
IPVS n’en fait tout simplement pas partie.
