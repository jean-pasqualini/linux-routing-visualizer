Très bonne question — et la réponse te montre **la limite de ce que fait un outil userspace (iptables)** par rapport à ce qui se passe **dans le kernel**.

---

## Quand tu écris une règle comme :

```bash
iptables -t nat -A POSTROUTING -o eth0 -j SNAT --to 10.0.0.42
```

Ce qui arrive vraiment, étape par étape :

---

## 🧠 1) iptables *ne fait rien au paquet lui-même*

`iptables` est **un outil userspace** qui :

1. Parse ta ligne de commande
2. Transforme ça en une structure interne
3. Appelle l’API **xtables / libxtables**
4. Passe des **requêtes ioctl** au kernel

👉 il ne redirige **aucun paquet** lui-même.

---

## 🧰 2) Transmission au kernel : appel à xtables

iptables utilise :

📌 **l’API x_tables** dans libc (libiptc)

Puis fait :

```text
ioctl(NFNL_SUBSYS_XTABLES, …)
```

ou si on est en mode nftables backend :

```text
netlink NFT_MSG_NEWRULE
```

Ce qui est transmis au kernel :

> “Ajoute cette règle dans la table NAT : POSTROUTING, target SNAT vers 10.0.0.42”

---

## 🧩 3) Dans le kernel : construction d’une règle Netfilter

Le kernel Netfilter reçoit la règle et :

1. Vérifie sa syntaxe
2. Crée une **structure interne** (xt_entry_match, xt_entry_target)
3. Enregistre la règle dans la chaîne NAT (`POSTROUTING`)
4. Stocke **un pointeur vers une fonction cible (target)**
   → ici la fonction SNAT

La règle ressemble conceptuellement à :

```
if (packet matches criteria) {
    call nf_nat_snat_fn(skb, state, &info);
}
```

`info` inclut l’IP de sortie (10.0.0.42), port range, etc.

---

## 🔄 4) Quand un paquet traverse POSTROUTING

Chaque paquet traverse plusieurs **hooks** Netfilter :

```
PREROUTING → FORWARD/INPUT/OUTPUT → POSTROUTING
```

Quand il arrive dans **POSTROUTING** :

* le noyau appelle un itérateur sur les règles NAT
* pour chaque règle :

    * **test de correspondance**
    * si match → **appel vers la target SNAT**

---

## 🧠 5) La target SNAT dans le kernel

La fonction interne ressemble à :

```
nf_nat_setup_info(...)
nf_nat_packet(...)
```

Concrètement, elle va :

### 🧱 a) Regarder le paquet

* protocole (TCP/UDP/ICMP)
* IP source
* port source

### 🔁 b) Faire le NAT

Pour SNAT :

* modifier **l’adresse source**
* modifier **le port source** (si configuré)
* mettre à jour les **checksums** (IP + L4)

### 📌 c) Conntrack

Le kernel utilise **conntrack** :

* crée (ou met à jour) une **entrée de suivi de flux**
* stocke :

    * tuple original : `(src, dst, sport, dport, proto)`
    * tuple traduit : `(newsrc, dst, new_sport, dport, proto)`
* la traduction inverse (DNAT du retour) se base sur cette entrée

Sans conntrack, SNAT ne fonctionnerait qu’unidirectionnellement.

---

## 🧪 6) Résultat dans la table d’exécution

À la fin, quand le paquet quitte l’interface :

* **IP source est réécrite**
* **Port source est réécrit si demandé**
* Il est routé vers la destination finale

---

## 🧠 En résumé

| Étape                            | Où ça se passe       |
| -------------------------------- | -------------------- |
| Tu écris la règle `-j SNAT`      | userspace (iptables) |
| La règle est traduite            | xtables / netlink    |
| Le kernel stocke la règle        | netfilter NAT table  |
| Le paquet MATCHE → SNAT appliqué | kernel (Netfilter)   |
| Conntrack mémorise la traduction | kernel conntrack     |

---

## 🧩 Visualisation simple

```
Userspace:
iptables → parse → netlink/xtables → kernel

Kernel:
PACKET →
  Netfilter PRE/POSTROUTING →
    test règles →
    target SNAT (réécriture + conntrack)
```

---

## Point clé

👉 **iptables ne fait pas lui-même le NAT.**
Il **configure** le kernel pour que ce soit le kernel qui le fasse au moment de l’émission des paquets.

---

Si tu veux, je peux te montrer **le pseudo-code approximatif de la target SNAT dans le kernel**, ou te pointer vers les fichiers sources exacts dans le code Linux.

