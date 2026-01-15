Voici le **vrai “Hello World” d’un module noyau Linux** minimal, exactement ce qu’utilisent les devs kernel pour débuter.

---

## 1️⃣ `hello.c`

```c
#include <linux/module.h>
#include <linux/kernel.h>

static int __init hello_init(void)
{
    printk(KERN_INFO "Hello from kernel module!\n");
    return 0;
}

static void __exit hello_exit(void)
{
    printk(KERN_INFO "Goodbye from kernel module!\n");
}

module_init(hello_init);
module_exit(hello_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("toi");
MODULE_DESCRIPTION("Hello World kernel module");
```

---

## 2️⃣ `Makefile`

```makefile
obj-m += hello.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```

---

## 3️⃣ Compilation

```bash
make
```

Tu obtiens :

```
hello.ko
```

---

## 4️⃣ Chargement / déchargement

```bash
sudo insmod hello.ko
dmesg | tail

sudo rmmod hello
dmesg | tail
```

Tu verras :

```
Hello from kernel module!
Goodbye from kernel module!
```

---

## Ce que tu viens d’apprendre

| Ligne           | Rôle                             |
| --------------- | -------------------------------- |
| `module_init()` | fonction appelée au chargement   |
| `module_exit()` | fonction appelée au déchargement |
| `printk()`      | équivalent kernel de `printf`    |
| `.ko`           | module ELF chargé dans le noyau  |

---

Prochaine étape naturelle :
➡️ ajouter un **hook Netfilter** dans ce module pour intercepter des paquets.


--- Voici le **vrai “Hello World” d’un module noyau Linux** minimal, exactement ce qu’utilisent les devs kernel pour débuter.

---

## 1️⃣ `hello.c`

```c
#include <linux/module.h>
#include <linux/kernel.h>

static int __init hello_init(void)
{
    printk(KERN_INFO "Hello from kernel module!\n");
    return 0;
}

static void __exit hello_exit(void)
{
    printk(KERN_INFO "Goodbye from kernel module!\n");
}

module_init(hello_init);
module_exit(hello_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("toi");
MODULE_DESCRIPTION("Hello World kernel module");
```

---

## 2️⃣ `Makefile`

```makefile
obj-m += hello.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```

---

## 3️⃣ Compilation

```bash
make
```

Tu obtiens :

```
hello.ko
```

---

## 4️⃣ Chargement / déchargement

```bash
sudo insmod hello.ko
dmesg | tail

sudo rmmod hello
dmesg | tail
```

Tu verras :

```
Hello from kernel module!
Goodbye from kernel module!
```

---

## Ce que tu viens d’apprendre

| Ligne           | Rôle                             |
| --------------- | -------------------------------- |
| `module_init()` | fonction appelée au chargement   |
| `module_exit()` | fonction appelée au déchargement |
| `printk()`      | équivalent kernel de `printf`    |
| `.ko`           | module ELF chargé dans le noyau  |

---

Prochaine étape naturelle :
➡️ ajouter un **hook Netfilter** dans ce module pour intercepter des paquets.
