package game

import "github.com/RGood/game_engine/pkg/gamestate"

type MoveAction struct {
	Unit     Unit
	Position Position
}

func (ma *MoveAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	ma.Unit.Move(ma.Position)

	return gs
}

type PlaceUnitAction struct {
	Owner    *Player
	Board    *UnitBoard
	Unit     Unit
	Position Position
}

func (ma *PlaceUnitAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	ma.Unit.Place(ma.Owner, ma.Position)
	ma.Unit.Subscribe(gs)

	return gs
}

type RemoveUnitAction struct {
	Unit Unit
}

func (ra *RemoveUnitAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	ra.Unit.Remove()
	ra.Unit.Unsubscribe(gs)

	return gs
}

type DamageAction struct {
	Unit   Unit
	Damage int
}

func (da *DamageAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	da.Unit.Damage(da.Damage)
	if da.Unit.GetHp() <= 0 {
		gs.QueueAction(&RemoveUnitAction{
			Unit: da.Unit,
		})
	}

	return gs
}

type HealAction struct {
	Unit Unit
	Heal int
}

func (ha *HealAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	ha.Unit.Damage(-ha.Heal)

	return gs
}

type AttackAction struct {
	Attacker Unit
	Defender Unit
}

func (aa *AttackAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {

	if aa.Attacker.InRange(aa.Defender) {
		allUnits := []Unit{}
		for unit, _ := range aa.Attacker.GetBoard().Units {
			allUnits = append(allUnits, unit)
		}

		posDiff := aa.Attacker.GetPosition().Diff(aa.Defender.GetPosition())
		wasBackstabbed := aa.Attacker.HasAttribute("backstab") && (posDiff.Y == 0 && ((posDiff.X == -1 && aa.Defender.FacesRight()) || (posDiff.X == 1 && !aa.Defender.FacesRight())))

		collateralDamage := map[Unit]int{}
		collateralDamage[aa.Defender] = aa.Attacker.GetAttack()
		if wasBackstabbed {
			collateralDamage[aa.Defender] += aa.Attacker.GetAttributeValue("backstab")
		}

		if aa.Attacker.HasAttribute("blast") {
			isInline, inlineFunc := aa.Attacker.IsInline(aa.Defender)
			if isInline {
				for _, unit := range filterUnits(filterUnits(allUnits, inlineFunc), aa.Attacker.IsEnemy) {
					collateralDamage[unit] = aa.Attacker.GetAttack()
				}
			}
			// Calculate other enemy units in that line and add to collateral damage
		}

		if aa.Attacker.HasAttribute("frenzy") && aa.Attacker.IsNear(aa.Defender) {
			for _, unit := range filterUnits(filterUnits(allUnits, aa.Attacker.IsNear), aa.Attacker.IsEnemy) {
				collateralDamage[unit] = aa.Attacker.GetAttack()
			}
		}

		for unit, damage := range collateralDamage {
			gs.QueueAction(&DamageAction{Unit: unit, Damage: damage})
		}

		// Do counter-attack check
		// Not backstabbed and ranged or near
		if !wasBackstabbed && (aa.Defender.HasAttribute("ranged") || aa.Defender.IsNear(aa.Attacker)) {
			gs.QueueAction(&DamageAction{Unit: aa.Attacker, Damage: aa.Defender.GetAttack()})
		}
	}

	return gs
}

type EffectAction struct {
	Unit   Unit
	Effect func(Unit)
}

func (ea *EffectAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	ea.Effect(ea.Unit)

	return gs
}

type DispelAction struct {
	Unit Unit
}

func (dispAction *DispelAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	dispAction.Unit.Dispel()

	return gs
}

type SpellAction struct {
	Owner  *Player
	Spell  Spell
	Units  []Unit
	Tiles  []Position
	Effect func()
}

func (sp *SpellAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	sp.Effect()

	return gs
}

type EquipArtifactAction struct {
	Owner    *Player
	Artifact *Artifact
}

func (eaa *EquipArtifactAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	eaa.Artifact.Equip(eaa.Owner, gs)

	return gs
}

type RemoveArtifactAction struct {
	Artifact *Artifact
}

func (raa *RemoveArtifactAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	raa.Artifact.Remove(gs)

	return gs
}

type EndTurnAction struct {
	Owner *Player
}

func (eta *EndTurnAction) Execute(gs *gamestate.Gamestate) *gamestate.Gamestate {
	if gs.ActivePlayer == eta.Owner {
		return gs.EndTurn()
	}

	return gs
}
