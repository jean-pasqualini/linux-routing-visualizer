Excellent — là on passe du *“je sais copier des commandes”* au **“je comprends ce que fait Linux sous le capot”**.
Tu vas voir que le NAT, conntrack, nftables, Kubernetes, ta box Internet… tout est la *même mécanique*.

---

# 🌍 Big Picture — ce que tu construis vraiment

```
Client privé (10.0.0.2)
        │
        ▼
[ ROUTEUR LINUX ]
 ├─ Routage IP
 ├─ Conntrack (state)
 ├─ NAT (SNAT/MASQUERADE)
 └─ Firewall FORWARD
        │
        ▼
Serveur "public" (192.168.0.2)
```

Linux devient une **box Internet complète**.

---

# 🧱 Briques fondamentales à connaître

| Brique                   | Rôle                       |
| ------------------------ | -------------------------- |
| **Namespace réseau**     | Un mini-Linux réseau isolé |
| **veth**                 | Câble Ethernet virtuel     |
| **Routage IP**           | Décide où part un paquet   |
| **Netfilter / nftables** | Pare-feu + NAT             |
| **Conntrack**            | Mémoire du NAT (stateful)  |
| **MASQUERADE / SNAT**    | Change l’IP source         |

---

# 1️⃣ Namespace = machines virtuelles réseau

```
ip netns add ns-client
```

Tu crées **un kernel réseau isolé** :
interfaces, routes, firewall **propres à lui**.

C’est ce que fait Kubernetes avec les pods.

---

# 2️⃣ veth = câble Ethernet

```
ip link add veth-client type veth peer name veth-r1
```

C’est un câble :

```
veth-client <─────── câble ───────> veth-r1
```

Tout paquet entrant d’un côté sort **instantanément** de l’autre.

---

# 3️⃣ IP + routes = routage L3

```
10.0.0.2 → default via 10.0.0.1
```

Quand le client veut parler à `192.168.0.2` :

1. Pas dans son subnet
2. Il envoie au **gateway 10.0.0.1**
3. Le routeur choisit `veth-r2`
4. Le paquet part

Sans routes → pas d’Internet.

---

# 4️⃣ Pourquoi NAT est obligatoire

Le serveur reçoit :

```
SRC=10.0.0.2 DST=192.168.0.2
```

Mais **10.0.0.0/24 est inconnu** pour lui.

Il répondrait :

```
DST=10.0.0.2
```

et ce paquet **meurt immédiatement**.

---

# 5️⃣ NAT POSTROUTING = réécriture IP source

```
nft add rule ip nat postrouting oifname veth-r2 ip saddr 10.0.0.0/24 masquerade
```

Traduction humaine :

> Si un paquet vient de 10.0.0.0/24
> et sort par veth-r2
> alors remplace son IP source par l’IP locale de veth-r2.

Donc :

```
10.0.0.2 → devient → 192.168.0.1
```

---

# 6️⃣ Conntrack = mémoire du NAT

Quand le paquet passe, Linux enregistre :

| Original                          | Traduit                              |
| --------------------------------- | ------------------------------------ |
| 10.0.0.2:54321 → 192.168.0.2:8080 | 192.168.0.1:45012 → 192.168.0.2:8080 |

Quand le serveur répond :

```
192.168.0.2:8080 → 192.168.0.1:45012
```

Conntrack intercepte et retransforme :

```
→ 10.0.0.2:54321
```

Sans conntrack → NAT impossible.

---

# 7️⃣ FORWARD = pare-feu L3/L4

Par défaut Fedora bloque tout :

```
policy drop
```

Tu dis :

```
veth-r1 → veth-r2 : ACCEPT
ESTABLISHED,RELATED : ACCEPT
```

C’est exactement ce que fait une box :
**autoriser l’aller + laisser revenir seulement ce qui correspond à une session existante.**

---

# 🧠 Ce que tu viens d’apprendre

Tu maîtrises maintenant :

* Le même NAT que :

    * Kubernetes kube-proxy
    * Cilium masquerade
    * ta box Internet
    * firewalld
    * Docker bridge network

---

# 🧩 Règle d’or à retenir

> NAT n’est pas une fonction magique.
>
> C’est **3 briques Linux** :
>
> * Routage
> * Conntrack
> * Réécriture d’adresses

Tout le reste n’est que des wrappers.

----

Très bonne question — et elle montre que tu as compris l’essence du truc.

