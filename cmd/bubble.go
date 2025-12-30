/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"
)

var roundFullBorder = lipgloss.Border{
	Top:         "▀",
	Bottom:      "▄",
	Left:        "▌",
	Right:       "▐",
	TopLeft:     "▛",
	TopRight:    "▜",
	BottomLeft:  "▙",
	BottomRight: "▟",
}

var btnStyle = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("15")).
	Background(lipgloss.Color("2")).
	Foreground(lipgloss.Color("0")).
	Padding(0, 3).
	Bold(true)

type model struct {
	t table.Model
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok && k.String() == "ctrl+c" {
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	//button := btnStyle.Render("VALIDER")
	return baseStyle.Render(m.t.View()) + "\n"
	//return lipgloss.Place(40, 10, lipgloss.Center, lipgloss.Center, button)
}

// bubbleCmd represents the bubble command
var bubbleCmd = &cobra.Command{
	Use:   "bubble",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		columns := []table.Column{
			{Title: "Rank", Width: 4},
			{Title: "City", Width: 10},
			{Title: "Country", Width: 10},
			{Title: "Population", Width: 10},
		}

		rows := []table.Row{
			{"1", "Tokyo", "Japan", "37,274,000"},
			{"2", "Delhi", "India", "32,065,760"},
			{"3", "Shanghai", "China", "28,516,904"},
			{"4", "Dhaka", "Bangladesh", "22,478,116"},
			{"5", "São Paulo", "Brazil", "22,429,800"},
			{"6", "Mexico City", "Mexico", "22,085,140"},
			{"7", "Cairo", "Egypt", "21,750,020"},
			{"8", "Beijing", "China", "21,333,332"},
			{"9", "Mumbai", "India", "20,961,472"},
			{"10", "Osaka", "Japan", "19,059,856"},
			{"11", "Chongqing", "China", "16,874,740"},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(7),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		p := tea.NewProgram(model{t})
		if err := p.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(bubbleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bubbleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bubbleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
