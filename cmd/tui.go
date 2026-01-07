/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/log"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var debugMode bool

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "A tui test",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//tview.Borders.Horizontal = '⎺'
		//tview.Borders.Vertical = '▌'

		app := tview.NewApplication()

		internal.Register(app)

		tabPanel := ui.NewSidePanel()
		mainPanel := ui.NewMainPanel()

		layout := tview.NewFlex()
		if debugMode {
			logPanel := log.NewLogPanel()
			layout.AddItem(logPanel, 20, 0, false)
		}

		layout.AddItem(tabPanel, 50, 0, true)
		layout.AddItem(mainPanel, 0, 1, false)

		frame := tview.NewFrame(layout).
			SetBorders(0, 0, 0, 0, 0, 0).
			AddText("Routing Visualizer", true, tview.AlignCenter, tcell.ColorWhite)

		app.SetRoot(frame, true).EnableMouse(true)
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)

	tuiCmd.Flags().BoolVar(&debugMode, "debug", false, "show the log pane")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tuiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tuiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
