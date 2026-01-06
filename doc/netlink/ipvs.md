Voici ton **cheatsheet IPVS (Linux Virtual Server)** – le moteur de load-balancing **L4 du kernel Linux**.

---

# 🧠 IPVS – C’est quoi ?

IPVS = module kernel qui fait du **load-balancing stateful ultra-performant** (niveau TCP/UDP).

Utilisé par :

* Kubernetes kube-proxy (mode IPVS)
* HAProxy kernel-mode
* Gros load balancers L4

---

# 🔧 Activation

```bash
modprobe ip_vs
modprobe ip_vs_rr
modprobe ip_vs_wrr
modprobe ip_vs_sh
```

---

# 🧱 Concepts

| Objet           | Rôle                      |
| --------------- | ------------------------- |
| **Service**     | IP:PORT virtuel           |
| **Real Server** | backend réel              |
| **Scheduler**   | algorithme de répartition |
| **Conntrack**   | sticky session kernel     |

---

# ⚙️ Algorithmes

| Algo  | Effet                   |
| ----- | ----------------------- |
| `rr`  | round-robin             |
| `wrr` | weighted RR             |
| `lc`  | least connections       |
| `wlc` | weighted LC             |
| `sh`  | source hashing          |
| `dh`  | dest hashing            |
| `sed` | shortest expected delay |
| `nq`  | never queue             |

---

# 🏗 Modes

| Mode    | Description        |
| ------- | ------------------ |
| **NAT** | modifie src/dst    |
| **DR**  | direct routing     |
| **TUN** | encapsulation IPIP |

---

# 🛠 Commandes ipvsadm

### Créer service virtuel

```bash
ipvsadm -A -t 1.1.1.1:80 -s rr
```

### Ajouter backend

```bash
ipvsadm -a -t 1.1.1.1:80 -r 10.0.0.5:80 -m
```

`-m` = NAT mode

---

### Voir config

```bash
ipvsadm -Ln --stats --rate
```

---

### Supprimer backend

```bash
ipvsadm -d -t 1.1.1.1:80 -r 10.0.0.5:80
```

---

# 📊 Monitoring kernel

| Fichier                 | Rôle           |
| ----------------------- | -------------- |
| `/proc/net/ip_vs`       | services       |
| `/proc/net/ip_vs_conn`  | connexions     |
| `/proc/net/ip_vs_stats` | stats globales |

---

# 🧬 Exemple `/proc/net/ip_vs`

```
TCP  1.1.1.1:80 rr
  -> 10.0.0.5:80 Masq 1 123 456
```

---

# 🔥 Flush total

```bash
ipvsadm -C
```

---

# 🧠 Cas réel Kubernetes

kube-proxy IPVS crée dynamiquement :

* 1 service virtuel / Service
* 1 real server / Pod IP

---

# ⚠️ Pièges

| Problème                        | Cause           |
| ------------------------------- | --------------- |
| Backend reçoit IP client privée | NAT oublié      |
| TCP cassé                       | rp_filter actif |
| Perte sticky                    | Conntrack flush |

---

# 🎯 Résumé

IPVS = **load balancer kernel-grade**, rapide, stateful, invisible pour les apps.
