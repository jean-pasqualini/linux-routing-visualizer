
Userland (libpcap / tcpdump)
│
▼
AF_PACKET socket (PF_PACKET)
│
▼
PACKET_MMAP (ring buffer partagé)
│
▼
TPACKET_V3 (zero-copy kernel → user)
│
▼
NIC driver → DMA

----

Super question — l’évolution de **libpcap**, c’est littéralement l’histoire de la capture réseau sous Unix.

---

# 🕰️ L’évolution des mécanismes libpcap

---

## 1993 – 1999 : **BSD BPF device**

### Socket / API

```
/dev/bpfX   (character device)
```

* ioctl sur un device
* kernel buffer copié vers userland
* filtrage BPF interprété

### Avantages

* simple
* portable BSD

### Limites

* très lent
* 1 paquet = 1 syscall
* pertes à haut débit

---

## 2000 – 2003 : **Linux RAW socket**

```c
socket(AF_INET, SOCK_RAW, IPPROTO_TCP)
```

* niveau L3
* pas de VLAN
* filtrage userland

➡️ Très vite abandonné : trop lent.

---

## 2003 – 2008 : **PF_PACKET + recvmsg()**

```c
socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL))
```

* capture niveau Ethernet
* toujours `recvmsg()`
* copies mémoire multiples

➡️ plafonne à ~200k pps.

---

## 2008 – 2013 : **PACKET_MMAP + TPACKET_V1/V2**

```c
setsockopt(fd, SOL_PACKET, PACKET_RX_RING)
mmap()
```

* ring buffer partagé kernel/user
* batch de paquets
* énorme gain perf (1–5 Mpps)

---

## 2014 – aujourd’hui : **TPACKET_V3 + eBPF JIT**

* TPACKET_V3
* eBPF compilé JIT
* multi-queue NIC
* zero-copy total
* NUMA aware

➡️ 40+ Gbit/s sans perte.

---

## 🌐 Par OS

| OS      | Mécanisme                     |
| ------- | ----------------------------- |
| Linux   | PF_PACKET + PACKET_MMAP       |
| FreeBSD | `/dev/bpf`                    |
| macOS   | `/dev/bpf`                    |
| Windows | Npcap / WinPcap (NDIS driver) |
| Solaris | DLPI                          |

---

## 📈 Résumé visuel

```
1993   /dev/bpf
2000   AF_INET raw socket
2003   AF_PACKET recvmsg
2008   PACKET_MMAP V1/V2
2014   TPACKET_V3 + eBPF
2025   eBPF offload NIC
```

---

## 🧠 Moralité

libpcap n’est plus une “lib de capture”.

C’est aujourd’hui un **driver DMA userspace ultra-optimisé** caché derrière une API C de 30 ans.


-----

C’est une super remarque, parce que **libpcap n’est pas un driver au sens “module noyau”**, mais son fonctionnement moderne est **identique à celui d’un driver DMA userspace**.

Je t’explique précisément pourquoi.

---

## 🔌 Le chemin réel d’un paquet réseau aujourd’hui

Quand un paquet arrive sur la carte réseau :

```
NIC ──DMA──► ring buffer kernel ──► ring buffer PACKET_MMAP ──► libpcap
```

Avec TPACKET_V3, ce *ring buffer PACKET_MMAP* est **mappé directement dans l’espace mémoire du process libpcap**.

Donc le CPU ne fait **aucune copie** du paquet.

---

## 🧠 Ce que fait normalement un driver

Un driver réseau classique :

* programme la NIC avec des adresses mémoire DMA
* reçoit des paquets directement dans des buffers partagés
* gère des rings, offsets, flags, status bits

---

## 🧩 Ce que fait libpcap moderne

libpcap fait exactement pareil mais **depuis userland** :

| Driver kernel    | libpcap                     |
| ---------------- | --------------------------- |
| ring DMA NIC     | ring PACKET_MMAP            |
| descripteurs RX  | `struct tpacket_block_desc` |
| polling NIC      | `poll()` sur socket         |
| interruption NIC | wakeup epoll                |
| buffer partagé   | `mmap()`                    |

