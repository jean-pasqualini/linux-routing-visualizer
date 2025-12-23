
raw → mangle → nat → filter → security

prerouting → input / forward → output → postrouting

priority -300  → exécuté AVANT
priority 0     → exécuté APRÈS


Paquet arrive
↓
Hook (ex: prerouting)
↓
Chaîne(s) attachée(s) à ce hook
↓
Règles exécutées