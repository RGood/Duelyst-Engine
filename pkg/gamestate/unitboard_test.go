package gamestate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var foo = NewPlayer("Foo", "Lyonar", NewPosition(0, 2), true)
var bar = NewPlayer("Bar", "Songhai", NewPosition(8, 2), false)

func Test_positionEquality(t *testing.T) {
	assert.Equal(t, NewPosition(0, 0), NewPosition(0, 0))
}

func Test_doublePlacement(t *testing.T) {
	board := NewUnitBoard(9, 5)
	goblin := NewMinion("goblin", foo, 1, 1)
	assert.True(t, board.PlaceUnit(goblin, NewPosition(0, 0)))

	gremlin := NewMinion("gremlin", foo, 1, 1)
	assert.False(t, board.PlaceUnit(gremlin, NewPosition(0, 0)))
}

func Test_validMoves(t *testing.T) {
	p1, p2, gs := setupGamestate()
	p1units := gs.Board.GetPlayerUnits(p1)
	p1general := p1units[0]
	p2units := gs.Board.GetPlayerUnits(p2)
	p2general := p2units[0]

	assert.Equal(t, 8, len(gs.Board.GetValidMoves(p1general)))

	// When the generals are separate, they are not valid targets of each other
	assert.Equal(t, map[Unit]struct{}{}, gs.Board.GetValidTargets(p1general))

	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: NewPosition(2, 2),
	})

	assert.Equal(t, 12, len(gs.Board.GetValidMoves(p1general)))

	gs.MakeMove(&MoveAction{
		Unit:     p2general,
		Position: NewPosition(3, 2),
	})

	assert.Equal(t, 10, len(gs.Board.GetValidMoves(p1general)))

	// When the generals are next to each other, they are valid targets of each other
	assert.Equal(t, map[Unit]struct{}{p2general: struct{}{}}, gs.Board.GetValidTargets(p1general))
}
