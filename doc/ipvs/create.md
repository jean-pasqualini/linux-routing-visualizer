sudo lsmod | grep ip_vs
Tu veux donc **créer un service IPVS** (un virtual service) pour commencer à recevoir du trafic / de la donnée.
Voilà exactement comment faire proprement.

---

# 🧱 1️⃣ Charger les modules IPVS

Sur le host (pas dans le conteneur) :

```bash
modprobe ip_vs
modprobe ip_vs_rr
modprobe ip_vs_wlc
```

Vérifie :

```bash
lsmod | grep ip_vs
```

---

# 🧱 2️⃣ Créer un service virtuel

On crée un VIP `1.1.1.1:8080` en round-robin :

```bash
ipvsadm -A -t 1.1.1.1:8080 -s rr
```

---

# 🧱 3️⃣ Ajouter des backends

```bash
ipvsadm -a -t 1.1.1.1:8080 -r 10.0.0.5:8080 -m
ipvsadm -a -t 1.1.1.1:8080 -r 10.0.0.6:8080 -m
```

`-m` = NAT mode (le plus simple).

---

# 🧱 4️⃣ Activer le routage

```bash
sysctl -w net.ipv4.ip_forward=1
```

---

# 🧱 5️⃣ Autoriser NAT / forwarding

```bash
iptables -t nat -A POSTROUTING -j MASQUERADE
iptables -A FORWARD -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A FORWARD -d 10.0.0.0/24 -p tcp --dport 8080 -j ACCEPT
```

---

# 🧱 6️⃣ Tester

```bash
curl http://1.1.1.1:8080
```

Puis :

```bash
watch -n1 ipvsadm -Ln --stats --rate
```

Tu verras les compteurs monter.

---

# 🧠 Résultat

Chaque connexion vers `1.1.1.1:8080` est :

* interceptée par IPVS
* répartie sur `10.0.0.5` / `10.0.0.6`
* suivie par conntrack
* visible dans `/proc/net/ip_vs_conn`

---

# ⚠️ Important

Si tu fais ça dans Docker / Kubernetes → il faut :

* `--net=host`
* accès `/proc/net/ip_vs`

---

Tu viens de créer un **load-balancer kernel natif**.
