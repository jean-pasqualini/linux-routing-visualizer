══════════════════════════════════════════════════════════════════════
FIREWALL SIMULATION REPORT
══════════════════════════════════════════════════════════════════════

┌───────────── Packet Simulation ─────────────┐
│ SRC      : 192.168.1.10                     │
│ DST      : 10.0.0.5                         │
│ DPORT    : 443                              │
│ PROTOCOL : TCP                              │
└─────────────────────────────────────────────┘

Evaluation path:

1. [FILTER] INPUT chain
   ├─ Rule #3:  ACCEPT tcp  --  192.168.1.0/24 → ANY dport 443
   │   └─ MATCH ✓  (source subnet match)
   │   └─ MATCH ✓  (destination port = 443)
   │   └─ ACTION → ACCEPT
   │
   Final verdict:  ACCEPT



══════════════════════════════════════════════════════════════════════
1) TIMELINE
   ══════════════════════════════════════════════════════════════════════

INPUT chain:
Rule 1 → no match (protocol mismatch)
Rule 2 → no match (dport 22 ≠ 443)
Rule 3 → MATCH ✓  → ACCEPT


══════════════════════════════════════════════════════════════════════
2) DECISION TREE
   ══════════════════════════════════════════════════════════════════════

protocol == TCP ?
├─ no  → next rule
└─ yes
dport == 443 ?
├─ no  → next rule
└─ yes → ACCEPT


══════════════════════════════════════════════════════════════════════
3) DIFF VIEW (Why rules failed)
   ══════════════════════════════════════════════════════════════════════

Rule 7: DROP tcp 192.168.0.0/16 → 10.0.0.0/8 dport 22

Expected        | Received
----------------|----------------
DPORT: 22       | 443        ✗
SRC: 192.168/16 | 192.168.1.10 ✓


══════════════════════════════════════════════════════════════════════
4) FORENSIC REPORT
   ══════════════════════════════════════════════════════════════════════

VERDICT : DROP
REASON  : No rule allowing TCP traffic to 10.0.0.5:443

Closest matching rule:
Rule 5 ACCEPT tcp 192.168.1.0/24 → 10.0.0.0/8 dport 80
Missing condition: dport 443


══════════════════════════════════════════════════════════════════════
5) NATURAL LANGUAGE
   ══════════════════════════════════════════════════════════════════════

The TCP packet coming from 192.168.1.10 to 10.0.0.5:443
was blocked because no firewall rule authorizes
port 443 on the INPUT chain.


══════════════════════════════════════════════════════════════════════
6) FIREWALL LENS (Rule Perspective)
   ══════════════════════════════════════════════════════════════════════

Rule #3 ACCEPT tcp 192.168.1.0/24 → ANY dport 443

SRC      : 192.168.1.10   ∈ 192.168.1.0/24  ✔
DST      : 10.0.0.5       ∈ ANY             ✔
DPORT    : 443            == 443            ✔
PROTO    : TCP            == TCP            ✔

→ FINAL ACTION: ACCEPT

══════════════════════════════════════════════════════════════════════
