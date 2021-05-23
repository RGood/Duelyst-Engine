package gamestate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_equipUnequip(t *testing.T) {
	p1, p2, gs := setupGamestate()
	p2units := gs.Board.GetPlayerUnits(p2)
	p2general := p2units[0]

	brm := NewArtifact("Bloodrage Mask", 1).OnNotify(func(artifact *Artifact, action Action, gamestate *Gamestate) {
		sa, ok := action.(*SpellAction)
		if ok && sa.Owner == artifact.Owner {
			units := []Unit{}
			for unit, _ := range gamestate.Board.Units {
				units = append(units, unit)
			}

			targets := filterUnits(filterUnits(units,
				func(unit Unit) bool {
					return unit.GetOwner() != artifact.Owner
				}),
				func(unit Unit) bool {
					return unit.GetType() == "general"
				})

			for _, target := range targets {
				gs.QueueAction(&DamageAction{
					Unit:   target,
					Damage: 1,
				})
			}
		}
	})

	dummyArtifact := NewArtifact("Dummy", 0)
	nullSpell := NewGenericSpell("Test", 0, func(arg1 *Player, arg2 *Gamestate, arg3 []Unit, arg4 []Position) {

	})

	brm.Equip(p1, gs)
	dummyArtifact.Equip(p2, gs)
	nullSpell.Cast(p1, gs, nil, nil)
	nullSpell.Cast(p1, gs, nil, nil)
	nullSpell.Cast(p1, gs, nil, nil)

	assert.Equal(t, 0, dummyArtifact.Charges)
	assert.Equal(t, 22, p2general.GetHp())
}

func Test_arclyteRegalia(t *testing.T) {
	p1, p2, gs := setupGamestate()
	p1units := gs.Board.GetPlayerUnits(p1)
	p1general := p1units[0]
	p2units := gs.Board.GetPlayerUnits(p2)
	p2general := p2units[0]

	damageCounter := 0
	arclyteRegalia := NewArtifact(
		"Arclyte Regalia",
		4,
	).OnEquip(func(artifact *Artifact, gamestate *Gamestate) {
		general := filterUnits(gamestate.Board.GetPlayerUnits(artifact.Owner), func(unit Unit) bool {
			return unit.GetType() == "general"
		})[0]

		general.BuffAttack(2)
	}).OnUnEquip(func(artifact *Artifact, gamestate *Gamestate) {
		general := filterUnits(gamestate.Board.GetPlayerUnits(artifact.Owner), func(unit Unit) bool {
			return unit.GetType() == "general"
		})[0]

		general.BuffAttack(-2)
	}).OnIntercept(func(artifact *Artifact, action Action, gamestate *Gamestate) Action {
		_, ok := action.(*EndTurnAction)
		if ok {
			damageCounter = 0
			return action
		}

		general := filterUnits(gamestate.Board.GetPlayerUnits(artifact.Owner), func(unit Unit) bool {
			return unit.GetType() == "general"
		})[0]
		damageAction, ok := action.(*DamageAction)
		if ok {
			if damageCounter == 0 && damageAction.Unit == general {
				damageCounter++
				damageAction.Damage -= 2
				if damageAction.Damage < 0 {
					damageAction.Damage = 0
				}
				return damageAction
			}
		}

		return action
	})

	arclyteRegalia.Equip(p1, gs)

	assert.Equal(t, 4, p1general.GetAttack())

	infrontOfP2 := p2general.GetPosition().Diff(NewPosition(1, 0))
	gs.MakeMove(&MoveAction{
		Unit:     p1general,
		Position: infrontOfP2,
	})

	gs.MakeMove(&AttackAction{
		Attacker: p1general,
		Defender: p2general,
	})

	assert.Equal(t, 21, p2general.GetHp())
	assert.Equal(t, 25, p1general.GetHp())

	gs.MakeMove(&EndTurnAction{})

	gs.MakeMove(&AttackAction{
		Attacker: p2general,
		Defender: p1general,
	})

	assert.Equal(t, 17, p2general.GetHp())
	assert.Equal(t, 25, p1general.GetHp())

	gs.MakeMove(&AttackAction{
		Attacker: p2general,
		Defender: p1general,
	})

	assert.Equal(t, 13, p2general.GetHp())
	assert.Equal(t, 23, p1general.GetHp())
	assert.Equal(t, 2, arclyteRegalia.Charges)

	gs.MakeMove(&AttackAction{
		Attacker: p2general,
		Defender: p1general,
	})

	gs.MakeMove(&AttackAction{
		Attacker: p2general,
		Defender: p1general,
	})

	assert.Equal(t, 0, arclyteRegalia.Charges)
	assert.Equal(t, 2, p1general.GetAttack())

}
