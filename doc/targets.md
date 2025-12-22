L’option -j (jump) peut prendre deux grandes catégories de valeurs :
1.	une chaîne (de la même table)
2.	une target (action)

Voici la liste complète, structurée et pratique.

⸻

1️⃣ -j <CHAÎNE>

👉 Saut vers une autre chaîne de la même table

Chaînes possibles
•	Chaînes builtin :
•	INPUT, OUTPUT, FORWARD
•	PREROUTING, POSTROUTING
•	Chaînes utilisateur :

iptables -N MY_CHAIN
iptables -A INPUT -j MY_CHAIN

⚠️ Jamais vers une autre table

⸻

2️⃣ -j <TARGET> (actions)

🔹 Targets de base (toujours disponibles)

Target	Effet
ACCEPT	Accepte le paquet, fin du traitement
DROP	Jette le paquet silencieusement
REJECT	Jette le paquet + réponse
RETURN	Retour à la chaîne appelante


⸻

🔹 Targets NAT (table nat uniquement)

Target	Rôle
DNAT	Change l’IP/port destination
SNAT	Change l’IP source
MASQUERADE	SNAT dynamique
REDIRECT	DNAT vers la machine locale


⸻

🔹 Targets de log / debug

Target	Description
LOG	Log kernel
NFLOG	Log via netfilter
TRACE	Trace complète (debug)


⸻

🔹 Targets de modification (mangle)

Target	Effet
MARK	Marque le paquet
CONNMARK	Marque la connexion
DSCP	Modifie le DSCP
TTL	Modifie le TTL


⸻

🔹 Targets avancées / spécialisées

Target	Usage
QUEUE	Envoi vers userspace
NFQUEUE	Userspace avancé
CT	Manipulation conntrack
NOTRACK	Désactive conntrack
SECMARK	Sécurité (SELinux)
AUDIT	Audit LSM
CLASSIFY	QoS


⸻

🔹 Targets dépendantes de modules

Disponibles seulement si le module est chargé :

iptables -j TEE
iptables -j TCPMSS
iptables -j HL


⸻

3️⃣ Ce que -j NE PEUT PAS être

❌ Une table :

-j nat

❌ Une target inexistante :

-j ALLOW

❌ Une chaîne d’une autre table

⸻

4️⃣ Comment voir les targets disponibles sur ton système

iptables -j help

ou

iptables -L -j


⸻

5️⃣ Résumé ultra-court

-j = CHAÎNE (même table)
-j = TARGET (action)

	•	❌ jamais une table
	•	❌ jamais une policy
	•	✅ toujours local à la table

⸻

Si tu veux, je peux :
•	te donner la liste exacte par table
•	expliquer quand utiliser RETURN
•	comparer chaque target avec son équivalent nftables

Dis-moi 👍