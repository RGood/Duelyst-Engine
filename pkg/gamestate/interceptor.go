package gamestate

type Interceptor interface {
	Apply(Action, *Gamestate) Action
	Subscribe(*Gamestate)
	Unsubscribe(*Gamestate)
}
