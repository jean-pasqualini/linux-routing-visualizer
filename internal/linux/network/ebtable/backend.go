package ebtable

import (
	"bytes"
	"errors"
	"os/exec"
)

type EbtableBackend struct {
}

func NewEbtableBackend() *EbtableBackend {
	return &EbtableBackend{}
}

type Ebtable struct {
}

// ebtables -L
func (b *EbtableBackend) Fetch() (Ebtable, error) {
	var table Ebtable
	return table, nil
}

func (b *EbtableBackend) runProces() (string, error) {
	// -c add the counters , "-c"
	cmd := exec.Command("ebtables-save")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	if errRun := cmd.Run(); errRun != nil {
		return "", errors.New(err.String())
	}

	return out.String(), nil
}
