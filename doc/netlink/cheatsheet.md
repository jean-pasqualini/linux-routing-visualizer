Voici ton **cheatsheet Netlink Linux** – clair, pratique, orienté dev système / réseau.

---

# 🧠 Netlink – C’est quoi ?

Netlink est un **socket kernel ↔ userspace** pour :

* interfaces réseau
* adresses IP
* routes
* règles de routage
* voisins (ARP/NDP)
* traffic control
* monitoring temps réel

Alternative moderne à `ioctl`, `proc`, `sysfs`.

---

# 📡 Familles Netlink

| Famille                  | Utilité                         |
| ------------------------ | ------------------------------- |
| `NETLINK_ROUTE`          | Interfaces, IP, routes, ARP, tc |
| `NETLINK_GENERIC`        | Protocoles custom kernel        |
| `NETLINK_NETFILTER`      | Conntrack, firewall             |
| `NETLINK_AUDIT`          | Audit                           |
| `NETLINK_KOBJECT_UEVENT` | Udev events                     |
| `NETLINK_SOCK_DIAG`      | Infos sockets                   |
| `NETLINK_XFRM`           | IPsec                           |
| `NETLINK_CRYPTO`         | Crypto API                      |

---

# 🗺 Ce que tu peux faire avec `NETLINK_ROUTE`

| Objet      | Commande ip | Message                        |
| ---------- | ----------- | ------------------------------ |
| Interfaces | `ip link`   | `RTM_GETLINK`, `RTM_NEWLINK`   |
| IP         | `ip addr`   | `RTM_GETADDR`, `RTM_NEWADDR`   |
| Routes     | `ip route`  | `RTM_GETROUTE`, `RTM_NEWROUTE` |
| ARP / NDP  | `ip neigh`  | `RTM_GETNEIGH`                 |
| Rules      | `ip rule`   | `RTM_GETRULE`                  |
| TC         | `tc qdisc`  | `RTM_NEWQDISC`                 |

---

# 🔧 Création socket netlink

```c
int fd = socket(AF_NETLINK, SOCK_RAW, NETLINK_ROUTE);

struct sockaddr_nl local = {
    .nl_family = AF_NETLINK,
    .nl_pid = getpid(),
    .nl_groups = RTMGRP_LINK | RTMGRP_IPV4_IFADDR | RTMGRP_IPV4_ROUTE
};

bind(fd, (struct sockaddr*)&local, sizeof(local));
```

---

# 📥 Abonnement aux events

| Groupe               | Event                |
| -------------------- | -------------------- |
| `RTMGRP_LINK`        | UP/DOWN interfaces   |
| `RTMGRP_IPV4_IFADDR` | IP ajout/suppression |
| `RTMGRP_IPV4_ROUTE`  | Routes modifiées     |
| `RTMGRP_NEIGH`       | ARP                  |

---

# 📤 Exemple requête dump routes

```c
struct {
  struct nlmsghdr nlh;
  struct rtmsg rtm;
} req;

req.nlh.nlmsg_len = NLMSG_LENGTH(sizeof(struct rtmsg));
req.nlh.nlmsg_type = RTM_GETROUTE;
req.nlh.nlmsg_flags = NLM_F_REQUEST | NLM_F_DUMP;
req.rtm.rtm_family = AF_INET;
send(fd, &req, req.nlh.nlmsg_len, 0);
```

---

# 📦 Parser message Netlink

```c
for (nlh = (struct nlmsghdr*)buf; NLMSG_OK(nlh, len); nlh = NLMSG_NEXT(nlh, len)) {
  if (nlh->nlmsg_type == RTM_NEWROUTE) {
     struct rtmsg *rtm = NLMSG_DATA(nlh);
     struct rtattr *attr = RTM_RTA(rtm);
     int attrlen = RTM_PAYLOAD(nlh);

     for (; RTA_OK(attr, attrlen); attr = RTA_NEXT(attr, attrlen)) {
        if (attr->rta_type == RTA_DST) ...
        if (attr->rta_type == RTA_GATEWAY) ...
     }
  }
}
```

---

# 🧩 RT Attributes importants

| Objet            | Attribut       |
| ---------------- | -------------- |
| Route dest       | `RTA_DST`      |
| Gateway          | `RTA_GATEWAY`  |
| Interface sortie | `RTA_OIF`      |
| Prefsrc          | `RTA_PREFSRC`  |
| Table            | `RTA_TABLE`    |
| Metric           | `RTA_PRIORITY` |

---

# 🕵️‍♂️ Outils debug

```bash
ss -nlp
ip monitor
strace -e sendto,recvfrom ip route
```

---

# 📚 Lib utiles

| Lib           | Usage            |
| ------------- | ---------------- |
| `libnl`       | API haut niveau  |
| `pyroute2`    | Netlink Python   |
| `bpftool net` | lien BPF/netlink |

---

# 🧠 Cas d’usage

| Besoin                       | Netlink             |
| ---------------------------- | ------------------- |
| Firewall dynamique           | `NETLINK_NETFILTER` |
| Monitor réseau temps réel    | `RTMGRP_*`          |
| Manipulation routing         | `RTM_NEWROUTE`      |
| Création proto kernel custom | `NETLINK_GENERIC`   |
| Traffic shaping              | `RTM_NEWQDISC`      |

---

Netlink = **l’API réseau réelle de Linux**.
Tout ce que fait `iproute2`, NetworkManager ou Kubernetes passe par là.
