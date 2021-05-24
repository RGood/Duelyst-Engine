package game

type Event struct {
	Type   string
	Source Unit
	Value  int
	Owner  *Player
}

type UnitBoard struct {
	BoardX, BoardY int
	Units          map[Unit]Position
	Positions      map[Position]Unit
}

type Position struct {
	X, Y int
}

func (p Position) Diff(op Position) Position {
	return NewPosition(
		p.X-op.X,
		p.Y-op.Y,
	)
}

func (p Position) Add(op Position) Position {
	return NewPosition(
		p.X+op.X,
		p.Y+op.Y,
	)
}

func (p Position) Abs() Position {
	return NewPosition(
		abs(p.X),
		abs(p.Y),
	)
}

func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

func (p Position) IsOnBoard(ub *UnitBoard) bool {
	return p.X >= 0 && p.X < ub.BoardX && p.Y >= 0 && p.Y < ub.BoardY
}

func NewUnitBoard(x, y int) *UnitBoard {
	return &UnitBoard{
		BoardX:    x,
		BoardY:    y,
		Units:     map[Unit]Position{},
		Positions: map[Position]Unit{},
	}
}

func (ub *UnitBoard) IsOccupied(pos Position) bool {
	_, ok := ub.Positions[pos]
	return ok
}

func (ub *UnitBoard) GetPosition(unit Unit) Position {
	pos, _ := ub.Units[unit]
	return pos
}

func (ub *UnitBoard) PlaceUnit(unit Unit, pos Position) bool {
	if ub.IsOccupied(pos) {
		return false
	}

	if pos.X < 0 || pos.X >= ub.BoardX || pos.Y < 0 || pos.Y >= ub.BoardY {
		return false
	}

	ub.Units[unit] = pos
	ub.Positions[pos] = unit
	unit.SetBoard(ub)

	return true
}

func (ub *UnitBoard) RemoveUnit(unit Unit) {
	pos := ub.Units[unit]
	delete(ub.Units, unit)
	delete(ub.Positions, pos)
	unit.SetBoard(nil)
}

func (ub *UnitBoard) GetPlayerUnits(owner *Player) []Unit {
	playerUnits := []Unit{}
	for unit, _ := range ub.Units {
		if unit.GetOwner() == owner {
			playerUnits = append(playerUnits, unit)
		}
	}

	return playerUnits
}

func (ub *UnitBoard) GetValidTargets(unit Unit) map[Unit]struct{} {
	validTargets := map[Unit]struct{}{}
	for otherUnit, _ := range ub.Units {
		if otherUnit.GetOwner() != unit.GetOwner() && unit.InRange(otherUnit) {
			validTargets[otherUnit] = struct{}{}
		}
	}

	return validTargets
}

func (ub *UnitBoard) GetValidMoves(unit Unit) map[Position]struct{} {
	validMoves := map[Position]struct{}{}

	// Add starting tile
	validMoves[unit.GetPosition()] = struct{}{}

	// For Unit.WalkRange
	for i := 0; i < unit.GetWalkDistance(); i++ {
		// Add tiles that don't have enemies on them 1 tile away from all added tiles
		nextValidMove := map[Position]struct{}{}
		for move, _ := range validMoves {
			possibleMoves := []Position{
				move.Add(NewPosition(0, 1)),
				move.Add(NewPosition(1, 0)),
				move.Add(NewPosition(0, -1)),
				move.Add(NewPosition(-1, 0)),
			}

			for _, pm := range possibleMoves {
				_, vmOk := validMoves[pm]
				_, nvmOk := nextValidMove[pm]
				if (!vmOk && !nvmOk) && pm.IsOnBoard(ub) && (ub.Positions[pm] == nil || ub.Positions[pm].GetOwner() == unit.GetOwner()) {
					nextValidMove[pm] = struct{}{}
				}
			}
		}

		for nvm, _ := range nextValidMove {
			validMoves[nvm] = struct{}{}
		}
	}

	// Remove all tiles with units on them
	for vm, _ := range validMoves {
		if ub.Positions[vm] != nil {
			delete(validMoves, vm)
		}
	}

	return validMoves
}
