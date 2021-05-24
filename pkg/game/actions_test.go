package game

import (
	"testing"

	"github.com/RGood/game_engine/pkg/gamestate"
	"github.com/stretchr/testify/assert"
)

func setupGamestate() (*Player, *Player, *gamestate.Gamestate) {
	board := NewUnitBoard(9, 5)
	p1 := NewPlayer("Foo", "Lyonar", board, NewPosition(0, 2), true)
	p2 := NewPlayer("Bar", "Songhai", board, NewPosition(8, 2), false)

	return p1, p2, gamestate.NewGamestate(p1, p2)
}

func Test_gameSetup(t *testing.T) {
	p1, p2, _ := setupGamestate()
	p1units := p1.GetUnits()

	assert.Equal(t, 1, len(p1units))
	assert.Equal(t, "general", p1units[0].GetType())
	assert.Equal(t, 2, p1units[0].GetAttack())
	assert.Equal(t, 25, p1units[0].GetHp())

	p2units := p2.Board.GetPlayerUnits(p2)
	assert.Equal(t, 1, len(p2units))
	assert.Equal(t, "general", p2units[0].GetType())
	assert.Equal(t, 2, p2units[0].GetAttack())
	assert.Equal(t, 25, p2units[0].GetHp())

	assert.True(t, p1.IsAlive())
	assert.True(t, p2.IsAlive())
}

func Test_move(t *testing.T) {
	p1, _, gs := setupGamestate()
	units := p1.Board.GetPlayerUnits(p1)

	assert.Equal(t, 1, len(units))
	unit := units[0]

	assert.Equal(t, "general", unit.GetType())
	assert.Equal(t, NewPosition(0, 2), unit.GetPosition())

	// This wasn't always true, and it was horrible to debug
	assert.Equal(t, unit.GetPosition(), p1.Board.GetPosition(unit))

	gs.MakeMove(&MoveAction{
		Unit:     unit,
		Position: NewPosition(2, 2),
	})

	assert.Equal(t, NewPosition(2, 2), unit.GetPosition())
}

func Test_health(t *testing.T) {
	p1, _, gs := setupGamestate()
	units := p1.Board.GetPlayerUnits(p1)
	general := units[0]

	assert.Equal(t, 25, general.GetHp())

	gs.MakeMove(&DamageAction{
		Unit:   general,
		Damage: 2,
	})

	assert.Equal(t, 23, general.GetHp())

	// Test overheal
	gs.MakeMove(&HealAction{
		Unit: general,
		Heal: 5,
	})

	assert.Equal(t, 25, general.GetHp())
}

func Test_attack(t *testing.T) {
	p1, p2, gs := setupGamestate()

	p1units := p1.Board.GetPlayerUnits(p1)
	p1general := p1units[0]

	p2units := p2.Board.GetPlayerUnits(p2)
	p2general := p2units[0]

	infrontOfP2 := p2general.GetPosition().Diff(NewPosition(1, 0))

	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: infrontOfP2,
	})

	gs.MakeMove(&AttackAction{
		Attacker: p1general,
		Defender: p2general,
	})

	assert.Equal(t, 23, p1general.GetHp())
	assert.Equal(t, 23, p2general.GetHp())
}