**On n’a pas *techniquement* besoin de 3 namespaces…
mais pédagogiquement, c’est indispensable.**

---

# Ce que représentent vraiment les 3 namespaces

| Namespace     | Rôle réel                                         |
| ------------- | ------------------------------------------------- |
| **ns-client** | Ta machine privée (PC, téléphone, pod Kubernetes) |
| **ns-router** | Ta box / nœud NAT / kube-node                     |
| **ns-server** | Internet / serveur externe                        |

Tu reconstruis **l’architecture complète d’Internet** dans ton kernel.

---

# Ce qui casserait avec seulement 2 namespaces

Si tu mets client + routeur ensemble, tu perds ça :

| Ce que tu veux observer | Pourquoi ça disparaît     |
| ----------------------- | ------------------------- |
| Vrai routage L3         | Plus de “next hop”        |
| Vraie réécriture SNAT   | Source déjà locale        |
| Table conntrack utile   | Pas de frontière réseau   |
| FORWARD chain           | Plus de trafic transitant |

Tu te retrouves à faire du NAT **sur la même machine**, ce qui n’existe pas dans la vraie vie.

---

# Le NAT a besoin de 3 mondes

Le NAT n’est pas une translation IP.

C’est un **contrat entre trois univers réseau distincts** :

```
[ réseau privé ] → [ routeur ] → [ réseau externe ]
```

S’il manque un monde →
il n’y a plus de frontière → il n’y a plus de NAT réel.

---

# Ce que tu vois grâce aux 3 namespaces

| Élément            | Tu l’observes vraiment                 |
| ------------------ | -------------------------------------- |
| POSTROUTING SNAT   | paquet traverse un autre kernel réseau |
| Conntrack stateful | entrée créée uniquement sur le routeur |
| FORWARD chain      | trafic non local                       |
| Asymétrie retour   | dé-NAT indispensable                   |

---

# Phrase clé à retenir

> Le NAT n’existe que **quand un paquet traverse une machine qui n’est ni la source ni la destination**.

C’est pour ça qu’on a besoin de **3 namespaces**.


---


🔥 **Bingo.**
Tu viens de toucher au seul vrai “bug conceptuel” du NAT : **le hairpin NAT**.

Et tu as raison :

> le NAT peut s’appliquer **même quand source et destination sont dans le même LAN**.

Mais regarde bien *ce que Linux est obligé de faire pour tricher*.

---

# 🌪 Qu’est-ce que le hairpin NAT

Un client du LAN accède au **nom public** de son propre réseau :

```
Client LAN → 86.x.x.x → NAT → Serveur LAN
```

Alors que serveur et client sont **dans le même réseau privé**.

---

# 🧠 Pourquoi c’est un hack

Normalement :

```
10.0.0.10 → 10.0.0.20
```

➡️ Pas de NAT, routage direct.

Mais si le client utilise l’IP publique :

```
10.0.0.10 → 86.x.x.x
```

La box doit :

1. DNAT → 10.0.0.20
2. **SNAT quand même le client**
3. Maintenir une entrée conntrack spéciale

Sinon le serveur répondrait directement au client → le retour contourne le NAT → connexion cassée.

---

# 🔄 Traduction réelle

| Sens          | Paquet                   |
| ------------- | ------------------------ |
| Client → Box  | 10.0.0.10 → 86.x.x.x     |
| Après DNAT    | 10.0.0.10 → 10.0.0.20    |
| Après SNAT    | **10.0.0.1 → 10.0.0.20** |
| Serveur → Box | 10.0.0.20 → 10.0.0.1     |
| Dé-NAT        | 10.0.0.20 → 10.0.0.10    |

La box se **met volontairement au milieu du LAN**.

---

# 💡 Pourquoi ça confirme la règle

Même en hairpin :

> Le NAT n’existe que s’il y a **3 rôles distincts** :
> client – routeur – serveur

La box est forcée de :

* se SNATer elle-même
* intercepter les retours
* casser le routage direct

---

# 🧠 Conclusion

Le hairpin NAT n’est pas une exception au modèle.
C’est un **contournement artificiel** pour maintenir une illusion.

Et s’il n’était pas implémenté →
la moitié des réseaux domestiques seraient inutilisables.

---

Il y a une réponse “théorique” simple…
et une réponse “réalité Linux / box Internet”.

---

# 📐 Réponse courte

