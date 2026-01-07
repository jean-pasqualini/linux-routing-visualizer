Très bonne question, là tu es au cœur de la logique réseau Linux.

Quand **une interface a plusieurs IP**, Linux choisit l’IP source selon un **algorithme précis**.

---

## Ordre réel de sélection de l’IP source

Quand tu envoies un paquet sans préciser la source :

1. **`ip rule` (policy routing)**
   Si une règle impose une table selon la source → utilisée.

2. **Champ `src` de la route**
   S’il existe :

   ```
   ... src 172.25.2.1
   ```

   → cette IP est choisie immédiatement.

3. **Sinon : choix automatique via RFC 6724**

Linux applique des règles de préférence :

| Priorité | Critère                                               |
| -------- | ----------------------------------------------------- |
| 1        | IP dans le **même subnet** que la destination         |
| 2        | IP ayant la **plus longue correspondance de préfixe** |
| 3        | IP avec le **scope le plus petit** (link > global)    |
| 4        | IP la plus “stable” (non temporaire, non deprecated)  |

---

## Exemple réel

Interface `eth0` :

```
192.168.1.10/24
192.168.1.11/24
10.0.0.5/24
```

Destination : `192.168.1.20`

Linux choisit :

```
SRC = 192.168.1.10 ou .11   (même subnet /24)
```

mais jamais `10.0.0.5`.

---

## Cas Docker typique

Bridge :

```
br-xxxx:
  172.25.2.1/24
  172.25.3.1/24
```

Destination : `172.25.3.5`

Linux choisit :

```
SRC = 172.25.3.1
```

car correspondance de subnet la plus précise.

---

## Comment forcer le choix

### 1. Dans la route

```bash
ip route replace 172.25.2.0/24 dev br-xxxx src 172.25.2.1
```

### 2. Par policy routing

```bash
ip rule add from 172.25.2.1 table docker
```

---

## Résumé simple

Quand une interface a plusieurs IP :

Linux choisit l’IP source qui est **la plus logique par rapport à la destination**
(same subnet, meilleur préfixe, meilleure portée).
