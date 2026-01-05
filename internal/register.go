package internal

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/socket"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/simulator"
	uiiptable "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/iptable"
	uisimulator "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
)

func Register() {
	logger := logging.NewFilelogger()

	state.Subscribe("logger", func(name string, event any) {
		if msg, ok := event.(string); ok {
			logger.Debug(msg)
		}
	})

	state.Subscribe("iptables:request", func(name string, event any) {
		ipt := iptable.NewBackend()
		tables, _ := ipt.ListChains("aeaze")
		raw := ipt.GetStdout()
		state.Dispatch("iptables:response", uiiptable.IpTableResponse{
			Parsed: tables,
			Raw:    raw,
		})
	})
	state.Subscribe("simulator_request", func(name string, event any) {
		if event, ok := event.(uisimulator.FormEvent); ok {
			ipt := iptable.NewBackend()
			tables, _ := ipt.ListChains("aeaze")
			sim := simulator.NewSimulator(simulator.SimulatorQuery{
				State:    event.State,
				Protocol: event.Protocol,
				Source: simulator.SimulatorQueryTarget{
					Device: event.Source.Device,
					IP:     event.Source.IP,
					Port:   event.Source.Port,
				},
				Target: simulator.SimulatorQueryTarget{
					Device: event.Target.Device,
					IP:     event.Target.IP,
					Port:   event.Target.Port,
				},
			}, tables)
			state.Dispatch("simulator_result", sim.Simulate())
		}
	})

	state.Subscribe("socket:request", func(name string, event any) {
		sBackend := socket.NewSocketBackend()
		state.Dispatch("socket:response", sBackend.ListListeners())
	})
}
