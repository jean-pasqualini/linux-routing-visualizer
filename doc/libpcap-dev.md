Très bonne question — c’est là que ça devient vraiment intéressant 😊
`libpcap` ne “sniffe” pas par magie : elle s’appuie **directement sur le noyau** et sur des mécanismes très bas niveau, différents selon l’OS.

Je te détaille le **vrai pipeline interne**.

---

## 1️⃣ Le principe fondamental

Une carte réseau reçoit des trames Ethernet.
Normalement, le noyau ne remet à ton programme **que les paquets qui lui sont destinés**.

libpcap va :

1. Mettre l’interface en **mode promiscuous** (optionnel)
2. S’accrocher **dans la pile réseau du noyau**
3. Copier chaque paquet **avant qu’il soit traité par TCP/IP**
4. Appliquer un **filtre BPF en kernel-space**
5. Remonter seulement les paquets filtrés vers ton programme

---

## 2️⃣ Sous Linux : PF_PACKET + mmap ring buffer

Aujourd’hui libpcap utilise principalement :

```
AF_PACKET / PF_PACKET sockets
```

C’est une API Linux qui permet d’accéder aux trames **au niveau Ethernet**.

### Étapes internes

```
Carte réseau
   ↓
Driver NIC
   ↓
Kernel (RX ring buffer)
   ↓
AF_PACKET socket
   ↓
BPF filter (dans le kernel)
   ↓
mmap() ring buffer partagé
   ↓
libpcap
   ↓
tcpdump / Wireshark / ton sniffeur
```

### Code noyau (simplifié)

libpcap fait en gros :

```c
int fd = socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
setsockopt(fd, SOL_PACKET, PACKET_VERSION, TPACKET_V3);
setsockopt(fd, SOL_PACKET, PACKET_RX_RING, &ring_params);
mmap(fd, ring_buffer);
```

Résultat :
Le noyau écrit directement les paquets dans une mémoire partagée → **zéro copie**, performances énormes.

---

## 3️⃣ Le filtre BPF : clé de la performance

Quand tu fais :

```bash
tcpdump "tcp and port 443"
```

libpcap compile ce filtre en **bytecode BPF** :

```
EtherType == IP
→ IP.protocol == TCP
→ TCP.dst_port == 443
```

Ce programme BPF est injecté dans le noyau :

```c
setsockopt(fd, SO_ATTACH_FILTER, &bpf_prog);
```

👉 Le noyau jette les paquets non pertinents **avant même qu’ils arrivent en userland**.

Sans ça, Wireshark sur un lien 10Gb/s serait inutilisable.

---

## 4️⃣ Mode promiscuous

```c
struct packet_mreq mreq;
mreq.mr_type = PACKET_MR_PROMISC;
setsockopt(fd, SOL_PACKET, PACKET_ADD_MEMBERSHIP, &mreq);
```

Ça dit à la carte réseau :

> « Envoie-moi TOUT ce que tu vois sur le fil, pas seulement mon MAC. »

---

## 5️⃣ Autres OS

| OS          | Implémentation                      |
| ----------- | ----------------------------------- |
| Linux       | `AF_PACKET` + `TPACKET_V3` + mmap   |
| BSD / macOS | `/dev/bpf` (Berkeley Packet Filter) |
| Windows     | Npcap → driver NDIS kernel          |
| Solaris     | DLPI                                |
| AIX / HP-UX | modules kernel propriétaires        |

---

## 6️⃣ Pourquoi il faut être root ?

Parce que tu ouvres un **socket raw Ethernet** :

```c
socket(AF_PACKET, SOCK_RAW, ...)
```

Ça donne accès aux trames brutes → équivalent d’une carte d’écoute réseau matérielle.
C’est donc restreint pour des raisons de sécurité.

---

## 7️⃣ Résumé en une phrase

`libpcap` ne lit pas “le réseau” depuis l’OS :
elle se branche **dans la couche la plus basse du noyau**, intercepte les trames Ethernet via des sockets spéciaux, applique un filtre BPF **dans le kernel**, et te donne un flux ultra-rapide de paquets avant qu’ils soient traités par TCP/IP.

C’est pour ça que tcpdump, Wireshark et ton futur sniffeur Go ont exactement la même puissance.


----

Voici un **résumé clair pour tes notes** sur le sniffing réseau sous Linux et libpcap.

---

# 🧠 libpcap & sniffing réseau — résumé

## 1️⃣ Rôle de libpcap

`libpcap` est la bibliothèque bas niveau utilisée par :

* tcpdump
* Wireshark / tshark
* gopacket / scapy / nmap (partiellement)

Elle intercepte les paquets **dans le noyau**, avant leur traitement par TCP/IP.

---

## 2️⃣ Pipeline de capture

```
Carte réseau
  ↓
Driver NIC
  ↓
RX ring buffer kernel
  ↓
AF_PACKET socket
  ↓
BPF filter (dans le kernel)
  ↓
mmap ring buffer partagé
  ↓
libpcap
  ↓
tcpdump / Wireshark / ton programme
```

---

## 3️⃣ AF_PACKET

* Famille de sockets Linux pour accéder au **niveau Ethernet (L2)**
* Création :

```c
socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
```

* Permet de :

    * voir ARP, VLAN, LLDP, etc.
    * injecter des trames Ethernet
    * sniffer avant firewall / IP