---

## 🧪 Exemple : buffer TPACKET_V3

Chaque block du ring contient :

```c
struct tpacket_block_desc {
    uint32_t version;
    uint32_t offset_to_priv;
    struct tpacket_hdr_v1 hdr;
};
```

Avec :

* timestamps kernel
* status bits OWNERSHIP (kernel ↔ userland)
* offsets vers les paquets

C’est **exactement une table de descripteurs DMA**.

---

## 🚀 Pourquoi c’est révolutionnaire

Avant PACKET_MMAP :

```
NIC → kernel buffer → copy_to_user() → libpcap
```

Aujourd’hui :

```
NIC → DMA → ring partagé mmap → libpcap
```

Le kernel ne “copie” plus les données, il **donne la propriété du buffer** à ton process.

---

## 🧠 Conclusion

libpcap moderne est :

> un moteur DMA programmable depuis userland,
> utilisant la NIC comme backend matériel,
> sans copie, sans syscall par paquet.

Ce n’est pas un driver Linux,
mais **dans son architecture, il se comporte comme un driver DMA userspace**.

---

Un **ring buffer** (ou *buffer circulaire*) est une zone mémoire organisée en **cercle**, utilisée pour faire transiter des données très rapidement entre deux acteurs : un producteur et un consommateur.

Dans le cas réseau :
➡️ le producteur = la carte réseau / le noyau
➡️ le consommateur = ton programme (libpcap, tcpdump…)

---

## 🌀 Principe simple

Imagine un tableau de cases :

```
[0][1][2][3][4][5][6][7]
```

On a deux pointeurs :

* `head` → où on **écrit**
* `tail` → où on **lit**

Quand on arrive à la fin, on **reboucle au début**.

```
head → avance → → → revient à 0
tail → avance → → → revient à 0
```

C’est pour ça qu’on parle de *ring* (anneau).

---

## 🔁 Exemple concret

État initial :

```
Buffer : [ ][ ][ ][ ][ ][ ][ ][ ]
head = 0   tail = 0
```

La carte réseau écrit 3 paquets :

```
Buffer : [A][B][C][ ][ ][ ][ ][ ]
head = 3   tail = 0
```

Ton programme lit 2 paquets :

```
Buffer : [ ][ ][C][ ][ ][ ][ ][ ]
head = 3   tail = 2
```

Puis la carte écrit encore 4 paquets (et reboucle) :

```
Buffer : [ ][ ][C][D][E][F][G][ ]
head = 7   tail = 2
```

---

## ⚙️ Pourquoi c’est ultra rapide

Parce que :

* pas d’allocation mémoire
* pas de déplacement de données
* juste des pointeurs qui avancent
* parfaitement adapté au DMA

---

## 🧠 Dans libpcap / PACKET_MMAP

Le ring buffer est :

* **alloué par le noyau**
* **mappé en mémoire (`mmap`) dans ton process**
* rempli directement par la NIC via DMA
* partagé kernel ↔ userland

Ton programme ne fait que :

1. regarder les flags du block
2. traiter les paquets
3. rendre le block au noyau

---

## 🚀 Résultat

| Sans ring buffer   | Avec ring buffer               |
| ------------------ | ------------------------------ |
| copie mémoire      | zéro copie                     |
| 1 syscall / paquet | 1 syscall / batch              |
| lent               | capable de saturer du 40Gbit/s |

---

## 🧩 En résumé

Un ring buffer est une **autoroute mémoire circulaire** qui permet à deux mondes (hardware / kernel / userland) d’échanger des données **à vitesse maximale sans copie**.

---

| Capability      | Rôle                                                               |
| --------------- | ------------------------------------------------------------------ |
| `cap_net_raw`   | Autorise la création de sockets `AF_PACKET / SOCK_RAW`             |
| `cap_net_admin` | Autorise `PACKET_RX_RING`, promiscuité, filtres eBPF, réglages NIC |
