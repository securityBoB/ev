package event

type Listener struct {
	Handler Handler
}

func NewListener(handler Handler) *Listener {
	return &Listener{Handler: handler}
}
