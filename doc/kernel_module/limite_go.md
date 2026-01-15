Oui — et là on est dans le **vrai terrain de jeu moderne** 😄
👉 écrire du **eBPF en Rust** est **déjà utilisé en production** (Cilium, Aya, Falco…).

Tu écris :

* eBPF **dans le kernel** en Rust `no_std`
* le control-plane en Rust userspace

---

## Deux grandes stacks Rust eBPF

| Stack   | Usage                |
| ------- | -------------------- |
| **Aya** | 100% Rust, zéro LLVM |
| RedBPF  | ancien, abandonné    |

👉 **Aya est le standard actuel.**

---

## Exemple complet : DNAT-like tc ingress en Rust

### 1️⃣ Crée le projet

```bash
cargo generate --git https://github.com/aya-rs/aya-template
```

Choisis `tc`.

---

## 2️⃣ Programme eBPF — `dnat-ebpf/src/main.rs`

```rust
#![no_std]
#![no_main]

use aya_bpf::{
    macros::tc,
    programs::TcContext,
};
use aya_bpf::bindings::{ETH_P_IP, IPPROTO_TCP};
use core::mem;
use network_types::{
    eth::EthHdr,
    ip::Ipv4Hdr,
    tcp::TcpHdr,
};

#[tc(name="dnat_ingress")]
pub fn dnat_ingress(ctx: TcContext) -> i32 {
    match try_dnat(ctx) {
        Ok(ret) => ret,
        Err(_) => 0,
    }
}

fn try_dnat(ctx: TcContext) -> Result<i32, ()> {
    let eth: EthHdr = ctx.load(0).map_err(|_| ())?;
    if eth.ether_type != u16::to_be(ETH_P_IP as u16) {
        return Ok(0);
    }

    let ip: Ipv4Hdr = ctx.load(mem::size_of::<EthHdr>()).map_err(|_| ())?;
    if ip.proto != IPPROTO_TCP as u8 {
        return Ok(0);
    }

    let tcp_off = mem::size_of::<EthHdr>() + (ip.ihl() as usize * 4);
    let mut tcp: TcpHdr = ctx.load(tcp_off).map_err(|_| ())?;

    if tcp.dest != u16::to_be(80) {
        return Ok(0);
    }

    // rewrite dest port
    ctx.store(tcp_off + 2, &u16::to_be(8080), 0).map_err(|_| ())?;

    Ok(0)
}
```

---

## 3️⃣ Loader userspace — `dnat/src/main.rs`

```rust
use aya::{Bpf, programs::Tc, util::online_cpus};
use std::convert::TryInto;

fn main() -> Result<(), anyhow::Error> {
    let mut bpf = Bpf::load_file("dnat-ebpf/target/bpfel-unknown-none/release/dnat-ebpf")?;
    let program: &mut Tc = bpf.program_mut("dnat_ingress")?.try_into()?;

    program.load()?;
    program.attach("eth0", aya::programs::TcAttachType::Ingress)?;

    loop { std::thread::park(); }
}
```

---

## 4️⃣ Build

```bash
cargo xtask build-ebpf
cargo run
```

---

## Ce que tu obtiens

* eBPF kernel en **Rust memory-safe**
* chargement, attachement tc via Rust
* zéro C

---

## Limites

| Élément               | Rust eBPF        |
| --------------------- | ---------------- |
| Conntrack NAT complet | ❌ encore à coder |
| Fragmentation IP      | ❌                |
| Helpers avancés NAT   | ⚠️               |

---

## Conclusion

> Tu ne peux pas faire du DNAT kernel en Go.
> Tu peux le faire **en Rust + eBPF** — c’est la voie moderne, safe, et utilisée en prod.


----

Oui — et tu mets le doigt sur **la vraie exception** à la règle.

👉 **Netlink est la “porte dérobée officielle” vers l’intérieur du kernel.**

Ce n’est pas un syscall isolé, c’est un **bus RPC bidirectionnel** entre userspace et des sous-systèmes kernel.

---

## Ce que Netlink te permet vraiment de faire

| Sous-système kernel      | Accessible en Go via Netlink        |
| ------------------------ | ----------------------------------- |
| nftables / iptables-nft  | ✅ créer règles DNAT, SNAT, firewall |
| Routing table (ip route) | ✅                                   |
| Interfaces réseau        | ✅ créer / supprimer / config        |
| Neighbors (ARP, ND)      | ✅                                   |
| Conntrack                | ✅ lister, tuer des flows            |
| XFRM (IPsec)             | ✅                                   |
| tc (qdisc, filters)      | ⚠️ partiel                          |
| WireGuard                | ✅                                   |
| cgroups v2               | ⚠️                                  |
| eBPF                     | ⚠️ charger / maps                   |

Netlink est ce qui permet à `ip`, `tc`, `nft`, `ss`, `conntrack` de fonctionner.

---

## Pourquoi Netlink est si puissant

Quand tu fais en Go :

```go
c := &nftables.Conn{}
c.AddRule(...)
c.Flush()
```

Tu ne fais pas du NAT toi-même.

Tu dis au kernel :

> “Voici un programme NAT. Exécute-le dans ton moteur interne.”

Donc tu **programmes le kernel**, sans y injecter ton propre code.

---

## Limite fondamentale de Netlink

| Action                               | Possible |
| ------------------------------------ | -------- |
| Ajouter une règle NAT                | ✅        |
| Supprimer un flow conntrack          | ✅        |
| Réécrire un paquet arbitraire        | ❌        |
| Implémenter un nouvel algorithme NAT | ❌        |
| Accéder à `skb->data`                | ❌        |

Netlink te donne accès aux **API stables** du kernel — pas aux structures internes.

---

## Conclusion

> Go + Netlink te permet de **configurer le comportement du kernel**.
> eBPF / modules te permettent de **changer ce comportement**.

Netlink = *control plane*
eBPF / modules = *data plane*
