Parfait 👍
Voici la traduction exacte et complète iptables → hooks Netfilter, en mode copiable / mémo.

⸻

📘 Traduction iptables → Netfilter hooks

iptables = vue utilisateur
Netfilter = réalité noyau

⸻

🔥 Hooks Netfilter (rappel)

NF_INET_PRE_ROUTING
NF_INET_LOCAL_IN
NF_INET_FORWARD
NF_INET_LOCAL_OUT
NF_INET_POST_ROUTING


⸻

🧱 Table raw

iptables -t raw

Chaîne iptables	Hook Netfilter
PREROUTING	NF_INET_PRE_ROUTING
OUTPUT	NF_INET_LOCAL_OUT

Priorité : très basse (avant conntrack)

⸻

🧱 Table mangle

iptables -t mangle

Chaîne iptables	Hook Netfilter
PREROUTING	NF_INET_PRE_ROUTING
INPUT	NF_INET_LOCAL_IN
FORWARD	NF_INET_FORWARD
OUTPUT	NF_INET_LOCAL_OUT
POSTROUTING	NF_INET_POST_ROUTING


⸻

🧱 Table nat

iptables -t nat

Chaîne iptables	Hook Netfilter
PREROUTING	NF_INET_PRE_ROUTING
OUTPUT	NF_INET_LOCAL_OUT
POSTROUTING	NF_INET_POST_ROUTING

⚠️ Appliquée uniquement au premier paquet d’une connexion

⸻

🧱 Table filter

iptables -t filter

Chaîne iptables	Hook Netfilter
INPUT	NF_INET_LOCAL_IN
FORWARD	NF_INET_FORWARD
OUTPUT	NF_INET_LOCAL_OUT


⸻

🧱 Table security

iptables -t security

Chaîne iptables	Hook Netfilter
INPUT	NF_INET_LOCAL_IN
FORWARD	NF_INET_FORWARD
OUTPUT	NF_INET_LOCAL_OUT


⸻

🧠 Vision globale (pipeline)

PREROUTING:
raw     → mangle → nat
(NF_INET_PRE_ROUTING)

ROUTING DECISION

INPUT:
mangle → filter → security
(NF_INET_LOCAL_IN)

FORWARD:
mangle → filter → security
(NF_INET_FORWARD)

OUTPUT:
raw → mangle → nat → filter → security
(NF_INET_LOCAL_OUT)

POSTROUTING:
mangle → nat
(NF_INET_POST_ROUTING)


⸻

🔑 À retenir absolument
•	iptables regroupe les hooks
•	Netfilter n’a que des hooks
•	tables = priorités + hooks
•	nftables expose ça directement

⸻

TL;DR

iptables tables/chaînes
↓
Netfilter hooks + priorité


⸻

Si tu veux, je peux :
•	faire la traduction inverse hooks → iptables
•	te donner la version nftables équivalente
•	expliquer pourquoi NAT n’est qu’au 1er paquet

Dis-moi 👍