👉 **Un NAT a *au minimum* 2 sous-réseaux différents.**

| Rôle                     | Sous-réseau                      |
| ------------------------ | -------------------------------- |
| Réseau privé (inside)    | ex: `10.0.0.0/24`                |
| Réseau externe (outside) | ex: `192.168.0.0/24` ou Internet |

Sans **au moins 2 sous-réseaux**, il n’y a **pas de frontière réseau**, donc pas de NAT réel.

---

# 🧠 Pourquoi 2 est le minimum

Le NAT existe uniquement pour résoudre ce problème :

> “Une adresse n’est pas routable dans l’autre monde.”

Donc il faut :

```
[ réseau A non routable depuis B ]
              ↓
            NAT
              ↓
[ réseau B qui ne connaît pas A ]
```

---

# 🌍 Réalité terrain : souvent 3 sous-réseaux

Dans une vraie box Internet tu as souvent :

| Zone                      | Exemple                     |
| ------------------------- | --------------------------- |
| LAN                       | 192.168.1.0/24              |
| DMZ / bridge / containers | 172.17.0.0/16               |
| WAN                       | IP publique / /30 opérateur |

Donc parfois **3 voire plus** :

```
LAN ─┐
DMZ ─┼─ NAT ─ WAN
IoT ─┘
```

Chaque zone est NATée indépendamment vers l’extérieur.

---

# 🌪 Cas spécial : hairpin NAT

Client et serveur sont dans le **même sous-réseau**…

Mais la box s’insère artificiellement entre eux,
ce qui recrée virtuellement **2 mondes**.

---

# 🧠 Phrase clé

> Le NAT n’est pas défini par le nombre de machines,
> mais par le nombre de **frontières réseau**.

**Minimum : 2 sous-réseaux distincts.**


----

Là on va relier **ta box, Docker, Kubernetes et le kernel Linux** — tu vas reconnaître exactement les mêmes patterns partout.

---

# 🌍 1. Vie réelle — ta box Internet

| Cas                                                       | Ce qui se passe        |
| --------------------------------------------------------- | ---------------------- |
| Tu vas sur Google                                         | **MASQUERADE / SNAT**  |
| Tu ouvres un port (NAT 443 → NAS)                         | **DNAT + SNAT retour** |
| Tu accèdes à ton NAS via ton IP publique depuis ton salon | **Hairpin NAT**        |

### Exemple réel

```
192.168.1.12 → 142.250.74.78 (google.com)
→ SNAT vers 86.43.xx.xx
```

---

# 🐳 2. Docker bridge

| Action                     | Type NAT        |
| -------------------------- | --------------- |
| Conteneur → Internet       | **MASQUERADE**  |
| `docker run -p 8080:80`    | **DNAT**        |
| Conteneur → localhost:8080 | **Hairpin NAT** |

### docker -p 8080:80

```
PREROUTING DNAT 0.0.0.0:8080 → 172.17.0.2:80
POSTROUTING SNAT/MASQ
```

---

# ☸️ 3. Kubernetes

| Fonction                    | NAT utilisé           |
| --------------------------- | --------------------- |
| Pod → Internet              | **MASQUERADE / SNAT** |
| Service ClusterIP           | **DNAT**              |
| Service NodePort            | **DNAT**              |
| Pod → Service sur même node | **Hairpin NAT**       |

### Pod → Service ClusterIP

```
PREROUTING DNAT 10.96.0.10 → 10.244.1.12
```

### Pod → Internet

```
POSTROUTING MASQUERADE
```

---

# 🔁 4. Hairpin partout

| Environnement | Exemple                                       |
| ------------- | --------------------------------------------- |
| Box maison    | Accéder au NAS via IP publique depuis le WiFi |
| Docker        | Conteneur → localhost:8080                    |
| Kubernetes    | Pod → Service → Pod sur le même node          |

Sans hairpin → ça ne marche tout simplement pas.

---

# 🧠 Règle universelle

Tous les environnements modernes utilisent **les mêmes 4 mécanismes kernel** :

| Mécanisme   | Où                       |
| ----------- | ------------------------ |
| SNAT        | Sortie vers autre réseau |
| DNAT        | Redirection d’entrée     |
| MASQUERADE  | SNAT dynamique           |
| Hairpin NAT | Boucles internes         |

Ce ne sont pas des features Docker ou Kubernetes.

👉 C’est **Linux netfilter + conntrack** partout.

