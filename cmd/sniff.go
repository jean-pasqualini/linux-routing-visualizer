/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/handlers"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"

	"github.com/spf13/cobra"
)

var sniffOpts handlers.LibpcapOptions

// sniffCmd represents the libpcapDev command
var sniffCmd = &cobra.Command{
	Use: "sniff",
	Run: func(cmd *cobra.Command, args []string) {
		l := logging.NewFilelogger()
		c := logging.WithLogger(context.Background(), l)

		// contexte annulable (CTRL+C)
		c, cancel := context.WithCancel(c)
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sig
			cancel()
		}()

		h := handlers.NewLibpcapHandler(sniffOpts)
		h.Handle(c)
	},
}

func init() {
	rootCmd.AddCommand(sniffCmd)

	sniffCmd.Flags().StringVarP(&sniffOpts.Interface, "interface", "i", "",
		"Interface réseau (ex: eth0, wlan0)")
	sniffCmd.Flags().StringVarP(&sniffOpts.Filter, "filter", "f", "",
		"Filtre BPF (ex: tcp and port 443)")
	sniffCmd.Flags().Int32VarP(&sniffOpts.Snaplen, "snaplen", "s", 1600,
		"Taille max capturée par paquet")
	sniffCmd.Flags().BoolVarP(&sniffOpts.Promisc, "promisc", "p", false,
		"Mode promiscuous")
	sniffCmd.Flags().DurationVarP(&sniffOpts.Timeout, "timeout", "t", 20*time.Second,
		"Timeout de capture")
}
