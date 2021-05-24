package gamestate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPlayer struct {
	Alive bool
}

func (player *TestPlayer) IsAlive() bool {
	return player.Alive
}

func NewTestPlayer(alive bool) *TestPlayer {
	return &TestPlayer{
		Alive: alive,
	}
}

func Test_gamestate(t *testing.T) {
	p1 := NewTestPlayer(true)
	p2 := NewTestPlayer(true)

	gamestate := NewGamestate(p1, p2)

	assert.False(t, gamestate.HasEnded())

	hasEnded, winner := gamestate.Winner()
	assert.False(t, hasEnded)
	assert.Nil(t, winner)

	assert.Equal(t, p1, gamestate.ActivePlayer)

	gamestate.EndTurn()

	assert.Equal(t, p2, gamestate.ActivePlayer)

	p2.Alive = false

	assert.True(t, gamestate.HasEnded())

	hasEnded, winner = gamestate.Winner()
	assert.True(t, hasEnded)
	assert.Equal(t, p1, winner)

	gamestate.EndTurn()
	assert.Equal(t, p2, gamestate.ActivePlayer)
}
