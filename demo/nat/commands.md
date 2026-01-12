Parfait. Là tu vas avoir **la version Fedora “from scratch” complète**, 100 % nftables, 0 magie, 0 firewalld.

À la fin tu comprendras **exactement comment Fedora fait du NAT** au niveau kernel.

---

# 🧠 Architecture

```
ns-client ──10.0.0.2──▶ veth-client ─▶ veth-r1 ─▶ ns-router ─▶ veth-r2 ─▶ veth-server ──192.168.0.2──▶ ns-server
                         10.0.0.1                       192.168.0.1
                              ←────────── SNAT / MASQUERADE ──────────
```

---

# 0️⃣ Paquets Fedora

```bash
sudo dnf install -y iproute nftables conntrack-tools tcpdump
```

---

# 1️⃣ Activer le routage kernel

```bash
sudo sysctl -w net.ipv4.ip_forward=1
```

---

# 2️⃣ Créer les namespaces

```bash
sudo ip netns add ns-client
sudo ip netns add ns-router
sudo ip netns add ns-server
```

---

# 3️⃣ Créer les liens veth

```bash
sudo ip link add veth-client type veth peer name veth-r1
sudo ip link add veth-r2 type veth peer name veth-server

sudo ip link set veth-client netns ns-client
sudo ip link set veth-r1     netns ns-router
sudo ip link set veth-r2     netns ns-router
sudo ip link set veth-server netns ns-server
```

---

# 4️⃣ IPs

```bash
sudo ip netns exec ns-client ip addr add 10.0.0.2/24 dev veth-client
sudo ip netns exec ns-router ip addr add 10.0.0.1/24 dev veth-r1

sudo ip netns exec ns-router ip addr add 192.168.0.1/24 dev veth-r2
sudo ip netns exec ns-server ip addr add 192.168.0.2/24 dev veth-server
```

---

# 5️⃣ Monter les interfaces

```bash
for ns in ns-client ns-router ns-server; do
  sudo ip netns exec $ns ip link set lo up
done

sudo ip netns exec ns-client ip link set veth-client up
sudo ip netns exec ns-router ip link set veth-r1 up
sudo ip netns exec ns-router ip link set veth-r2 up
sudo ip netns exec ns-server ip link set veth-server up
```

---

# 6️⃣ Routes

```bash
sudo ip netns exec ns-client ip route add default via 10.0.0.1
sudo ip netns exec ns-server ip route add default via 192.168.0.1
```

---

# 7️⃣ nftables NAT réel Fedora

Entre dans le routeur :

```bash
sudo ip netns exec ns-router bash
```

### Créer table NAT

```bash
nft add table ip nat
nft add chain ip nat postrouting { type nat hook postrouting priority 100 \; }
```

### Ajouter la règle SNAT dynamique

```bash
nft add rule ip nat postrouting oifname "veth-r2" ip saddr 10.0.0.0/24 masquerade
```

➡️ Ça veut dire :
**“Tout paquet venant de 10.0.0.0/24 et sortant par veth-r2 → change son IP source vers l’IP locale de veth-r2.”**

---

# 8️⃣ Table FORWARD (obligatoire sinon tout est bloqué)

```bash
nft add table ip filter
nft add chain ip filter forward { type filter hook forward priority 0 \; policy drop \; }
```

Autoriser l’aller :

```bash
nft add rule ip filter forward iifname "veth-r1" oifname "veth-r2" accept
```

Autoriser le retour grâce à conntrack :

```bash
nft add rule ip filter forward ct state established,related accept
```

---

# 9️⃣ Serveur TCP

```bash
sudo ip netns exec ns-server nc -l 8080
```

---

# 🔟 Client

```bash
sudo ip netns exec ns-client nc 192.168.0.2 8080
```

Écris du texte → il s’affiche côté serveur.

---

# 1️⃣1️⃣ Observer conntrack (le cœur du NAT)

```bash
sudo ip netns exec ns-router conntrack -L
```

Tu verras :

```
tcp 6 431999 ESTABLISHED
src=10.0.0.2 dst=192.168.0.2 sport=54321 dport=8080
src=192.168.0.2 dst=192.168.0.1 sport=8080 dport=45012
```

---

# 🔬 Ce que tu viens de voir

| Sens             | IP réelle sur le fil                     |
| ---------------- | ---------------------------------------- |
| client → serveur | **192.168.0.1:45012 → 192.168.0.2:8080** |
| serveur → client | **192.168.0.2:8080 → 192.168.0.1:45012** |

Conntrack sait que `192.168.0.1:45012` doit redevenir `10.0.0.2:54321`.

---

# 🧠 Phrase clé

> Fedora ne “fait pas du NAT”.
>
> Fedora **programme conntrack + nftables**, et c’est le kernel Linux qui fait toute la magie.
exp