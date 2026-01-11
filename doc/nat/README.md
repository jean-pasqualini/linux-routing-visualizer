Voici une explication complète et claire du **NAT sur LINX (Linux)**, de ce qui se passe réellement dans le noyau jusqu’aux règles `iptables` / `nftables`.

---

## 1️⃣ Qu’est-ce que le NAT ?

Le **NAT – Network Address Translation** permet de **modifier les adresses IP et / ou ports dans les paquets** quand ils traversent une machine Linux.

Cas typique :
Un réseau privé `192.168.1.0/24` sort sur Internet via une seule IP publique.

| Machine privée | IP vue sur Internet |
| -------------- | ------------------- |
| 192.168.1.10   | 82.64.23.11         |
| 192.168.1.20   | 82.64.23.11         |

Le routeur Linux modifie les paquets à la volée.

---

## 2️⃣ Les 3 types de NAT

| Type           | Rôle                                              |
| -------------- | ------------------------------------------------- |
| **SNAT**       | Change l’adresse source (sortie vers Internet)    |
| **DNAT**       | Change l’adresse de destination (port-forwarding) |
| **MASQUERADE** | Variante dynamique de SNAT (IP publique variable) |

---

## 3️⃣ Où le NAT agit dans le noyau ?

Le NAT est traité par **netfilter** dans des *hooks* précis.

```
        PREROUTING --> DNAT
              |
           ROUTING
              |
INPUT       FORWARD        OUTPUT
              |
        POSTROUTING --> SNAT / MASQUERADE
```

| Chaîne      | Moment                     |
| ----------- | -------------------------- |
| PREROUTING  | Avant le routage → DNAT    |
| POSTROUTING | Après le routage → SNAT    |
| OUTPUT      | Paquets générés localement |

---

## 4️⃣ Exemple concret

Réseau interne : `192.168.1.0/24`
Interface LAN : `eth1`
Interface WAN : `eth0`

---

### 🔹 Activer le routage IP

```bash
echo 1 > /proc/sys/net/ipv4/ip_forward
```

Permanent :

```bash
net.ipv4.ip_forward=1
```

---

## 5️⃣ MASQUERADE (le plus courant)

```bash
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
```

Effet :

| Avant NAT                       | Après NAT                      |
| ------------------------------- | ------------------------------ |
| 192.168.1.10:51423 → 8.8.8.8:53 | 82.64.23.11:43122 → 8.8.8.8:53 |

Linux garde une table de correspondance :

```
192.168.1.10:51423 ↔ 82.64.23.11:43122
```

---

## 6️⃣ DNAT – Port forwarding

Rediriger le port 80 public vers un serveur interne :

```bash
iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 \
-j DNAT --to-destination 192.168.1.10:80
```

Autoriser le forward :

```bash
iptables -A FORWARD -p tcp -d 192.168.1.10 --dport 80 -j ACCEPT
```

---

## 7️⃣ Table de conntrack

Linux suit chaque connexion avec **conntrack** :

```bash
conntrack -L
```

Exemple :

```
tcp  6 431999 ESTABLISHED
src=192.168.1.10 dst=1.1.1.1 sport=51233 dport=443
src=1.1.1.1 dst=82.64.23.11 sport=443 dport=43122
```

---

## 8️⃣ nftables (nouvelle méthode)

Même chose avec nftables :

```bash
nft add table ip nat
nft add chain ip nat postrouting { type nat hook postrouting priority 100 \; }
nft add rule ip nat postrouting oif "eth0" masquerade
```

DNAT :

```bash
nft add rule ip nat prerouting iif "eth0" tcp dport 80 dnat to 192.168.1.10:80
```

---

## 9️⃣ Différence SNAT vs MASQUERADE

| SNAT        | MASQUERADE            |
| ----------- | --------------------- |
| IP fixe     | IP dynamique          |
| Plus rapide | Un peu plus lent      |
| Serveurs    | Connexions PPP / DHCP |

---

## 🔟 Résumé

Le NAT Linux :

