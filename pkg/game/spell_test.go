package game

import (
	"testing"

	"github.com/RGood/game_engine/pkg/gamestate"
	"github.com/stretchr/testify/assert"
)

func Test_boardSpell(t *testing.T) {
	p1, _, gs := setupGamestate()
	units := p1.GetUnits()
	p1general := units[0]

	assert.Equal(t, NewPosition(0, 2), p1general.GetPosition())
	p1general.AddAttribute("ranged", 0)
	assert.True(t, p1general.HasAttribute("ranged"))

	NewGenericSpell("EMP", 4, func(player *Player, gs *gamestate.Gamestate, _ []Unit, _ []Position) {
		for unit, _ := range player.Board.Units {
			gs.QueueAction(&EffectAction{
				Unit: unit,
				Effect: func(u Unit) {
					u.Dispel()
				},
			})
		}
	}).Cast(p1, gs, nil, nil)

	assert.False(t, p1general.HasAttribute("ranged"))
}

func Test_singleTargetSpell(t *testing.T) {
	p1, p2, gs := setupGamestate()
	p1units := p1.GetUnits()
	p1general := p1units[0]

	p2units := p2.GetUnits()
	p2general := p2units[0]

	assert.Equal(t, NewPosition(0, 2), p1general.GetPosition())
	assert.Equal(t, NewPosition(8, 2), p2general.GetPosition())
	assert.Equal(t, 25, p2general.GetHp())

	phoenixFire := NewDamageSpell("Phoenix Fire", 2, 3, func(owner *Player, game *gamestate.Gamestate, damage int, targets []Unit, _ []Position) {
		if len(targets) == 1 {
			game.QueueAction(&DamageAction{
				Unit:   targets[0],
				Damage: damage,
			})
		}
	})

	phoenixFire.Cast(p1, gs, []Unit{p2general}, nil)
	assert.Equal(t, 22, p2general.GetHp())
}

func Test_multiTargetSpell(t *testing.T) {
	p1, p2, gs := setupGamestate()

	gremlin := NewMinion("gremlin", 1, 1)
	goblin := NewMinion("goblin", 1, 1)

	gremlin.Place(p1, NewPosition(1, 2))
	goblin.Place(p2, NewPosition(7, 2))

	juxtaposition := NewGenericSpell("Juxtaposition", 0, func(owner *Player, gs *gamestate.Gamestate, units []Unit, _ []Position) {
		if len(units) == 2 {
			p1 := units[0].GetPosition()
			p2 := units[1].GetPosition()

			gs.QueueAction(&MoveAction{Unit: units[0], Position: p2})
			gs.QueueAction(&MoveAction{Unit: units[1], Position: p1})
		}
	})

	juxtaposition.Cast(p1, gs, []Unit{gremlin, goblin}, nil)

	assert.Equal(t, NewPosition(7, 2), gremlin.GetPosition())
	assert.Equal(t, NewPosition(1, 2), goblin.GetPosition())
}

func Test_tileSpell(t *testing.T) {
	p1, p2, gs := setupGamestate()

	gremlin := NewMinion("gremlin", 1, 1)
	goblin := NewMinion("goblin", 1, 1)
	gremlin.Place(p1, NewPosition(1, 2))
	goblin.Place(p2, NewPosition(7, 2))

	chromaticCold := NewGenericSpell("Chromatic Cold", 2, func(owner *Player, gs *gamestate.Gamestate, _ []Unit, tiles []Position) {
		if len(tiles) == 1 {
			targetTile := tiles[0]
			if _, ok := owner.Board.Positions[targetTile]; ok {
				unit := owner.Board.Positions[targetTile]
				if unit.GetOwner() != owner {
					gs.QueueAction(&DamageAction{
						Unit:   unit,
						Damage: 1,
					})
				}

				gs.QueueAction(&DispelAction{
					Unit: unit,
				})
			}
		}

	})

	chromaticCold.Cast(p2, gs, nil, []Position{NewPosition(1, 2)})
	chromaticCold.Cast(p2, gs, nil, []Position{NewPosition(7, 2)})

	assert.False(t, gremlin.IsAlive())
	assert.True(t, goblin.IsAlive())
}