* Mode promiscuous :

```c
PACKET_MR_PROMISC
```

---

## 4️⃣ Filtrage BPF

Les filtres (`tcp and port 443`) sont compilés en **bytecode BPF** injecté dans le noyau :

```c
setsockopt(fd, SO_ATTACH_FILTER, &bpf_prog);
```

👉 Les paquets inutiles sont rejetés **avant l’userspace**.

---

## 5️⃣ Versions TPACKET

`TPACKET` est le format mémoire interne d’AF_PACKET.

| Version | Caractéristiques                                   |
| ------- | -------------------------------------------------- |
| **V1**  | 1 paquet = 1 frame, timestamps faibles             |
| **V2**  | timestamps ns, VLAN, rxhash                        |
| **V3**  | batching multi-paquets par block, ultra haute perf |

### Différences clés

|                 | V1        | V2      | V3       |
| --------------- | --------- | ------- | -------- |
| Zéro copie mmap | ✅         | ✅       | ✅        |
| Timestamp ns    | ❌         | ✅       | ✅        |
| Batching        | ❌         | ❌       | ✅        |
| Perf max        | ~300k pps | ~1M pps | >10M pps |

---

## 6️⃣ AF_PACKET ≠ TPACKET_V3

* **AF_PACKET** = API noyau
* **TPACKET_V3** = format mémoire haute performance utilisé par AF_PACKET

AF_PACKET n’est pas déprécié — il a simplement évolué.

---

## 7️⃣ Anciennes implémentations

Avant AF_PACKET :

| Période      | Mécanisme                |
| ------------ | ------------------------ |
| BSD / macOS  | `/dev/bpf`               |
| Linux 2.x    | `SOCK_PACKET` (supprimé) |
| Linux ancien | raw sockets IP           |

---

## 8️⃣ Pourquoi root est nécessaire

Parce que tu ouvres un socket Ethernet brut :

```c
socket(AF_PACKET, SOCK_RAW, ...)
```

👉 accès direct aux trames réseau = privilège noyau.

---

## 9️⃣ Conclusion

`libpcap` n’est qu’un **wrapper élégant** autour d’AF_PACKET.
AF_PACKET + BPF + TPACKET_V3 = le moteur qui permet à Wireshark et tcpdump de sniffer des liens 10–40 Gb/s en temps réel.


----

Aujourd’hui les deux coexistent — mais **libpcap n’a absolument pas disparu**.

La vraie réponse est :

> **libpcap reste la couche de capture standard**,
> **eBPF est utilisé comme accélérateur / moteur d’analyse avancé.**

---

## 1️⃣ Modèle classique (toujours dominant)

```
NIC → driver → AF_PACKET → TPACKET_V3 → libpcap → userspace
```

Utilisé par :

* Wireshark
* tcpdump
* Zeek (mode standard)
* Suricata (mode pcap)

Avantages :

* portable
* stable
* simple
* très bien optimisé

---

## 2️⃣ Modèle moderne eBPF

```
NIC → driver → eBPF (kernel) → ring buffer / perf events → userspace
```

On accroche un programme eBPF sur :

* XDP (avant la pile réseau)
* TC ingress / egress
* socket filters

Le programme eBPF :

* filtre
* agrège
* transforme
* n’envoie à l’userspace que des **événements utiles**

---

## 3️⃣ Pourquoi eBPF est apparu

libpcap te donne :

* **tout** le trafic
* analyse en userspace
* coûteux à 40–100 Gb/s

eBPF permet :

* filtrage sémantique (DNS, HTTP, TLS SNI…)
* réduction massive du volume
* analyse temps réel kernel-side

---

## 4️⃣ Exemples réels

| Outil             | Techno           |
| ----------------- | ---------------- |
| tcpdump           | libpcap          |
| Wireshark         | libpcap          |
| Zeek (haute perf) | AF_PACKET + eBPF |
| Suricata XDP      | eBPF             |
| Cilium / Pixie    | eBPF only        |
| Falco             | eBPF             |

---

## 5️⃣ Pourquoi libpcap est encore partout

| libpcap                     | eBPF                          |
| --------------------------- | ----------------------------- |
| portable                    | Linux only                    |
| simple                      | très complexe                 |
| analyse complète possible   | logique limitée               |
| pas de vérification stricte | vérificateur BPF contraignant |

---

## 6️⃣ Architecture hybride moderne

De plus en plus d’outils font :

```
NIC → eBPF filtre primaire → AF_PACKET/libpcap → analyse fine userspace
```

👉 eBPF enlève 90% du bruit
👉 libpcap garde la richesse protocolaire

---

## 7️⃣ Conclusion

Les sniffeurs modernes **n’ont pas remplacé libpcap**.

Ils l’**entourent avec eBPF** pour survivre à l’ère du 100 Gb/s.


----
|                    | BPF classique (libpcap) | eBPF                                 |
| ------------------ | ----------------------- | ------------------------------------ |
| Lieu               | socket AF_PACKET        | XDP / TC / kprobes                   |
| État entre paquets | ❌                       | ✅                                    |
| Filtrage           | simple (offsets)        | logique riche                        |
| Performance        | très bonne              | extrême                              |
| Complexité         | faible                  | élevée                               |
| Usage typique      | tcpdump, Wireshark      | observabilité cloud, IDS kernel-side |
