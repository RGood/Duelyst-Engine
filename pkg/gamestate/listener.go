package gamestate

type Listener interface {
	Notify(Action, *Gamestate)
	Subscribe(*Gamestate)
	Unsubscribe(*Gamestate)
}
