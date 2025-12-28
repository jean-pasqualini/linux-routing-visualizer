package state

type EventHandler func(name string, event any)

var handlers = make(map[string][]EventHandler)

func Dispatch(name string, event any) {
	for _, handler := range handlers[name] {
		handler(name, event)
	}
}
