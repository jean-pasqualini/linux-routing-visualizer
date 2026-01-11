Tu es en train de regarder **la structure interne exacte d’une entrée conntrack Linux**.
Ce n’est plus du NAT « utilisateur », c’est **le cœur du moteur stateful du noyau**.

Je t’explique **chaque champ**, dans l’ordre réel d’utilisation du kernel.

---

# 🧠 Structure `Con` = une connexion réseau complète

Chaque `Con` correspond à **une session réelle suivie par netfilter** :

TCP, UDP, ICMP, FTP, SIP, etc.

---

## 1️⃣ Info *InfoSource

Métadonnées d’origine de l’entrée :

| Champ  | Rôle                                       |
| ------ | ------------------------------------------ |
| kernel | créée par netfilter                        |
| user   | injectée depuis userspace (conntrack tool) |
| event  | générée par un event kernel                |

---

## 2️⃣ Origin *IPTuple

Flux **initial** (client → serveur)

| Élément | Exemple      |
| ------- | ------------ |
| SrcIP   | 192.168.1.10 |
| DstIP   | 8.8.8.8      |
| SrcPort | 51423        |
| DstPort | 53           |
| Proto   | UDP          |

---

## 3️⃣ Reply *IPTuple

Flux **retour** (serveur → client) **après NAT**

| Élément | Exemple     |
| ------- | ----------- |
| SrcIP   | 8.8.8.8     |
| DstIP   | 82.64.23.11 |
| SrcPort | 53          |
| DstPort | 43122       |

---

## 4️⃣ ProtoInfo *ProtoInfo

Infos spécifiques au protocole.

### TCP :

| Champ           | Signification                   |
| --------------- | ------------------------------- |
| state           | SYN_SENT, ESTABLISHED, FIN_WAIT |
| flags           | SYN/ACK vus                     |
| window tracking | ajustement NAT                  |

### UDP :

| Champ     |           |
| --------- | --------- |
| last_seen | timestamp |
| timeout   | 30s       |

---

## 5️⃣ CounterOrigin / CounterReply

Statistiques par direction :

| Champ   | Rôle              |
| ------- | ----------------- |
| packets | nombre de paquets |
| bytes   | volume transféré  |

---

## 6️⃣ Helper *Helper

Module applicatif actif.

Exemples :

| Protocole | Helper |
| --------- | ------ |
| FTP       | ftp    |
| SIP       | sip    |
| H323      | h323   |
| PPTP      | pptp   |

Ils inspectent le payload et créent **connexions secondaires dynamiques**.

---

## 7️⃣ NatSrc *Nat

Infos de **NAT réellement appliqué**.

| Champ | Rôle                |
| ----- | ------------------- |
| ip    | IP source NAT       |
| port  | port source NAT     |
| range | plage utilisée      |
| flags | static / masquerade |

C’est là que MASQUERADE devient réel.

---

## 8️⃣ SeqAdjOrig / SeqAdjRepl

Correction des numéros de séquence TCP.

Quand NAT modifie la taille des paquets (FTP, SIP),
le noyau **réécrit les SEQ / ACK** pour ne rien casser.

---

## 9️⃣ ID *uint32

Identifiant interne unique de la connexion.

---

## 🔟 Status / StatusMask

Bits de statut de la connexion :

| Flag       | Signification  |
| ---------- | -------------- |
| ASSURED    | flux stable    |
| SEEN_REPLY | retour observé |
| CONFIRMED  | accepté        |
| NAT        | NAT actif      |
| DYING      | en destruction |

---

## 1️⃣1️⃣ Use *uint32

Nombre de références kernel vers cette entrée.
Empêche la suppression tant qu’elle est utilisée.

---

## 1️⃣2️⃣ Mark / MarkMask

Marquage firewall (tc, routing policy, QoS).

Permet de faire :

```bash
iptables -t mangle -A PREROUTING -j CONNMARK --set-mark 10
```

---

## 1️⃣3️⃣ Timeout *uint32

Temps restant avant expiration.

| TCP ESTABLISHED | 5 jours |
| UDP | 30s |
| ICMP | 10s |

---

## 1️⃣4️⃣ Zone *uint16

Séparation de tables conntrack.

Utilisé pour :

* conteneurs
* VRF
* namespaces réseau

---

## 1️⃣5️⃣ Timestamp *Timestamp

Heure de création / dernière activité.

---

## 1️⃣6️⃣ SecCtx *SecCtx

Contexte de sécurité (SELinux / AppArmor).

---

## 1️⃣7️⃣ Exp *Exp

Connexion **attendue** (related).

Exemple FTP :

```
PORT 192,168,1,10,195,80
```

→ création automatique d’une future connexion.

---

# 🧩 Résumé mental

Un objet `Con` =

> **Une session réseau vivante, suivie, NATée, analysée, marquée, comptée et corrigée en temps réel par le noyau.**

C’est **le cerveau du firewall stateful Linux**.
