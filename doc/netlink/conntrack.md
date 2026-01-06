`conntrack` est **le cerveau mémoire du firewall Linux**.
Sans lui : pas de NAT fiable, pas de stateful firewall, pas d’anti-DDoS propre.

---

# 🧠 C’est quoi Conntrack ?

Conntrack = **suivi d’état des connexions réseau** au niveau kernel.

Chaque flux réseau devient un objet :

```
(src IP, src port, dst IP, dst port, proto)  
→ état TCP, timeout, NAT mapping, mark, helper
```

---

# 🔄 États gérés

| Protocole | États                                    |
| --------- | ---------------------------------------- |
| TCP       | NEW → ESTABLISHED → FIN_WAIT → TIME_WAIT |
| UDP       | UNREPLIED → ASSURED                      |
| ICMP      | REQUEST → REPLY                          |
| Autres    | GENERIC                                  |

---

# 🌍 Sans Conntrack

| Fonction          | Résultat     |
| ----------------- | ------------ |
| NAT               | ❌ impossible |
| Firewall stateful | ❌            |
| Load-balancing    | ❌            |
| FTP / SIP helpers | ❌            |
| Anti-DDoS précis  | ❌            |

---

# 🛠️ Ce que tu peux faire avec Conntrack

## 1️⃣ Voir toutes les connexions actives

```bash
conntrack -L
```

---

## 2️⃣ Supprimer des connexions en live

```bash
conntrack -D -s 10.0.0.5
conntrack -D -d 1.2.3.4
```

🔥 Coupe toutes les connexions sans toucher aux règles firewall.

---

## 3️⃣ Monitorer en temps réel

```bash
conntrack -E
```

Tu vois :

* nouvelles connexions
* timeouts
* fermetures

Parfait pour IDS maison.

---

## 4️⃣ Modifier le comportement NAT

| Action          | Effet              |
| --------------- | ------------------ |
| Changer timeout | accélère nettoyage |
| Flush partiel   | stoppe un DDoS     |
| Modifier mark   | routage avancé     |

---

## 5️⃣ Load-balancing stateful

Conntrack garantit que :

```
client → serveur A
tous les paquets suivants → serveur A
```

Sans ça, round-robin casserait TCP.

---

## 6️⃣ Détection DDoS intelligente

Tu peux détecter :

* trop de SYN sans ACK
* explosion d’UNREPLIED UDP
* flood ICMP

Et tuer **seulement** les flux toxiques.

---

# 🔬 Exemple entrée Conntrack

```
tcp  6 431999 ESTABLISHED src=10.0.0.5 dst=8.8.8.8 sport=55322 dport=53
     src=8.8.8.8 dst=1.1.1.1 sport=53 dport=55322 [ASSURED] mark=0 use=1
```

Ici :

* mapping NAT
* état TCP
* durée de vie
* compteur d’usage

---

# ⚠️ Limites critiques

| Problème            | Impact                         |
| ------------------- | ------------------------------ |
| Table pleine        | plus aucune nouvelle connexion |
| Timeouts trop longs | DoS facile                     |
| Pas de flush        | fuite mémoire kernel           |

---

# 🧬 En Netlink

Avec `NETLINK_NETFILTER` tu peux :

* recevoir `NEW / DESTROY` events
* supprimer des entrées
* créer un firewall *réactif* au trafic réel

---

Conntrack n’est pas un outil.
C’est la **mémoire vivante d’Internet sur ton serveur**.
