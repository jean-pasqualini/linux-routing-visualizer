Très bonne question — comprendre ça te fait passer du mode *“je tape des règles”* au mode *“je conçois un firewall”* 😄

Voici le rôle **réel** de chaque table iptables.

---

## 1️⃣ `filter` — le vrai firewall

👉 **Autoriser / bloquer le trafic**

Chaînes : `INPUT`, `OUTPUT`, `FORWARD`

Exemples :

```bash
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -s 1.2.3.4 -j DROP
```

C’est **la seule table faite pour filtrer**.

---

## 2️⃣ `nat` — traduction d’adresses

👉 Modifier les IP / ports

Chaînes : `PREROUTING`, `POSTROUTING`, `OUTPUT`

Exemples :

```bash
# redirection de port
iptables -t nat -A PREROUTING -p tcp --dport 80 -j DNAT --to 192.168.1.10

# masquerading
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
```

❌ Jamais de DROP ici.

---

## 3️⃣ `mangle` — marquage / QoS

👉 Ajouter des marques aux paquets pour :

* routage avancé
* qualité de service (tc)
* priorisation

Exemples :

```bash
iptables -t mangle -A PREROUTING -p tcp --dport 22 -j MARK --set-mark 10
```

---

## 4️⃣ `raw` — avant conntrack

👉 Contrôler le suivi d’état

Chaînes : `PREROUTING`, `OUTPUT`

Exemple :

```bash
iptables -t raw -A PREROUTING -p udp --dport 123 -j NOTRACK
```

Uniquement pour cas très avancés.

---

## 5️⃣ `security` — intégration SELinux / LSM

👉 Contexte de sécurité

Rarement utilisé manuellement.

---

## Résumé clair

| Table      | À quoi elle sert           |
| ---------- | -------------------------- |
| `filter`   | Autoriser / bloquer        |
| `nat`      | DNAT / SNAT / redirections |
| `mangle`   | Marquer / QoS              |
| `raw`      | Gérer conntrack            |
| `security` | SELinux                    |

---

### Phrase à retenir

> Si tu bloques un paquet ailleurs que dans `filter`,
> tu n’es plus en train de faire un firewall — tu fais de la chirurgie réseau.
