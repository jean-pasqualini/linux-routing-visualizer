Le **hairpin NAT** (aussi appelé *NAT loopback*) c’est un cas très particulier où :

> un client **dans le même réseau** accède à un service **via son IP publique / VIP**,
> et non via son IP privée réelle.

---

## 🎯 Exemple simple

Tu as :

* Client : `192.168.1.10`
* Serveur réel : `192.168.1.20`
* VIP publique : `203.0.113.50` → DNAT vers `192.168.1.20`
* Firewall/NAT : `203.0.113.1`

Le client interne fait :

```
192.168.1.10 → 203.0.113.50:443
```

Alors que le serveur est **dans le même LAN**.

---

## ❌ Ce qui se passe sans hairpin

Le firewall fait le DNAT :

```
192.168.1.10 → 192.168.1.20
```

Le serveur répond :

```
192.168.1.20 → 192.168.1.10
```

Mais le client attend une réponse venant de `203.0.113.50`,
pas de `192.168.1.20`.

➡️ Conntrack ne reconnaît pas le flux → le paquet retour est rejeté.

---

## ✅ Ce que fait le hairpin NAT

Le firewall fait **deux translations** :

1. **DNAT**

   ```
   203.0.113.50 → 192.168.1.20
   ```
2. **SNAT forcé** (vers sa propre IP)

   ```
   192.168.1.10 → 203.0.113.1
   ```

Résultat vu par le serveur :

```
203.0.113.1 → 192.168.1.20
```

Réponse du serveur :

```
192.168.1.20 → 203.0.113.1
```

Le firewall renvoie ensuite vers le client original.

---

## 🧠 Signature conntrack

Dans ce cas tu vois toujours :

```
SRC_NAT_DONE | DST_NAT_DONE
```

même si source et destination sont dans le même subnet.

---

## Résumé ultra-court

Hairpin NAT =
**accéder à un service local via son IP publique**,
et le firewall doit faire **DNAT + SNAT** pour que ça fonctionne.
