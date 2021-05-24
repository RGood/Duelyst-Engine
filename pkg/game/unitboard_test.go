package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var board = NewUnitBoard(9, 5)
var foo = NewPlayer("Foo", "Lyonar", board, NewPosition(0, 2), true)
var bar = NewPlayer("Bar", "Songhai", board, NewPosition(8, 2), false)

func Test_positionEquality(t *testing.T) {
	assert.Equal(t, NewPosition(0, 0), NewPosition(0, 0))
}

func Test_doublePlacement(t *testing.T) {
	board := NewUnitBoard(9, 5)
	goblin := NewMinion("goblin", 1, 1)
	assert.True(t, board.PlaceUnit(goblin, NewPosition(0, 0)))

	gremlin := NewMinion("gremlin", 1, 1)
	assert.False(t, board.PlaceUnit(gremlin, NewPosition(0, 0)))
}

func Test_validMoves(t *testing.T) {
	p1, p2, gs := setupGamestate()
	p1general := p1.GetGeneral()
	p2general := p2.GetGeneral()

	assert.Equal(t, 8, len(p1general.GetValidMoves()))

	// When the generals are separate, they are not valid targets of each other
	assert.Equal(t, map[Unit]struct{}{}, p1.Board.GetValidTargets(p1general))

	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: NewPosition(2, 2),
	})

	assert.Equal(t, 12, len(p1.Board.GetValidMoves(p1general)))

	gs.MakeMove(&MoveAction{
		Unit:     p2general,
		Position: NewPosition(3, 2),
	})

	assert.Equal(t, 10, len(p1.Board.GetValidMoves(p1general)))

	// When the generals are next to each other, they are valid targets of each other
	assert.Equal(t, map[Unit]struct{}{p2general: struct{}{}}, p1.Board.GetValidTargets(p1general))
}
