package game

type Player struct {
	id          *string
	General     string
	StartingPos Position
	FacesRight  bool
	Board       *UnitBoard
}

func NewPlayer(id string, general string, board *UnitBoard, pos Position, right bool) *Player {
	player := &Player{
		id:         &id,
		General:    general,
		FacesRight: right,
		Board:      board,
	}

	generalUnit := NewUnitFactory().SetName(general).SetHealth(25).SetAttack(2).SetUnitType("general").Create()
	generalUnit.Place(player, pos)

	board.PlaceUnit(generalUnit, pos)

	return player
}

func (p *Player) IsAlive() bool {
	for unit, _ := range p.Board.Units {
		if unit.GetOwner() == p && unit.GetType() == "general" {
			return true
		}
	}

	return false
}

func (p *Player) GetUnits() []Unit {
	units := []Unit{}
	for unit, _ := range p.Board.Units {
		if unit.GetOwner() == p {
			units = append(units, unit)
		}
	}

	return units
}

func (p *Player) GetGeneral() Unit {
	if p.IsAlive() {
		for unit, _ := range p.Board.Units {
			if unit.GetOwner() == p && unit.GetType() == "general" {
				return unit
			}
		}
	}

	return nil
}