✔ modifie IP/ports en temps réel
✔ repose sur conntrack
✔ fonctionne via `PREROUTING` et `POSTROUTING`
✔ permet partage d’IP publique, redirections, DMZ, etc.

---

Si tu veux, dis-moi ton schéma réseau et je te fais une config NAT complète prête à l’emploi.


---

C’est **LE point le plus important du NAT**, et beaucoup de gens se trompent dessus. Je t’explique simplement.

---

## 1️⃣ IP fixe vs IP dynamique – c’est quoi ?

### 🔹 IP fixe

Ton routeur a **toujours la même IP publique**.

| Aujourd’hui | Demain      | Dans 6 mois |
| ----------- | ----------- | ----------- |
| 82.64.23.11 | 82.64.23.11 | 82.64.23.11 |

👉 Cas typique : fibre pro, IP réservée.

---

### 🔹 IP dynamique

Ton FAI te donne une IP qui **change régulièrement** (reboot box, reconnexion, timeout DHCP…).

| Aujourd’hui | Demain     |
| ----------- | ---------- |
| 82.64.23.11 | 90.12.44.8 |

👉 Cas typique : box ADSL / fibre grand public.

---

## 2️⃣ Pourquoi ça change tout pour le NAT ?

Quand Linux fait du NAT sortant, il doit savoir :

> **Quelle IP source utiliser pour remplacer 192.168.1.10 ?**

---

## 3️⃣ Avec IP fixe → SNAT

Tu écris la règle **une fois** :

```bash
iptables -t nat -A POSTROUTING -o eth0 -j SNAT --to-source 82.64.23.11
```

Linux remplace toujours par `82.64.23.11`.

✔ rapide
✔ stable
❌ si l’IP change → plus rien ne sort

---

## 4️⃣ Avec IP dynamique → MASQUERADE

```bash
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
```

Linux fait :

1. « Quelle est mon IP actuelle sur eth0 ? »
2. Il l’utilise **à la volée** pour chaque paquet.

Donc quand l’IP change :

| Avant         | Après reconnexion    |
| ------------- | -------------------- |
| 82.64.23.11   | 90.12.44.8           |
| Le NAT marche | Le NAT marche encore |

✔ aucun script à modifier
✔ parfait pour box internet
❌ un peu plus lent

---

## 5️⃣ Résumé ultra simple

| Ton IP publique  | Tu utilises  |
| ---------------- | ------------ |
| Toujours la même | `SNAT`       |
| Change souvent   | `MASQUERADE` |

---

## 6️⃣ Phrase magique à retenir

> **SNAT = je connais mon IP à l’avance.**
> **MASQUERADE = je demande à Linux quelle est mon IP à chaque paquet.**

---

Si tu me dis quel type de connexion tu as (box, fibre pro, IP fixe ou non), je te donne **la règle exacte adaptée à ton cas**.


---

Oui. **MASQUERADE est un SNAT automatique.**

Il n’existe **que** en sortie, dans `POSTROUTING`.
Il ne peut **jamais** faire du DNAT.

---

## Ce que fait vraiment MASQUERADE

Quand tu écris :

```bash
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
```

Linux fait exactement comme :

```bash
SNAT --to-source <IP_actuelle_de_eth0>
```

mais **sans que tu aies besoin de connaître l’IP à l’avance**.

---

## Comparaison directe

| SNAT                  | MASQUERADE                      |
| --------------------- | ------------------------------- |
| Tu fournis l’IP       | Linux la découvre               |
| IP fixe requise       | IP dynamique OK                 |
| Plus rapide           | Légèrement plus lent            |
| Connexions longues OK | Connexions cassées si IP change |
| Serveurs              | Box / PPP / DHCP                |

---

## Où il agit dans le noyau

```
POSTROUTING  --->  SNAT / MASQUERADE
```

Jamais ailleurs.

---

## Règle d’or

> 🔹 **DNAT = PREROUTING**
> 🔹 **SNAT / MASQUERADE = POSTROUTING**
> 🔹 **MASQUERADE = SNAT dynamique**

---

Donc ta phrase est parfaite :

> **Le MASQUERADE se fait uniquement sur du SNAT.**
