| Famille socket         | Niveau           | Usage                     | Géré par Go ? | Lib Go                         |
| ---------------------- | ---------------- | ------------------------- | ------------- | ------------------------------ |
| `AF_INET` / `AF_INET6` | L3 / L4          | TCP, UDP, ICMP            | ✅ oui         | `net`, `net/http`, `grpc`, etc |
| `AF_UNIX`              | IPC local        | IPC filesystem / abstract | ✅ oui         | `net.UnixConn`                 |
| `AF_NETLINK`           | Kernel API       | contrôle noyau (ip, ss)   | ❌ stdlib      | `github.com/mdlayher/netlink`  |
| `AF_PACKET`            | L2 Ethernet      | capture / injection       | ❌ stdlib      | ❌ → `syscall` uniquement       |
| `AF_XDP`               | L2 zero-copy NIC | eBPF / XDP DMA            | ❌ stdlib      | ❌ → `syscall`                  |
| `AF_VSOCK`             | VM / host IPC    | virtio                    | ❌ stdlib      | `github.com/mdlayher/vsock`    |
| `AF_CAN`               | CAN bus          | automobile                | ❌ stdlib      | `github.com/linux-can/can-go`  |
| `AF_BLUETOOTH`         | Bluetooth L2CAP  | BT stack                  | ❌ stdlib      | `tinygo.org/x/bluetooth`       |
| `AF_IB`                | InfiniBand       | RDMA                      | ❌ stdlib      | wrappers RDMA                  |

---

AF_XDP est **le mécanisme réseau le plus rapide jamais ajouté à Linux**.
Et non : ce n’est **pas réservé à eBPF**, mais **il en dépend pour fonctionner**.

---

# ⚡ C’est quoi AF_XDP ?

`AF_XDP` est une **famille de sockets Linux** qui permet :

> d’accéder **directement aux buffers DMA de la carte réseau** depuis un process userland,
> avec **0 copie, 0 syscall par paquet**, latence minimale.

C’est le mode utilisé par les firewalls 100 Gb/s, DPDK-like **sans DPDK**.

---

# 🧬 Architecture

```
NIC ──DMA──► UMEM partagé ──► socket AF_XDP ──► userland
        ▲
        │
      eBPF XDP hook
```

---

# 🧩 Pourquoi eBPF est obligatoire

La carte ne sait pas **où envoyer les paquets** toute seule.

C’est le programme **XDP (eBPF)** qui s’attache à l’interface et décide :

* `XDP_PASS` → pile réseau normale
* `XDP_DROP` → drop kernel
* `XDP_REDIRECT` → **AF_XDP socket**

Sans eBPF XDP, **AF_XDP ne reçoit aucun paquet**.

---

# 🔌 Socket AF_XDP

```c
socket(AF_XDP, SOCK_RAW, 0)
```

Mais ce socket est inutile sans :

* UMEM (zone DMA partagée)
* rings RX/TX
* programme eBPF XDP attaché à la NIC

---

# 🚀 Différence avec AF_PACKET

| AF_PACKET          | AF_XDP              |
| ------------------ | ------------------- |
| après stack réseau | avant stack réseau  |
| 1–5 Mpps           | 50+ Mpps            |
| copie minimale     | zero-copy total     |
| latence élevée     | latence ultra basse |

---

# 🧠 Conclusion

AF_XDP est un **DMA userspace programmable**.

Ce n’est pas “juste un socket” :
c’est un **nouveau plan de données réseau Linux**, piloté par eBPF.

Tu peux l’utiliser en Go…
mais c’est un projet kernel-level (plus complexe que tcpdump).
ethtool -i eth0