func Test_multiTileSpell(t *testing.T) {
	p1, _, gs := setupGamestate()

	bonechillBarrier := NewGenericSpell("Bonechill Barrier", 2, func(owner *Player, gs *gamestate.Gamestate, _ []Unit, tiles []Position) {
		if len(tiles) <= 3 {
			for _, tile := range tiles {
				token := NewWall("Bonechill Barrier", owner, 2, 0)
				gs.QueueAction(&PlaceUnitAction{
					Owner:    owner,
					Unit:     token,
					Position: tile,
				})
			}
		}
	})

	bonechillBarrier.Cast(p1, gs, nil, []Position{
		NewPosition(1, 1),
		NewPosition(2, 1),
		NewPosition(3, 0),
	})

	assert.Equal(t, 4, len(p1.GetUnits()))

	walls := []Unit{}
	for _, unit := range p1.GetUnits() {
		if unit.HasSubtype("wall") {
			walls = append(walls, unit)
		}
	}

	assert.Equal(t, 3, len(walls))

	chromaticCold := NewGenericSpell("Chromatic Cold", 2, func(owner *Player, gs *gamestate.Gamestate, _ []Unit, tiles []Position) {
		if len(tiles) == 1 {
			targetTile := tiles[0]
			if _, ok := owner.Board.Positions[targetTile]; ok {
				unit := owner.Board.Positions[targetTile]
				if unit.GetOwner() != owner {
					gs.QueueAction(&DamageAction{
						Unit:   unit,
						Damage: 1,
					})
				}

				gs.QueueAction(&DispelAction{
					Unit: unit,
				})
			}
		}
	})

	wall := walls[0]

	chromaticCold.Cast(p1, gs, nil, []Position{wall.GetPosition()})
	assert.False(t, wall.IsAlive())

	assert.Equal(t, 3, len(p1.GetUnits()))
}

func Test_dispelMovedWall(t *testing.T) {
	p1, _, gs := setupGamestate()

	bonechillBarrier := NewGenericSpell("Bonechill Barrier", 2, func(owner *Player, gs *gamestate.Gamestate, _ []Unit, tiles []Position) {
		if len(tiles) <= 3 {
			for _, tile := range tiles {
				token := NewWall("Bonechill Barrier", owner, 2, 0)
				gs.QueueAction(&PlaceUnitAction{
					Owner:    owner,
					Unit:     token,
					Position: tile,
				})
			}
		}
	})

	bonechillBarrier.Cast(p1, gs, nil, []Position{
		NewPosition(1, 1),
		NewPosition(2, 1),
		NewPosition(3, 0),
	})

	assert.Equal(t, 4, len(p1.GetUnits()))

	walls := []Unit{}
	for _, unit := range p1.GetUnits() {
		if unit.HasSubtype("wall") {
			walls = append(walls, unit)
		}
	}

	assert.Equal(t, 3, len(walls))

	chromaticCold := NewGenericSpell("Chromatic Cold", 2, func(owner *Player, gs *gamestate.Gamestate, _ []Unit, tiles []Position) {
		if len(tiles) == 1 {
			targetTile := tiles[0]
			if _, ok := owner.Board.Positions[targetTile]; ok {
				unit := owner.Board.Positions[targetTile]
				if unit.GetOwner() != owner {
					gs.QueueAction(&DamageAction{
						Unit:   unit,
						Damage: 1,
					})
				}

				gs.QueueAction(&DispelAction{
					Unit: unit,
				})
			}
		}
	})

	wall := walls[0]

	gs.MakeMove(&MoveAction{
		wall,
		NewPosition(4, 2),
	})

	chromaticCold.Cast(p1, gs, nil, []Position{wall.GetPosition()})
	assert.False(t, wall.IsAlive())

	assert.Equal(t, 3, len(p1.GetUnits()))
	assert.Equal(t, nil, p1.Board.Positions[NewPosition(4, 2)])
}

