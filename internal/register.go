package internal

import (
	"context"
	"os"
	"time"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/ipvs"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/routing"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sniffing"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/socket"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/simulator"
	uiiptable "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/iptable"
	uiipvs "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/ipvs"
	uisimulator "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

func Register(app *tview.Application) {
	logger := logging.NewFilelogger()
	ctx := logging.WithLogger(context.Background(), logger)

	state.Subscribe("app:log", func(name string, event any) {
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
				State:            event.State,
				Protocol:         event.Protocol,
				IncludeUnmatched: event.ShowUnmatched,
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

	state.Subscribe("ipvs:request", func(name string, event any) {
		sBackend := ipvs.NewIPVSBackend()
		raw, services, err := sBackend.Fetch()
		if err == nil {
			state.Dispatch("ipvs:response", uiipvs.IPVSEvent{raw, services})
		}
	})

	state.Subscribe("routing:request", func(name string, event any) {
		backend := routing.NewRoutingBackend()
		config, err := backend.Fetch()
		if err == nil {
			state.Dispatch("routing:response", config)
		}
	})

	state.Subscribe("monitor:start", func(name string, event any) {
		go func() {
			backend := sniffing.NewSniffingBackend()
			packetChan, err := backend.Sniff(ctx, os.Getenv("NETDEVICE"), "", 1600, false, 2*time.Second)
			if err != nil {
				state.Dispatch("app:log", err.Error())
				return
			}

			state.Dispatch("app:log", "monitor:start")
			for packet := range packetChan {
				app.QueueUpdate(func() {
					state.Dispatch("app:log", "packet received")
					state.Dispatch("monitor:packet", packet)
				})
			}
		}()
	})
}
