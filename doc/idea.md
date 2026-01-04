▶ PREROUTING
├─ rule: -p tcp --dport 21 ✔
│     ▶ chain: FTP
│        ├─ rule: -s 10.0.0.0/8 ✔
│        │     ▶ chain: FTP-INTERNAL
│        │        └─ rule: -j ACCEPT
