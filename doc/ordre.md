C’est LA question fondamentale 👍
Et ta confusion est totalement normale.
La réponse clé est :

Le paquet ne “choisit” pas une table.
Le noyau le fait, selon une timeline fixe.

Je t’explique pas à pas, sans jargon inutile.

⸻

1️⃣ Le paquet ne décide rien

Un paquet :
•	n’est pas conscient des tables
•	ne “voit” pas iptables
•	ne sait pas où il va

👉 C’est le noyau Linux (Netfilter) qui applique un pipeline prédéfini.

⸻

2️⃣ Les tables = des phases du pipeline

Pense aux tables non pas comme des boîtes,
mais comme des étapes obligatoires dans le temps.

[ Arrivée paquet ]
↓
raw       (avant conntrack)
↓
mangle    (modifier)
↓
nat       (traduire)
↓
filter    (décider)
↓
[ Sortie paquet ]

➡️ Le paquet passe automatiquement par ces étapes.

⸻

3️⃣ Pourquoi plusieurs tables existent

Parce que chaque phase a un rôle unique :

Phase	Question posée
raw	“Je le track ou pas ?”
mangle	“Je modifie le paquet ?”
nat	“Je change les adresses ?”
filter	“Je l’autorise ?”
security	“LSM autorise ?”

👉 Mélanger tout ça serait ingérable.

⸻

4️⃣ Comment le noyau sait quelle table appliquer

Réponse courte

👉 Il ne choisit pas,
👉 il suit un ordre codé en dur dans le kernel.

⸻

Exemple : paquet entrant vers ta machine

1. raw PREROUTING
2. mangle PREROUTING
3. nat PREROUTING
4. décision de routage
5. mangle INPUT
6. filter INPUT
7. security INPUT

Aucune condition.
Aucune décision dynamique.
C’est automatique.

⸻

Exemple : paquet en transit

raw PREROUTING
mangle PREROUTING
nat PREROUTING
routing
mangle FORWARD
filter FORWARD
security FORWARD
mangle POSTROUTING
nat POSTROUTING


⸻

Exemple : paquet généré localement

raw OUTPUT
mangle OUTPUT
nat OUTPUT
filter OUTPUT
security OUTPUT
mangle POSTROUTING
nat POSTROUTING


⸻

5️⃣ Pourquoi on ne peut pas “sauter” de table

Parce que :
•	les tables sont le pipeline
•	ce serait comme “sauter le temps”

👉 Techniquement impossible et dangereux.

⸻

6️⃣ Analogie simple (très importante)

Imagine un contrôle aéroportuaire :
1.	Scan bagages (raw)
2.	Contrôle technique (mangle)
3.	Changement d’étiquette (nat)
4.	Contrôle passeport (filter)
5.	Sécurité finale (security)

Tu ne peux pas dire :

“Après le scanner, va directement au passeport”

C’est imposé par l’aéroport (le kernel).

⸻

7️⃣ Et nftables dans tout ça

nftables ne supprime pas la timeline.
Il la rend plus explicite :

chain prerouting {
type filter hook prerouting priority -300;
}

👉 Le hook = moment du pipeline
👉 Le priority = ordre précis

⸻

8️⃣ TL;DR (à retenir absolument)
•	❌ Le paquet ne choisit rien
•	✅ Le noyau applique un pipeline fixe
•	🧠 Les tables = phases temporelles
•	🔒 Ordre imposé par le kernel
•	🧩 nftables rend ça plus clair

⸻

Si tu veux, je peux :
•	te dessiner le pipeline ASCII ultra détaillé
•	te montrer un paquet réel tracer pas à pas
•	expliquer pourquoi NAT n’est pas du filtrage

Dis-moi 👌