### From the network
NIC → raw:PREROUTING → mangle → nat → routing → filter:INPUT → socket
### Localhost
process → raw:OUTPUT → mangle → nat → filter:OUTPUT → socket


# WTF

Très bonne question — et elle est **subtile**.
La clé est ici 👇

---

## 🧠 La règle fondamentale (celle qui débloque tout)

> **Netfilter raisonne sur l’ORIGINE du paquet, pas sur l’interface physique.**

Même si le paquet passe par `lo`, **il n’est pas considéré comme “reçu”**.

---

## 🔁 Ce qui se passe vraiment avec `localhost`

Quand tu fais :

```bash
curl localhost:9090
```

### Étapes réelles (simplifiées) :

1. Le processus `curl` crée un paquet TCP
2. Le noyau marque ce paquet comme **`LOCAL_OUT`**
3. Le routage décide : *destination = local*
4. Le paquet est injecté **directement dans la stack locale**
5. Il est livré au socket serveur

👉 **Il n’entre jamais par le chemin “réception” (`RX`)**
👉 Il passe par `lo`, mais **comme raccourci interne**, pas comme paquet entrant

---

## 🧭 Différence RX vs LOOPBACK

| Concept               | INPUT |
| --------------------- | ----- |
| RX (réception réseau) | ✅     |
| Loopback interne      | ❌     |

`lo` n’est pas une “carte qui reçoit” ici,
c’est une **optimisation interne du noyau**.

---

## 🔍 Schéma mental correct

```
Process A
   ↓ (LOCAL_OUT)
OUTPUT
   ↓
[ shortcut interne ]
   ↓
Process B
```

❌ Pas de :

* PREROUTING
* INPUT

---

## 🧪 Preuve simple par LOG

```bash
iptables -I INPUT -i lo -j LOG --log-prefix "INPUT_LO: "
iptables -I OUTPUT -o lo -j LOG --log-prefix "OUTPUT_LO: "
```

Puis :

```bash
curl localhost:9090
```

Résultat :

* ✅ `OUTPUT_LO`
* ❌ jamais `INPUT_LO`

---

## 🧠 Pourquoi le noyau fait ça ?

Pour :

* éviter des copies inutiles
* éviter une double traversée de la stack
* améliorer les performances

C’est un **raccourci volontaire**.

---

## 🧠 Résumé en une phrase

> **`lo` n’implique pas INPUT — seule une vraie réception réseau le fait.**

---

Si tu veux, je peux te montrer :

* le **flag kernel exact** (`skb->pkt_type`)
* un **schéma complet RX vs LOCAL_OUT**
* comment **forcer un passage INPUT avec un namespace réseau**

Dis-moi 👍
