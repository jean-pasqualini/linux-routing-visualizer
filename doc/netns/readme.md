| Emplacement                    | Qui          |
| ------------------------------ | ------------ |
| `/proc/<pid>/ns/net`           | tout process |
| `/var/run/netns/*`             | `ip netns`   |
| `/run/systemd/nspawn/*/ns/net` | systemd      |
| `/tmp/foo`                     | toi          |
| `/myapp/namespaces/router1`    | toi          |
| FD anonyme                     | process Go   |