func Test_endOfTurnListener(t *testing.T) {
	p1, p2, gs := setupGamestate()
	p1units := p1.GetUnits()
	p1general := p1units[0]
	p2units := p2.GetUnits()
	p2general := p2units[0]

	// Define Eight Gates spell
	eightGates := NewGenericSpell("Eight Gates", 2, func(owner *Player, game *gamestate.Gamestate, _ []Unit, _ []Position) {
		NewUntilEndOfTurnInterceptor(
			func(action gamestate.Action, validateGame *gamestate.Gamestate) bool {
				spellAction, ok := action.(*SpellAction)
				if ok {
					_, dsOk := spellAction.Spell.(*DamageSpell)
					return dsOk
				}

				return false
			},
			func(_ gamestate.Interceptor, action gamestate.Action, executionGamestate *gamestate.Gamestate) gamestate.Action {
				sa, _ := action.(*SpellAction)
				ds, _ := sa.Spell.(*DamageSpell)

				return &SpellAction{
					Owner: sa.Owner,
					Spell: ds,
					Effect: func() {
						ds.Effect(sa.Owner, executionGamestate, ds.Damage+2, sa.Units, sa.Tiles)
					},
				}
			},
		).Subscribe(game)
	})

	// Define Phoenix Fire spell
	phoenixFire := NewDamageSpell("Phoenix Fire", 2, 3, func(owner *Player, game *gamestate.Gamestate, damage int, units []Unit, _ []Position) {
		if len(units) == 1 {
			game.QueueAction(&DamageAction{
				Unit:   units[0],
				Damage: damage,
			})
		}
	})

	// Cast Eight Gates
	eightGates.Cast(p1, gs, nil, nil)
	// Cast Phoenix Fire (should be buffed) (25 -> 15)
	phoenixFire.Cast(p1, gs, []Unit{p2general}, nil)
	phoenixFire.Cast(p1, gs, []Unit{p2general}, nil)

	// End the turn
	gs.MakeMove(&EndTurnAction{Owner: p1})

	// Cast Phoenix Fire (should not be buffed) (25 -> 22)
	phoenixFire.Cast(p2, gs, []Unit{p1general}, nil)

	assert.Equal(t, 15, p2general.GetHp())
	assert.Equal(t, 22, p1general.GetHp())
}

func Test_saberspineSeal(t *testing.T) {
	p1, p2, gs := setupGamestate()
	p1units := p1.GetUnits()
	p1general := p1units[0]
	p2units := p2.GetUnits()
	p2general := p2units[0]

	infrontOfP2 := p2general.GetPosition().Diff(NewPosition(1, 0))

	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: infrontOfP2,
	})

	assert.True(t, p1general.IsNear(p2general))

	saberspineSeal := NewGenericSpell("Saberspine Seal", 1, func(owner *Player, game *gamestate.Gamestate, targets []Unit, _ []Position) {
		if len(targets) == 1 {
			targetUnit := targets[0]
			targetUnit.BuffAttack(3)

			NewUntilEndOfTurnListener(
				func(action gamestate.Action, gamestate *gamestate.Gamestate) bool {
					_, ok := action.(*EndTurnAction)
					return ok
				},
				func(_ gamestate.Listener, action gamestate.Action, gamestate *gamestate.Gamestate) {
					targetUnit.BuffAttack(-3)
				},
			).Subscribe(game)
		}
	})

	saberspineSeal.Cast(p1, gs, []Unit{p1general}, nil)

	gs.MakeMove(&AttackAction{
		Attacker: p1general,
		Defender: p2general,
	})

	assert.Equal(t, 5, p1general.GetAttack())

	gs.MakeMove(&EndTurnAction{})

	assert.Equal(t, 2, p1general.GetAttack())

	assert.Equal(t, 23, p1general.GetHp())
	assert.Equal(t, 20, p2general.GetHp())

	saberspineSeal.Cast(p2, gs, []Unit{p2general}, nil)
	saberspineSeal.Cast(p2, gs, []Unit{p2general}, nil)

	assert.Equal(t, 8, p2general.GetAttack())

	gs.MakeMove(&AttackAction{
		Attacker: p2general,
		Defender: p1general,
	})

	assert.Equal(t, 15, p1general.GetHp())
	assert.Equal(t, 18, p2general.GetHp())

	gs.MakeMove(&EndTurnAction{})

	assert.Equal(t, 2, p1general.GetAttack())
	assert.Equal(t, 2, p2general.GetAttack())
}