func Test_blastAttack(t *testing.T) {

	p1, p2, gs := setupGamestate()

	p1units := p1.GetUnits()
	p1general := p1units[0]

	p1general.AddAttribute("blast", 0)

	gremlin1 := NewMinion("gremlin", 1, 1)
	gs.MakeMove(&PlaceUnitAction{Owner: p2, Unit: gremlin1, Position: NewPosition(7, 2)})
	gremlin2 := NewMinion("gremlin", 1, 1)
	gs.MakeMove(&PlaceUnitAction{Owner: p2, Unit: gremlin2, Position: NewPosition(6, 2)})
	gremlin3 := NewMinion("gremlin", 1, 1)
	gs.MakeMove(&PlaceUnitAction{Owner: p2, Unit: gremlin3, Position: NewPosition(5, 2)})

	gs.MakeMove(&AttackAction{
		Attacker: p1general,
		Defender: gremlin1,
	})

	assert.Equal(t, NewPosition(-1, -1), gremlin1.GetPosition())
	assert.False(t, gremlin1.IsAlive())
	assert.Equal(t, NewPosition(-1, -1), gremlin2.GetPosition())
	assert.False(t, gremlin2.IsAlive())
	assert.Equal(t, NewPosition(-1, -1), gremlin3.GetPosition())
	assert.False(t, gremlin3.IsAlive())
}

func Test_frenzyAttack(t *testing.T) {
	p1, p2, gs := setupGamestate()

	p1units := p1.Board.GetPlayerUnits(p1)
	p1general := p1units[0]

	p1general.AddAttribute("frenzy", 0)

	gremlin1 := NewMinion("gremlin", 1, 1)
	gs.MakeMove(&PlaceUnitAction{Owner: p2, Unit: gremlin1, Position: NewPosition(0, 1)})
	gremlin2 := NewMinion("gremlin", 1, 1)
	gs.MakeMove(&PlaceUnitAction{Owner: p2, Unit: gremlin2, Position: NewPosition(1, 2)})
	gremlin3 := NewMinion("gremlin", 1, 1)
	gs.MakeMove(&PlaceUnitAction{Owner: p2, Unit: gremlin3, Position: NewPosition(0, 3)})

	gs.MakeMove(&AttackAction{
		Attacker: p1general,
		Defender: gremlin1,
	})

	assert.Equal(t, NewPosition(-1, -1), gremlin1.GetPosition())
	assert.False(t, gremlin1.IsAlive())
	assert.Equal(t, NewPosition(-1, -1), gremlin2.GetPosition())
	assert.False(t, gremlin2.IsAlive())
	assert.Equal(t, NewPosition(-1, -1), gremlin3.GetPosition())
	assert.False(t, gremlin3.IsAlive())

	assert.Equal(t, 24, p1general.GetHp())
	assert.True(t, p1general.IsAlive())
}

func Test_rangedAttack(t *testing.T) {
	p1, p2, gs := setupGamestate()

	p1units := p1.Board.GetPlayerUnits(p1)
	p1general := p1units[0]

	p2units := p1.Board.GetPlayerUnits(p2)
	p2general := p2units[0]

	p1general.AddAttribute("ranged", 0)

	gs.MakeMove(&AttackAction{
		Attacker: p1general,
		Defender: p2general,
	})

	assert.Equal(t, 23, p2general.GetHp())
	assert.Equal(t, 25, p1general.GetHp())
}

func Test_backstabAttack(t *testing.T) {

	p1, p2, gs := setupGamestate()

	p1units := p1.GetUnits()
	p1general := p1units[0]

	p2units := p2.GetUnits()
	p2general := p2units[0]

	p1general.AddAttribute("backstab", 2)

	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: NewPosition(1, 2),
	})

	gs.MakeMove(&MoveAction{
		Unit:     p2general,
		Position: NewPosition(0, 2),
	})

	gs.MakeMove(&AttackAction{
		Attacker: p1general,
		Defender: p2general,
	})

	assert.Equal(t, 25, p1general.GetHp())
	assert.Equal(t, 21, p2general.GetHp())
}

func Test_endTurnAction(t *testing.T) {
	p1, p2, gs := setupGamestate()

	assert.Equal(t, p1, gs.ActivePlayer)

	// Non-active player cannot end the turn
	gs.MakeMove(&EndTurnAction{Owner: p2})
	assert.Equal(t, p1, gs.ActivePlayer)

	gs.MakeMove(&EndTurnAction{Owner: p1})

	assert.Equal(t, p2, gs.ActivePlayer)
}
