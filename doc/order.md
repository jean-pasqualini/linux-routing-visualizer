### Paquet entrant 
```
raw → mangle → nat (PREROUTING)
        ↓
   routage
        ↓
   mangle → filter → security (INPUT)
```

### Paquet transité (FORWARD)
```
raw → mangle → nat (PREROUTING)
        ↓
   routage
        ↓
   mangle → filter → security (FORWARD)
        ↓
   mangle → nat (POSTROUTING)
```

### Paquet sortant
```
raw → mangle → nat → filter → security (OUTPUT)
        ↓
   décision de routage
        ↓
   mangle → nat (POSTROUTING)
```