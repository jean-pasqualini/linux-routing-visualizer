package arptable

import (
	"bytes"
	"errors"
	"os/exec"
)

type arptableBackend struct {
}

func NewArpTableBackend() *arptableBackend {
	return &arptableBackend{}
}

type ArpEntry struct {
}

func (b *arptableBackend) Fetch() ([]ArpEntry, error) {
	arpList := []ArpEntry{}

	return arpList, nil
}

func (b *arptableBackend) runProces() (string, error) {
	// -c add the counters , "-c"
	cmd := exec.Command("arptables-save")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	if errRun := cmd.Run(); errRun != nil {
		return "", errors.New(err.String())
	}

	return out.String(), nil
}
