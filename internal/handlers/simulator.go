package handlers

import (
	"context"
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
)

type SimulatorHandler struct {
}

func NewSimulatorHandler() *SimulatorHandler {
	return &SimulatorHandler{}
}

func (h *SimulatorHandler) Handle(context context.Context) {
	logger := logging.FromContext(context)
	logger.Debug("Handling request")
	fmt.Println("Hello World")
	ipt := iptable.NewBackend()

	state.Subscribe("logger", func(name string, event any) {
		if msg, ok := event.(string); ok {
			logger.Debug(msg)
		}
	})

	tables, err := ipt.Fetch()
	if err != nil {
		logger.Error("an error: " + err.Error())
		return
	}
	iptSim := simulator.NewSimulator(
		simulator.SimulatorQuery{
			Target: simulator.SimulatorQueryTarget{
				IP: "8.8.8.8",
			},
			Source: simulator.SimulatorQueryTarget{
				IP: "8.8.4.4",
			},
		},
		tables,
	)

	pp.Println(tables)
	pp.Println(iptSim.Simulate())
}
