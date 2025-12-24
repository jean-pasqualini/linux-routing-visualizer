package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DiagramCanvas struct {
	*tview.Box
	width, height int
	buffer        [][]rune
	nodes         []*Node
	edges         []*Edge
}

type Node struct {
	X, Y  int
	W, H  int
	Title string
}

type Edge struct {
	From, To *Node
}

func NewDiagramCanvas(w, h int) *DiagramCanvas {
	d := &DiagramCanvas{
		Box:    tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Diagram"),
		width:  w,
		height: h,
		nodes:  []*Node{},
		edges:  []*Edge{},
	}

	d.buffer = make([][]rune, h)
	for y := 0; y < h; y++ {
		d.buffer[y] = make([]rune, w)
		for x := 0; x < w; x++ {
			d.buffer[y][x] = ' '
		}
	}
	return d
}

func (d *DiagramCanvas) AddNode(n *Node) {
	d.nodes = append(d.nodes, n)
}

func (d *DiagramCanvas) AddEdge(a, b *Node) {
	d.edges = append(d.edges, &Edge{From: a, To: b})
}

func (d *DiagramCanvas) clear() {
	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x++ {
			d.buffer[y][x] = ' '
		}
	}
}
func (d *DiagramCanvas) Draw(screen tcell.Screen) {
	d.Box.DrawForSubclass(screen, d)

	x0, y0, w, h := d.GetInnerRect()

	d.clear()

	for _, n := range d.nodes {
		d.drawNode(n)
	}
	for _, e := range d.edges {
		d.drawEdge(e.From, e.To)
	}

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	for y := 0; y < h && y < d.height; y++ {
		for x := 0; x < w && x < d.width; x++ {
			screen.SetContent(x0+x, y0+y, d.buffer[y][x], nil, style)
		}
	}
}

func (d *DiagramCanvas) drawNode(n *Node) {
	x, y, w, h := n.X, n.Y, n.W, n.H

	for i := 0; i < w; i++ {
		d.buffer[y][x+i] = '-'
		d.buffer[y+h-1][x+i] = '-'
	}
	for i := 0; i < h; i++ {
		d.buffer[y+i][x] = '|'
		d.buffer[y+i][x+w-1] = '|'
	}

	d.buffer[y][x] = '+'
	d.buffer[y][x+w-1] = '+'
	d.buffer[y+h-1][x] = '+'
	d.buffer[y+h-1][x+w-1] = '+'

	title := n.Title
	for i, r := range title {
		if i+1 < w-1 {
			d.buffer[y+h/2][x+1+i] = r
		}
	}
}

func (d *DiagramCanvas) drawEdge(a, b *Node) {
	x1 := a.X + a.W
	y1 := a.Y + a.H/2
	x2 := b.X
	_ = b.Y + b.H/2

	for x := x1; x < x2; x++ {
		d.buffer[y1][x] = '-'
	}
	d.buffer[y1][x2-1] = '>'
}
