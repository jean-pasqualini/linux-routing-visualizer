package ns

import (
	"github.com/k0kubun/pp"
	"github.com/prometheus/procfs"
)

func ListNamespace() ([]string, error) {
	var list []string

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
		pp.Println(namespaces)
	}

	return list, nil
}
