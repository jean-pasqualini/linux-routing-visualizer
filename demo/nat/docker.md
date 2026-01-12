Dans Docker, tu retrouves **exactement la même topologie que ton lab NAT Fedora**, juste automatisée.

Je te décris la configuration réelle, pas le marketing Docker.

---

# 🌍 Architecture réseau Docker par défaut

Quand tu fais :

```bash
docker run nginx
```

Docker construit :

```
[ container ]
   172.17.0.2
       │
     veth
       │
[ docker0 bridge ]
   172.17.0.1
       │
   HOST LINUX  ── NAT ── Internet
```

---

# 📦 Les sous-réseaux Docker

| Zone       | Sous-réseau               |
| ---------- | ------------------------- |
| Containers | `172.17.0.0/16`           |
| Host LAN   | `192.168.1.0/24` ou autre |
| WAN        | IP publique               |

---

# 🔄 Ce que fait Docker

### 1️⃣ Bridge L2

```
docker0 = switch Ethernet
```

Il ne route rien, il commute.

---

### 2️⃣ NAT POSTROUTING (MASQUERADE)

```bash
iptables -t nat -A POSTROUTING -s 172.17.0.0/16 ! -o docker0 -j MASQUERADE
```

➡️ Tout conteneur → SNAT vers IP du host.

---

### 3️⃣ DNAT -p 8080:80

```bash
iptables -t nat -A PREROUTING -p tcp --dport 8080 -j DNAT --to 172.17.0.2:80
```

---

### 4️⃣ Hairpin NAT

Quand un conteneur appelle :

```
curl http://localhost:8080
```

Docker applique DNAT + SNAT vers lui-même.

---

# 🧠 Docker = ton lab automatisé

| Ton lab Fedora | Docker    |
| -------------- | --------- |
| ns-client      | container |
| ns-router      | host      |
| ns-server      | Internet  |

C’est littéralement la même architecture kernel.
    