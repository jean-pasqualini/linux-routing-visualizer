package state

func Subscribe(name string, handler EventHandler) {
	if handlers[name] == nil {
		handlers[name] = make([]EventHandler, 0, 0)
	}
	handlers[name] = append(handlers[name], handler)
}
