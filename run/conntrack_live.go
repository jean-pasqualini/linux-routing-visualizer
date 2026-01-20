package main

import (
	"context"
	"fmt"
	"time"

	ct "github.com/florianl/go-conntrack"
)

func main() {
	nfct, err := ct.Open(&ct.Config{})
	if err != nil {
		fmt.Println("could not create nfct:", err)
		return
	}
	defer nfct.Close()

	monitor := func(c ct.Con) int {
		fmt.Printf("%#v\n", c)
		return 0
	}

	if err := nfct.Register(context.Background(), ct.Conntrack, ct.NetlinkCtNew, monitor); err != nil {
		fmt.Println("could not register callback:", err)
		return
	}

	fmt.Println("Waiting for events...")
	time.Sleep(60 * time.Second)

}
