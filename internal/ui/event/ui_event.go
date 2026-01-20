package event

type TabEventSubscribe interface {
	OnTabShow()
	OnTabHide()
}
