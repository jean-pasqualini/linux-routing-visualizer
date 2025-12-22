/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/handlers"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "linux-routing-visualizer",
	Short: "A linux routing visualizer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		l := logging.New("linux-routing-visualizer")
		c := logging.WithLogger(context.Background(), l)
		h := handlers.NewAppHandler()
		h.Handle(c)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
