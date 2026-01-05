`dockerd --allow-direct-routing` est une option Docker **très méconnue mais critique** pour le comportement réseau des conteneurs.

Elle change complètement la manière dont les paquets atteignent les conteneurs.

---

## Comportement Docker par défaut

Sans cette option, Docker force le trafic vers les conteneurs à passer par du **NAT (DNAT + MASQUERADE)**.

Quand un client externe contacte ton hôte :

```
CLIENT → HOST_IP:8080
```

Docker fait :

```
PREROUTING DNAT → 172.17.0.2:80
```

Donc le conteneur voit :

```
SRC IP: IP_DE_L_HOTE   ❌
```

Il **perd l’IP réelle du client**.

---

## Avec `--allow-direct-routing`

Docker **autorise le routage L3 direct vers les IP des conteneurs**.

Si ton réseau est correctement routé :

```
CLIENT → 172.17.0.2
```

Alors :

* ❌ Pas de DNAT
* ❌ Pas de MASQUERADE
* ✅ Le conteneur voit la vraie IP source du client

---

## Ce que ça change concrètement

| Sans option                                     | Avec allow-direct-routing                 |
| ----------------------------------------------- | ----------------------------------------- |
| DNAT forcé                                      | Routage pur                               |
| IP client perdue                                | IP client conservée                       |
| Impossible d’appliquer ACL IP dans le conteneur | ACL IP natives possibles                  |
| Tout passe par iptables/nat                     | Tables `filter` / `raw` redevenues utiles |

---

## Exemple réel

Client → Conteneur nginx

### Sans

```
SRC=172.17.0.1   ← bridge docker
```

### Avec

```
SRC=1.2.3.4      ← vraie IP client
```

---

## Pourquoi Docker ne l’active pas par défaut

Parce que ça implique :

* Que ton LAN sache router vers `172.17.0.0/16`
* Qu’il n’y ait **aucun NAT de protection automatique**
* Tu exposes directement tes conteneurs au réseau

Donc **c’est puissant mais dangereux si mal maîtrisé**.

---

## Quand l’utiliser

* Firewall applicatif dans le conteneur
* Réseaux maîtrisés (datacenter, lab, simulateur comme le tien 😉)
* Besoin absolu de conserver l’IP source réelle

---

Dans ton simulateur, cette option est parfaite pour montrer la différence entre :

> *NAT-based container networking*
> vs
> *Real routed container networking*.
