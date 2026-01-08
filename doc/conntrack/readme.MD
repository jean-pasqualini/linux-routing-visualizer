Tu mets exactement le doigt sur un point très mal compris :

conntrack n’utilise PAS les états iptables NEW / ESTABLISHED / RELATED.

Ce sont deux couches différentes.

⸻

1️⃣ États conntrack (noyau)

C’est ce que tu vois dans flow.Status :

Flag noyau	Rôle réel
SEEN_REPLY	On a vu un paquet retour
ASSURED	Connexion considérée fiable
CONFIRMED	Insérée définitivement
EXPECTED	Connexion “annexe” (FTP, SIP…)

Ce sont des bits internes, pas des états firewall.

⸻

2️⃣ États iptables (-m state / -m conntrack)

Quand tu fais :

-m conntrack --ctstate ESTABLISHED

Netfilter traduit en réalité :

SEEN_REPLY && ASSURED && CONFIRMED

Ce n’est qu’une projection logique.

⸻

Mapping réel

iptables	noyau
NEW	pas encore CONFIRMED
ESTABLISHED	SEEN_REPLY
RELATED	EXPECTED
INVALID	incohérent / hors suivi


⸻

Pourquoi ils ont inventé ESTABLISHED ?

Pour rendre ça humain.

Mais la vérité est au niveau du bitmask noyau.

⸻

Conclusion

ESTABLISHED n’existe pas dans conntrack.
C’est un raccourci firewall, pas un état réel.