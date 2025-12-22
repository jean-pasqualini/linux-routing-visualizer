# Netfilter / iptables — Tables, Chaînes et Usage

Ce document liste **toutes les tables iptables**, leurs **chaînes associées** et leur **usage**, dans un format **copiable tel quel**.

---

## TABLE: raw

**But**
- Traitement le plus tôt possible
- Avant le suivi de connexion (conntrack)
- Possibilité de désactiver le suivi (`NOTRACK`)

**Chaînes**
- `PREROUTING` : paquets entrants, avant toute décision
- `OUTPUT`     : paquets générés localement

**Usage typique**
- `TRACE`
- `NOTRACK`

---

## TABLE: mangle

**But**
- Modifier les paquets
- QoS, marquage, altération d’en-têtes

**Chaînes**
- `PREROUTING`
- `INPUT`
- `FORWARD`
- `OUTPUT`
- `POSTROUTING`

**Usage typique**
- `MARK`
- `CONNMARK`
- `DSCP`
- `TTL`

---

## TABLE: nat

**But**
- Traduction d’adresses (NAT)
- Appliquée uniquement au **premier paquet** d’une connexion

**Chaînes**
- `PREROUTING`  : DNAT avant décision de routage
- `OUTPUT`      : DNAT pour trafic local
- `POSTROUTING` : SNAT / `MASQUERADE` après routage

**Usage typique**
- `DNAT`
- `SNAT`
- `MASQUERADE`
- `REDIRECT`

---

## TABLE: filter

**But**
- Décision finale sur les paquets
- Autoriser ou bloquer le trafic

**Chaînes**
- `INPUT`   : trafic à destination de la machine locale
- `FORWARD`: trafic en transit
- `OUTPUT` : trafic généré localement

**Usage typique**
- `ACCEPT`
- `DROP`
- `REJECT`
- `LOG`

---

## TABLE: security

**But**
- Intégration avec les LSM (ex: SELinux)
- Contrôles de sécurité avancés

**Chaînes**
- `INPUT`
- `FORWARD`
- `OUTPUT`

**Usage typique**
- `SECMARK`
- `AUDIT`

---

## Ordre de passage simplifié (pipeline)

raw PREROUTING
mangle PREROUTING
nat PREROUTING
↓
décision de routage
↓
mangle INPUT      → filter INPUT        (trafic local)
mangle FORWARD    → filter FORWARD      (transit)
mangle OUTPUT     → nat OUTPUT          (trafic local sortant)
↓
nat POSTROUTING
mangle POSTROUTING

---

## Rappels importants

- Une règle **ne peut jamais sauter entre tables**
- `ACCEPT`, `DROP`, `REJECT`, `LOG`, `TRACE` sont des **targets**, pas des chaînes
- `LOG` est **non terminal**
- `TRACE` est **debug uniquement**
- En **nftables**, les tables n’existent que si elles sont créées

---

## TL;DR

raw       → très tôt, avant conntrack
mangle    → modifier les paquets
nat       → traduire les adresses
filter    → décider (autoriser / bloquer)
security  → LSM / SELinux



Si tu veux, je peux aussi te fournir :
•	la version nftables native équivalente
•	un schéma ASCII détaillé
•	une fiche mémo ultra-courte à imprimer