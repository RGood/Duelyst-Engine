package gamestate

type Action interface {
	Execute(*Gamestate) *Gamestate
}
