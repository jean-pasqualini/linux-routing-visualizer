/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/handlers"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"

	"github.com/spf13/cobra"
)

// sysctlCmd represents the sysctl command
var sysctlCmd = &cobra.Command{
	Use:   "sysctl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		l := logging.New("linux-routing-visualizer")
		c := logging.WithLogger(context.Background(), l)
		h := handlers.NewSysctlHandler()
		h.Handle(c)
	},
}

func init() {
	rootCmd.AddCommand(sysctlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sysctlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sysctlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
