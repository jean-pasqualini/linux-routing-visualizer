/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/handlers"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var libpcapOpts handlers.LibpcapOptions

// libpcapDevCmd represents the libpcapDev command
var libpcapDevCmd = &cobra.Command{
	Use:   "libpcap-dev",
	Short: "A brief description of your command",

	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("libpcapDev called")
		l := logging.New("linux-routing-visualizer")
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

		h := handlers.NewLibpcapHandler(libpcapOpts)
		h.Handle(c)
	},
}

func init() {
	rootCmd.AddCommand(libpcapDevCmd)

	libpcapDevCmd.Flags().StringVarP(&libpcapOpts.Interface, "interface", "i", "",
		"Interface réseau (ex: eth0, wlan0)")
	libpcapDevCmd.Flags().StringVarP(&libpcapOpts.Filter, "filter", "f", "",
		"Filtre BPF (ex: tcp and port 443)")
	libpcapDevCmd.Flags().Int32VarP(&libpcapOpts.Snaplen, "snaplen", "s", 1600,
		"Taille max capturée par paquet")
	libpcapDevCmd.Flags().BoolVarP(&libpcapOpts.Promisc, "promisc", "p", false,
		"Mode promiscuous")
	libpcapDevCmd.Flags().DurationVarP(&libpcapOpts.Timeout, "timeout", "t", 2*time.Second,
		"Timeout de capture")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// libpcapDevCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// libpcapDevCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
