Voici la liste complète des “tables” possibles, avec leur usage réel, en distinguant iptables et nftables (important pour éviter les confusions).

⸻

1️⃣ Tables iptables (modèle historique)

Ces tables sont fixes et ont un rôle précis dans le pipeline Netfilter.

Table	Rôle principal	Hooks
filter	Décision finale (ACCEPT / DROP)	INPUT, FORWARD, OUTPUT
nat	Traduction d’adresses (NAT)	PREROUTING, OUTPUT, POSTROUTING
mangle	Modification des paquets	tous
raw	Avant conntrack	PREROUTING, OUTPUT
security	LSM / SELinux	INPUT, FORWARD, OUTPUT

👉 En iptables-legacy, elles existent toutes même si vides.

⸻

2️⃣ Tables nftables (modèle moderne)

⚠️ En nftables, une “table” est juste un conteneur logique.
Elle n’existe que si tu la crées.

Familles possibles
•	ip (IPv4)
•	ip6 (IPv6)
•	inet (IPv4 + IPv6)
•	arp
•	bridge
•	netdev

⸻

Tables courantes (par convention)

Table (nom libre)	Usage typique
filter	Filtrage
nat	NAT
mangle	QoS / MARK / DSCP
raw	notrack / early
security	SELinux

👉 Ces noms imitent iptables, mais ne sont pas obligatoires.

⸻

Ce qui compte vraiment en nftables

Les chaînes avec hook

C’est là que le comportement est défini :

chain input {
type filter hook input priority 0;
}

	•	type → rôle (filter, nat, route)
	•	hook → moment du pipeline
	•	priority → ordre précis

➡️ Le nom de la table n’a aucune importance fonctionnelle.

⸻

3️⃣ Correspondance iptables → nftables

iptables	nftables
table filter	chain type filter
table nat	chain type nat
table mangle	chain type filter + priorité
table raw	chain filter priorité très basse
table security	chain filter priorité haute


⸻

4️⃣ Pourquoi tu ne vois que filter et nat

Parce que :
•	ce sont les seules tables nécessaires
•	Docker / système les créent
•	les autres ne sont pas utiles par défaut

⸻

5️⃣ TL;DR
•	iptables → 5 tables fixes
•	nftables → tables libres
•	ce sont les chaînes + hooks qui comptent
•	filter et nat suffisent dans 90% des cas

⸻

Si tu veux, je peux :
•	dessiner le pipeline complet
•	montrer une config nft minimale
•	traduire chaque table iptables en nftables natif

Dis-moi 👍