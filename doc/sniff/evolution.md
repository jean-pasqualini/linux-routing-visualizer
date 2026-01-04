Voici la **meilleure stratégie de sniffing par génération de kernel Linux**, basée sur ce qui est réellement exploitable en production.

---

# 🕰️ Évolution des meilleures méthodes de capture

| Kernel Linux           | Meilleure méthode                         | Pourquoi                        |
| ---------------------- | ----------------------------------------- | ------------------------------- |
| ≤ 2.6.27 (avant 2008)  | `AF_PACKET + recvmsg()`                   | PACKET_MMAP pas encore mature   |
| 2.6.27 – 3.1           | `AF_PACKET + PACKET_MMAP (TPACKET_V1/V2)` | premiers rings partagés         |
| 3.2 – 3.18             | `AF_PACKET + TPACKET_V2`                  | stable, utilisé par tcpdump     |
| 3.19 – 4.17            | `AF_PACKET + TPACKET_V3`                  | batching, zero-copy, timestamps |
| ≥ 4.18                 | `AF_PACKET + TPACKET_V3 + eBPF filter`    | standard moderne libpcap        |
| ≥ 5.3 + NIC compatible | `AF_XDP zero-copy`                        | ultra-haut débit (40–100Gbps)   |

---

# 🧩 Détail par époque

## 🏺 Kernels anciens (< 2.6.27)

* seule option : `recvfrom()` sur socket `AF_PACKET`
* 1 syscall / paquet
* très lent
* aujourd’hui inutilisable

---

## 🧱 PACKET_MMAP (2008+)

```c
socket(AF_PACKET, SOCK_RAW, ETH_P_ALL)
setsockopt(PACKET_RX_RING)
mmap()
```

| ABI        | Usage                 |
| ---------- | --------------------- |
| TPACKET_V1 | obsolète              |
| TPACKET_V2 | stable                |
| TPACKET_V3 | **actuel recommandé** |

C’est ce que libpcap utilise aujourd’hui.

---

## ⚡ AF_XDP (2019+)

```c
socket(AF_XDP, SOCK_RAW, 0)
```

* programme XDP eBPF obligatoire :

```
NIC → XDP hook → AF_XDP UMEM → userspace
```

Avantages :

* zéro copie réelle
* latence minimale
* 50M+ pps

Inconvénients :

* dépend fortement du matériel
* nécessite eBPF + root
* beaucoup plus complexe

---

# 🎯 Recommandation claire

| Objectif                          | Solution                              |
| --------------------------------- | ------------------------------------- |
| sniffer universel (comme tcpdump) | **AF_PACKET + TPACKET_V3**            |
| IDS / load balancer 100Gbps       | **AF_XDP**                            |
| portable / simple                 | libpcap                               |
| sniffer Go sans cgo               | **AF_PACKET + TPACKET_V3 en syscall** |

---

# 🧠 Conclusion

Il n’existe pas **une seule meilleure méthode**.

Mais pour 95% des usages :

> **AF_PACKET + TPACKET_V3 est le meilleur compromis performance / compatibilité / stabilité depuis 2014.**
