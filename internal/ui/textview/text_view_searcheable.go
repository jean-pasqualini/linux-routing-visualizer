package textview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"regexp"
)

type TextViewSearchable struct {
	*tview.Pages
	textView      *tview.TextView
	inputText     *tview.InputField
	visibleSearch bool
}

func NewTextViewSearchable() *TextViewSearchable {
	pages := tview.NewPages()
	textView := tview.NewTextView().
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true)
	inputText := tview.NewInputField()
	inputText.SetBorder(true)
	inputText.SetTitle("Search")
	inputText.SetFieldBackgroundColor(tcell.ColorBlack)
	pages.AddPage("textview", textView, true, true)
	pages.AddPage("inputText", inputText, false, true)
	return &TextViewSearchable{
		Pages:     pages,
		textView:  textView,
		inputText: inputText,
	}
}

func (t *TextViewSearchable) Draw(screen tcell.Screen) {
	t.Box.DrawForSubclass(screen, t)
	x, y, width, height := t.Pages.GetInnerRect()
	textView := t.Pages.GetPage("textview")
	textView.SetRect(x, y, width, height)
	textView.Draw(screen)

	if t.visibleSearch {
		inputText := t.Pages.GetPage("inputText")
		if inputText, ok := inputText.(*tview.InputField); ok {
			inputText.SetRect(x+width-30, y+height-3, 30, 3)
			inputText.Draw(screen)
		}
	}
}

func TagSelectionRegions(text, query string) string {
	if query == "" {
		return text
	}

	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(query))

	return re.ReplaceAllStringFunc(text, func(m string) string {
		return `["selection"]` + m + `[""]`
	})
}

// InputHandler returns the handler for this primitive.
func (t *TextViewSearchable) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Rune() == '/' {
			t.visibleSearch = true
			t.textView.Focus(setFocus)
			return
		}
		if event.Key() == tcell.KeyEnter {
			t.visibleSearch = false
			t.textView.SetText(TagSelectionRegions(t.textView.GetText(true), t.inputText.GetText()))
			t.textView.Highlight("selection")
			t.textView.ScrollToHighlight()
			return
		}
		if t.visibleSearch == false {
			if event.Rune() == 'n' {
				t.textView.ScrollToHighlight()
			}

			return
		}
		if event.Key() == tcell.KeyEscape {
			t.visibleSearch = false
			return
		}
		if handler := t.Pages.InputHandler(); handler != nil {
			handler(event, setFocus)
			return
		}
	}
}

func (t *TextViewSearchable) SetText(text string) {
	t.textView.SetText(text)
}
