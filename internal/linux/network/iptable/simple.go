package iptable

import (
	"context"
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
)

type IPtableReader struct {
}

func NewIPtableReader() *IPtableReader {
	return &IPtableReader{}
}

func (r *IPtableReader) Read(context context.Context) {
	logger := logging.FromContext(context)
	ipt, _ := NewBackend()

	chains, err := ipt.ListChains("filter")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Println(chains)
}
