/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "A tui test",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app := tview.NewApplication()

		buildPage := func(name string) *ui.DiagramCanvas {
			canvas := ui.NewDiagramCanvas(80, 24)
			r1 := &ui.Node{X: 10, Y: 3, W: 10, H: 5, Title: "Raw " + name}
			r2 := &ui.Node{X: 20 + 5, Y: 3, W: 10, H: 5, Title: "Mangle"}
			r3 := &ui.Node{X: 30 + 5 + 5, Y: 3, W: 10, H: 5, Title: "Nat"}
			r4 := &ui.Node{X: 40 + 5 + 5 + 5, Y: 3, W: 10, H: 5, Title: "Filter"}
			r5 := &ui.Node{X: 50 + 5 + 5 + 5 + 5, Y: 3, W: 10, H: 5, Title: "Security"}
			canvas.AddNode(r1)
			canvas.AddNode(r2)
			canvas.AddNode(r3)
			canvas.AddNode(r4)
			canvas.AddNode(r5)
			canvas.AddEdge(r1, r2)
			canvas.AddEdge(r2, r3)
			canvas.AddEdge(r3, r4)
			canvas.AddEdge(r4, r5)

			return canvas
		}

		pages := tview.NewPages()
		pages.AddPage("lol", buildPage("lol"), true, true)
		pages.AddPage("lal", buildPage("lal"), true, true)

		tabbar := ui.NewTabBar([]string{"lol", "lal"}).
			SetOnSelect(func(_ int, name string) {
				pages.SwitchToPage(name)
			})

		tabContainer := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tabbar, 1, 0, false).
			AddItem(pages, 0, 1, false)

		layout := tview.NewFlex().
			AddItem(tabContainer, 0, 1, true).
			AddItem(tview.NewBox().SetBorder(true).SetBackgroundColor(tcell.ColorBlue), 0, 1, false)

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
