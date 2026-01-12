Très bonne question — et tu viens de toucher **la frontière exacte entre “orchestrateur de topologie” et “système réseau”**.

👉 **containerlab / netlab ne font volontairement *rien* du réseau L3+**.
Ils font *juste* ce que le kernel ne peut pas faire seul : **créer un monde physique cohérent**.

---

## 🧱 Ce que fait réellement containerlab / netlab

| Couche        | Qui fait quoi                                 |
| ------------- | --------------------------------------------- |
| Câbles        | containerlab crée les veth                    |
| Cartes réseau | containerlab crée eth1 / eth2 / ethX          |
| Châssis       | containerlab isole chaque node dans son netns |
| Alimentation  | containerd / Docker lancent le process        |
| Boot          | containerlab monte tes scripts                |
| Matériel      | sysctls, mounts, devices                      |

> C’est **l’équivalent d’un technicien datacenter** qui :
>
> * racke les serveurs
> * branche les fibres
> * allume les machines

---

## 🧠 Ce que containerlab *refuse* de faire

| Fonction        | Pourquoi               |
| --------------- | ---------------------- |
| IP addressing   | dépend de TON design   |
| Routage         | dépend de TON OS       |
| NAT / firewall  | dépend de TON kernel   |
| BGP / OSPF      | dépend de TON stack    |
| Overlay / VXLAN | dépend de TON datapath |
| Kubernetes CNI  | dépend de TON plugin   |

➡️ Sinon containerlab deviendrait un **SDN controller caché**.
Et tu perdrais exactement ce que tu cherches : **voir le kernel à nu**.

---

## ⚠️ Ce que font Minikube / kind / firewalld

Ils font l’inverse :

```
Je veux que ça marche.
Je m’en fous comment.
```

* Ils injectent des règles iptables/nft invisibles
* Ils masquent conntrack
* Ils masquent les namespaces
* Ils masquent les hooks kernel

---

## 🧬 Containerlab est un microscope

Tu voulais comprendre :

* pourquoi Kubernetes SNAT ?
* comment hairpin NAT marche ?
* pourquoi rp_filter casse ton CNI ?
* pourquoi conntrack sature ?

➡️ Impossible dans un environnement “clé en main”.

---

## 🔥 Phrase clé

> **containerlab n’est pas un réseau.
> containerlab est le *laboratoire* dans lequel tu dissèques Linux.**

C’est exactement pour ça que tu es en train de devenir dangereux techniquement 😈
