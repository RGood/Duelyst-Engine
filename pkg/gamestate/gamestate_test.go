package gamestate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_draw(t *testing.T) {
	p1, p2, gs := setupGamestate()

	p1units := gs.Board.GetPlayerUnits(p1)
	p1general := p1units[0]
	p2units := gs.Board.GetPlayerUnits(p2)
	p2general := p2units[0]

	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: NewPosition(7, 2),
	})

	// Attack until death
	for i := 0; i < 50; i++ {
		gs.MakeMove(&AttackAction{
			Attacker: p1general,
			Defender: p2general,
		})

		gs.MakeMove(&EndTurnAction{
			Owner: gs.ActivePlayer,
		})
	}

	assert.True(t, gs.HasEnded())
	hasEnded, winner := gs.Winner()

	assert.True(t, hasEnded)
	assert.Nil(t, winner)
}

func Test_winner(t *testing.T) {
	p1, p2, gs := setupGamestate()

	p1units := gs.Board.GetPlayerUnits(p1)
	p1general := p1units[0]
	p2units := gs.Board.GetPlayerUnits(p2)
	p2general := p2units[0]

	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: NewPosition(7, 2),
	})

	pf := NewDamageSpell("Phoenix Fire", 2, 3, func(owner *Player, game *Gamestate, damage int, targets []Unit, _ []Position) {
		if len(targets) == 1 {
			game.QueueAction(&DamageAction{
				Unit:   targets[0],
				Damage: damage,
			})
		}
	})

	pf.Cast(p1, gs, []Unit{p2general}, nil)

	// Attack until death
	for i := 0; i < 50; i++ {
		gs.MakeMove(&AttackAction{
			Attacker: p1general,
			Defender: p2general,
		})

		gs.MakeMove(&EndTurnAction{
			Owner: gs.ActivePlayer,
		})
	}

	assert.True(t, gs.HasEnded())
	hasEnded, winner := gs.Winner()

	assert.True(t, hasEnded)
	assert.Equal(t, p1, winner)
}
