package gameengine

type Game interface {
	MakeMove(*Game) *Game
	HasEnded() bool
}
