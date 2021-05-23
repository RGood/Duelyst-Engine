package gamestate

type Spell interface {
	Cast(*Player, *Gamestate, []Unit, []Position)
}

type GenericSpell struct {
	Name   string
	Cost   int
	Effect func(*Player, *Gamestate, []Unit, []Position)
}

func NewGenericSpell(name string, cost int, effect func(*Player, *Gamestate, []Unit, []Position)) *GenericSpell {
	return &GenericSpell{
		Name:   name,
		Cost:   cost,
		Effect: effect,
	}
}

func (spell *GenericSpell) Cast(owner *Player, gs *Gamestate, units []Unit, positions []Position) {
	gs.MakeMove(&SpellAction{
		Owner: owner,
		Spell: spell,
		Effect: func() {
			spell.Effect(owner, gs, units, positions)
		},
	})
}

type DamageSpell struct {
	Name   string
	Cost   int
	Damage int
	Effect func(*Player, *Gamestate, int, []Unit, []Position)
}

func NewDamageSpell(name string, cost int, damage int, effect func(*Player, *Gamestate, int, []Unit, []Position)) *DamageSpell {
	return &DamageSpell{
		Name:   name,
		Cost:   cost,
		Damage: damage,
		Effect: effect,
	}
}

func (ds *DamageSpell) Cast(owner *Player, gs *Gamestate, units []Unit, positions []Position) {
	gs.MakeMove(&SpellAction{
		Owner: owner,
		Spell: ds,
		Units: units,
		Effect: func() {
			ds.Effect(owner, gs, ds.Damage, units, positions)
		},
	})
}
