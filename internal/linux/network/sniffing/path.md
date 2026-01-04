AF_INET (L3)
NIC → driver → stack IP → raw socket → copy_to_user → ton programme
NET_FILTER / NFQUEUE (L3/L4)
NIC → driver → stack IP → netfilter hook → queue → copy_to_user
AF_PACKET (L2)
NIC → DMA → ring mmap partagé → userland

---
| Méthode                | Niveau | Perf   |
| ---------------------- | ------ | ------ |
| AF_PACKET + TPACKET_V3 | L2     | 🚀🚀🚀 |
| AF_INET raw socket     | L3     | 🐌     |
| NFQUEUE                | L3/L4  | 🐌🐌   |

----

| Outil                    | Méthode                                |
| ------------------------ | -------------------------------------- |
| **tcpdump**              | `AF_PACKET + PACKET_MMAP (TPACKET_V3)` |
| **Wireshark / tshark**   | libpcap → `AF_PACKET + PACKET_MMAP`    |
| **Suricata (fast path)** | `AF_PACKET` ou `AF_XDP`                |
| **Zeek (mode perf)**     | `AF_PACKET`                            |
| **Snort 3**              | `AF_PACKET`                            |
| **ntopng / nProbe**      | `AF_PACKET`                            |

----

Oui — et c’est même la **tendance majeure actuelle**.

Il existe aujourd’hui des sniffeurs qui utilisent **eBPF directement dans le noyau**, parfois **sans passer par AF_PACKET**.

---

## 🧬 Deux grandes familles

### 1️⃣ eBPF comme **filtre**

C’est ce que fait libpcap moderne :

```
NIC → kernel → eBPF filter → AF_PACKET ring → userland
```

Exemples :

* `tcpdump`
* `Wireshark`
* `tshark`

Ici eBPF remplace l’ancien BPF interprété, mais la capture reste AF_PACKET.

---

### 2️⃣ eBPF comme **moteur de capture (XDP)**

Ici tu es **avant la stack réseau** :

```
NIC → XDP hook → eBPF program → ring AF_XDP → userland
```

Exemples :

* Suricata XDP
* Katran (Facebook)
* Cilium monitor
* bpftool prog tracelog

Là, tu n’utilises **plus AF_PACKET du tout**.

---

## 🚀 Avantages

| AF_PACKET    | AF_XDP + eBPF            |
| ------------ | ------------------------ |
| millions pps | dizaines de millions pps |
| stack réseau | avant stack              |
| copie minime | zero-copy total          |
| universel    | matériel dépendant       |

---

## 🧠 Conclusion

Oui, les sniffeurs eBPF existent.

* libpcap = eBPF filtre
* XDP / AF_XDP = eBPF moteur data-plane

Tu es déjà à mi-chemin du monde eBPF.
