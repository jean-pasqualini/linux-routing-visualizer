/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "A tui test",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pp.SetColorScheme(pp.ColorScheme{
			Bool:            pp.NoColor,
			Integer:         pp.NoColor,
			Float:           pp.NoColor,
			String:          pp.NoColor,
			StringQuotation: pp.NoColor,
			EscapedChar:     pp.NoColor,
			FieldName:       pp.NoColor,
			PointerAdress:   pp.NoColor,
			Nil:             pp.NoColor,
			Time:            pp.NoColor,
			StructName:      pp.NoColor,
			ObjectLength:    pp.NoColor,
		})

		app := tview.NewApplication()

		tabPanel := ui.NewSidePanel()
		mainPanel := ui.NewMainPanel()

		ipt := iptable.NewBackend()
		tables, _ := ipt.ListChains("aeaze")
		raw := ipt.GetStdout()
		go func() {
			mainPanel.ShowTables(app, tables, raw)
		}()

		layout := tview.NewFlex().
			AddItem(tabPanel, 50, 0, true).
			AddItem(mainPanel, 0, 1, false)

		frame := tview.NewFrame(layout).
			SetBorders(0, 0, 0, 0, 0, 0).
			AddText("Routing Visualizer", true, tview.AlignCenter, tcell.ColorWhite)

		app.SetRoot(frame, true).EnableMouse(true)
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tuiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tuiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
