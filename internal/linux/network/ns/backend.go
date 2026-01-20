package ns

import (
	"fmt"
	"github.com/prometheus/procfs"
	"path"
	"strings"
)

type NamespaceSpec struct {
	Path   string
	NodeID int
	Name   string
}

type NamespaceList map[int]NamespaceSpec

// /proc/pid/ns/net
// /var/run/docker/netns
// /run/containerd/netns
// /var/run/netns
// lsns -t net or ip netns ls
// mount | grep nsfs
// stat or readlink to get the inodeid
func UniqueSlice[T comparable](inputSlice []T) []T {
	uniqueSlice := make([]T, 0, len(inputSlice))
	seen := make(map[T]bool, len(inputSlice))
	for _, element := range inputSlice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}
	return uniqueSlice
}

func procIdentity(p procfs.Proc) string {
	name := ""
	groups, _ := p.Cgroups()
	cmdLine, _ := p.CmdLine()
	for _, group := range groups {
		if strings.HasSuffix(group.Path, ".service") {
			name += "[service]"

			return fmt.Sprintf("[service:%d] %s", p.PID, path.Base(group.Path))
		}
		if strings.Contains(group.Path, "docker-") {
			return fmt.Sprintf("[docker:%d] %s", p.PID, strings.Join(cmdLine, " "))
		}
	}
	return fmt.Sprintf("[process:%d] %s", p.PID, strings.Join(cmdLine, " "))
}

func ListNamespace() (NamespaceList, error) {
	list := make(NamespaceList)

	// Procs part
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return list, err
	}

	procs, err := fs.AllProcs()
	if err != nil {
		return list, err
	}

	for _, proc := range procs {
		namespaces, err := proc.Namespaces()
		if err != nil {
			continue
		}
		id := int(namespaces["net"].Inode)
		if _, ok := list[id]; !ok {
			list[id] = NamespaceSpec{
				Path:   fmt.Sprintf("/proc/%d/ns/net", proc.PID),
				NodeID: id,
				Name:   procIdentity(proc),
			}
		}
	}

	return list, nil
}
