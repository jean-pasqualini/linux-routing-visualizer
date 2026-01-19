runtime.LockOSThread()
defer runtime.UnlockOSThread()

origNS, _ := os.Open("/proc/self/ns/net")
targetNS, _ := os.Open("/run/docker/netns/2161bf945a7d")

unix.Setns(int(targetNS.Fd()), unix.CLONE_NEWNET)

// opérations réseau ici

unix.Setns(int(origNS.Fd()), unix.CLONE_NEWNET)

---

Oui — **et c’est un point très important**, mais il faut être **précis** sur *quoi* change et *quand*.

👉 **Ce n’est PAS le PID du process qui change**,
👉 **c’est le PID *visible*** quand tu changes de **PID namespace**.

Je détaille clairement 👇

---

## 1️⃣ Ce qui NE change PAS

Quand tu fais :

```c
setns(pidns_fd, CLONE_NEWPID);
```

* ❌ le **PID réel du process** (dans le kernel) **ne change pas**
* ❌ le **PID du thread courant** ne change pas
* ❌ ton process ne “renaît” pas

👉 le kernel garde le même `task_struct`.

---

## 2️⃣ Ce qui change (et seulement après `fork()`)

Le **PID namespace est spécial**.

### Règle clé :

> **Entrer dans un PID namespace n’a d’effet visible qu’après un `fork()`**

Donc :

```text
setns(pidns)
fork()
```

* le **process fils**

    * a un **nouveau PID** dans ce namespace
    * souvent `1`
* le parent garde son PID original

---

## 3️⃣ Comparaison avec network namespace

| Namespace      | PID change ? | Effet immédiat |
| -------------- | ------------ | -------------- |
| net            | ❌            | ✅ oui          |
| mount          | ❌            | ✅ oui          |
| uts            | ❌            | ✅ oui          |
| ipc            | ❌            | ✅ oui          |
| **pid**        | ❌ (parent)   | ❌              |
| **pid (fork)** | ✅ (child)    | ✅              |

---

## 4️⃣ Exemple concret (shell)

```bash
# entrer dans un pidns existant
nsenter --pid=/proc/123/ns/pid -- bash
echo $$        # PID différent
```

👉 `nsenter` fait :

1. `setns(pidns)`
2. `fork()`
3. `exec()`

---

## 5️⃣ Pourquoi ça fonctionne sans fork pour le netns

Le **network namespace** :

* est consulté **à chaque syscall réseau**
* donc le changement est immédiat

Pas besoin de fork.

---

## 6️⃣ En Go (important)

👉 Tu **NE peux PAS** correctement utiliser un PID namespace sans fork.

En Go :

* pas de fork pur
* `exec.Cmd` fait un fork + exec
* donc PID namespace = process dédié

---

## 🔑 Résumé ultra-clair

* ❌ `setns(netns)` → PID inchangé
* ❌ `setns(pidns)` → PID inchangé
* ✅ `setns(pidns)` + `fork()` → **nouveau PID visible**
* 🧠 seul le PID namespace fonctionne comme ça

---

## TL;DR

> **Tu avais raison de poser la question**
>
> Le PID ne change **que dans un PID namespace, et seulement après fork**.

Si tu veux, je peux :

* te montrer le code kernel exact
* expliquer pourquoi le PID namespace est unique
* comparer avec user namespace

Dis-moi 👍