----

```
struct xt_table  (table "nat")
 └── struct xt_table_info
      └── tableau de struct ipt_entry
            ├── matches
            └── target (ex: SNAT)

net/ipv4/netfilter/iptable_nat.c
net/netfilter/x_tables.c
include/linux/netfilter/x_tables.h

struct nft_table
 └── struct nft_chain
      └── struct nft_rule
           └── tableau de struct nft_expr
                 ├── payload
                 ├── cmp
                 └── nat
net/netfilter/nf_tables*.c
include/net/netfilter/nf_tables.h

```

```aiignore
struct nf_conn
 └── tuple original
 └── tuple NAT
 └── timeouts
 └── state TCP
```


-----

Très bonne question — tu veux savoir **qui fait quoi exactement dans le kernel**. Voilà la cartographie précise.

---

# 1️⃣ Qui parcourt les règles ?

Ça se passe dans **Netfilter hooks**.

Fonction clé (IPv4/IPv6 unifié) :

```
nf_nat_inet_fn()
```

📍 Fichier :

```
net/netfilter/nf_nat_core.c
```

Cette fonction est enregistrée comme **hook** sur :

* `NF_INET_PRE_ROUTING`
* `NF_INET_POST_ROUTING`
* `NF_INET_LOCAL_OUT`

Elle est appelée pour **chaque paquet** qui traverse ces hooks.

---

# 2️⃣ Qui exécute les règles nftables / xtables ?

### Backend nftables :

```
nft_do_chain()
```

📍

```
net/netfilter/nf_tables_core.c
```

Elle itère sur :

```
struct nft_chain → struct nft_rule → struct nft_expr
```

Quand elle rencontre une expression de type `nft_nat` :

```
nft_nat_eval()
```

📍

```
net/netfilter/nft_nat.c
```

---

# 3️⃣ Qui appelle nf_nat_setup_info() ?

C’est précisément :

```
nft_nat_eval()
```

Elle appelle :

```
nf_nat_setup_info(ct, range, maniptype)
```

📍

```
net/netfilter/nf_nat_core.c
```

---

# 4️⃣ Qui crée / met à jour le conntrack ?

Dans :

```
nf_nat_setup_info()
```

Elle appelle :

```
nf_nat_alloc_null_binding()
nf_nat_packet()
```

Ces fonctions :

* allouent / récupèrent un `struct nf_conn`
* remplissent :

    * tuple original
    * tuple traduit
* installent la **relation NAT bidirectionnelle**

---

# 5️⃣ Où est stockée la traduction active ?

Dans la structure :

```
struct nf_conn
```

📍

```
include/net/netfilter/nf_conntrack.h
```

Elle contient :

```
tuple original
tuple reply
nat info
timeouts
state TCP
```

---

# 6️⃣ Qui applique la réécriture sur chaque paquet ?

Quand le paquet suivant de la même connexion arrive :

```
nf_nat_packet(ct, skb, maniptype)
```

📍

```
net/netfilter/nf_nat_core.c
```

C’est elle qui :

* modifie IP / port
* appelle les helpers TCP/UDP
* recalcule checksums

---

# 7️⃣ Vue complète pipeline

```
skb arrive
 ↓
nf_hook_slow()
 ↓
nf_nat_inet_fn()
 ↓
nft_do_chain()
 ↓
nft_nat_eval()
 ↓
nf_nat_setup_info()   ← première fois
 ↓
nf_nat_packet()       ← chaque paquet
 ↓
skb sort modifié
```

---

## Résumé

| Action               | Fonction kernel       |
| -------------------- | --------------------- |
| Hook NAT             | `nf_nat_inet_fn()`    |
| Parcours règles      | `nft_do_chain()`      |
| Exécution règle NAT  | `nft_nat_eval()`      |
| Création NAT         | `nf_nat_setup_info()` |
| Réécriture paquet    | `nf_nat_packet()`     |
| Stockage translation | `struct nf_conn`      |

C’est ce pipeline que tu devrais lire si tu veux vraiment comprendre le NAT Linux de l’intérieur.

----

```
                Kernel networking stack
                         │
                         ▼
                Netfilter core
        (hooks NF_INET_*)
                         │
                         ▼
              nftables / xtables engine
                         │
                         ▼
                  nf_nat (NAT engine)
                         │
                         ▼
                   conntrack (nf_conntrack)
                         │
                         ▼
                 rewriting IP / TCP / UDP

```