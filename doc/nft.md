Excellente question — et elle montre que tu commences à voir le modèle mental derrière Netfilter 👌
La réponse courte :

nftables rend le pipeline explicite, alors qu’iptables le cache.

Je t’explique précisément pourquoi, avec des exemples concrets.

⸻

1️⃣ Ce qu’iptables te cache

En iptables, le pipeline est implicite.

Quand tu écris :

iptables -t nat -A PREROUTING -p tcp --dport 80 -j DNAT

Tu ne vois nulle part :
•	quand cette règle est exécutée
•	dans quel ordre exact
•	par rapport aux autres tables

👉 Tu dois le savoir par cœur.

⸻

2️⃣ Ce que nftables rend explicite

En nftables, tout est écrit dans la règle.

chain prerouting {
type nat hook prerouting priority -100;
}

Tu vois immédiatement :

Élément	Sens
type nat	rôle (NAT)
hook prerouting	moment exact
priority -100	ordre précis

👉 Plus de magie.

⸻

3️⃣ Plus de confusion “table = phase”

En iptables :
•	raw, mangle, nat, filter
•	semblent être des “choix”

En nftables :
•	le hook définit la phase
•	la table n’est qu’un conteneur

Tu peux écrire :

table inet firewall {
chain input {
type filter hook input priority 0;
}
}

➡️ Le nom firewall est libre
➡️ La phase est clairement hook input

⸻

4️⃣ Le pipeline devient lisible

iptables (mentalement)

raw → mangle → nat → filter → security

nftables (réellement visible)

hook prerouting priority -300
hook prerouting priority -150
hook prerouting priority -100
hook input priority 0

👉 Tu vois l’ordre exact.

⸻

5️⃣ Plusieurs chaînes au même hook (impossible en iptables)

En nftables :

chain early {
type filter hook prerouting priority -300;
}

chain later {
type filter hook prerouting priority -100;
}

Tu vois clairement qui passe avant qui.

⸻

6️⃣ Actions multiples, règles lisibles

iptables :

-j LOG
-j DROP

nftables :

log prefix "DROP " counter drop

👉 Le comportement est évident.

⸻

7️⃣ TRACE beaucoup plus clair

iptables :
•	logs noyau cryptiques

nftables :

nft monitor trace

👉 Sortie lisible :
•	table
•	chaîne
•	règle exacte
•	verdict

⸻

8️⃣ Résumé visuel

Aspect	iptables	nftables
Pipeline	caché	explicite
Ordre	implicite	priority
Phase	table	hook
Lisibilité	faible	élevée
Debug	pénible	propre


⸻

9️⃣ TL;DR (à graver)
•	iptables = magie + conventions
•	nftables = déclaration explicite
•	le paquet ne choisit rien
•	le kernel suit le hook + priorité

⸻

Si tu veux, je peux :
•	te montrer la même règle iptables et nft côte à côte
•	dessiner le pipeline nftables complet
•	t’aider à penser nftables “nativement”

Dis-moi 👍