package game

import "github.com/RGood/game_engine/pkg/gamestate"

type Unit interface {
	GetName() string
	GetOwner() *Player
	GetType() string
	HasSubtype(string) bool
	GetHp() int
	IsAlive() bool
	GetAttack() int
	GetWalkDistance() int
	AddAttribute(string, int)
	RemoveAttribute(string)
	GetAttributeValue(string) int
	HasAttribute(string) bool
	IsEnemy(Unit) bool
	InRange(Unit) bool
	IsNear(Unit) bool
	IsInline(Unit) (bool, func(Unit) bool)
	Damage(int)
	GetPosition() Position
	FacesRight() bool
	GetValidMoves() map[Position]struct{}
	Move(Position)
	Place(*Player, Position)
	Remove()
	Dispel()
	BuffAttack(int)
	BuffHealth(int)
	SetBoard(*UnitBoard)
	GetBoard() *UnitBoard
	AddActionTrigger(ActionTrigger) int
	RemoveActionTrigger(int)
	Subscribe(*gamestate.Gamestate)
	Unsubscribe(*gamestate.Gamestate)
	Notify(gamestate.Action, *gamestate.Gamestate)
	Apply(gamestate.Action, *gamestate.Gamestate) gamestate.Action
}

type Minion struct {
	name             string
	unitType         string
	subtypes         map[string]struct{}
	owner            *Player
	faceRight        bool
	walkDistance     int
	baseHp           int
	hpDelta          int
	baseAttack       int
	attackDelta      int
	damage           int
	attributes       map[string]int
	board            *UnitBoard
	triggerCount     int
	triggers         map[int]ActionTrigger
	interceptorCount int
	interceptors     map[int]InterceptTrigger
}

func Equal(u1, u2 Unit) bool {
	return u1.GetName() == u2.GetName() &&
		u1.GetBoard() == u2.GetBoard() &&
		u1.GetPosition() == u2.GetPosition()
}

type ActionTrigger struct {
	Trigger   func(gamestate.Action, *gamestate.Gamestate)
	CanDispel bool
}

type InterceptTrigger struct {
	Trigger   func(gamestate.Action, *gamestate.Gamestate) gamestate.Action
	CanDispel bool
}

type UnitFactory struct {
	name             string
	unitType         string
	subtypes         map[string]struct{}
	hp               int
	attack           int
	attributes       map[string]int
	triggerCount     int
	triggers         map[int]ActionTrigger
	interceptorCount int
	interceptors     map[int]InterceptTrigger
}

func NewUnitFactory() *UnitFactory {
	return &UnitFactory{
		subtypes:         map[string]struct{}{},
		attributes:       map[string]int{},
		triggerCount:     0,
		triggers:         map[int]ActionTrigger{},
		interceptorCount: 0,
		interceptors:     map[int]InterceptTrigger{},
	}
}

func (uf *UnitFactory) SetName(name string) *UnitFactory {
	uf.name = name
	return uf
}

func (uf *UnitFactory) SetUnitType(unitType string) *UnitFactory {
	uf.unitType = unitType
	return uf
}

func (uf *UnitFactory) AddSubtype(subtype string) *UnitFactory {
	uf.subtypes[subtype] = struct{}{}
	return uf
}

func (uf *UnitFactory) SetHealth(hp int) *UnitFactory {
	uf.hp = hp
	return uf
}

func (uf *UnitFactory) SetAttack(attack int) *UnitFactory {
	uf.attack = attack
	return uf
}

func (uf *UnitFactory) AddAttribute(name string, value int) *UnitFactory {
	uf.attributes[name] = value
	return uf
}

func (uf *UnitFactory) AddTrigger(trigger ActionTrigger) *UnitFactory {
	triggerId := uf.triggerCount
	uf.triggerCount++
	uf.triggers[triggerId] = trigger
	return uf
}

func (uf *UnitFactory) AddIntercept(trigger InterceptTrigger) *UnitFactory {
	interceptorId := uf.interceptorCount
	uf.interceptorCount++
	uf.interceptors[interceptorId] = trigger
	return uf
}

func (uf *UnitFactory) Create() Unit {
	return NewUnit(
		uf.name,
		uf.unitType,
		uf.subtypes,
		uf.attributes,
		uf.hp,
		uf.attack,
		uf.triggers,
		uf.interceptors,
	)
}

func NewUnit(name string, unitType string, subtypes map[string]struct{}, attributes map[string]int, hp int, attack int, triggers map[int]ActionTrigger, interceptors map[int]InterceptTrigger) Unit {
	return &Minion{
		name:             name,
		unitType:         unitType,
		subtypes:         subtypes,
		attributes:       attributes,
		baseHp:           hp,
		baseAttack:       attack,
		walkDistance:     2,
		triggerCount:     len(triggers),
		triggers:         triggers,
		interceptorCount: len(interceptors),
		interceptors:     interceptors,
	}
}

func NewMinion(name string, hp int, attack int) *Minion {

	minion := &Minion{
		name:         name,
		unitType:     "minion",
		walkDistance: 2,
		baseHp:       hp,
		hpDelta:      0,
		baseAttack:   attack,
		attackDelta:  0,
		damage:       0,
		attributes:   map[string]int{},
		triggers:     map[int]ActionTrigger{},
	}

	return minion
}

func NewGeneral(name string, owner *Player) *Minion {
	general := &Minion{
		name:         name,
		owner:        owner,
		unitType:     "general",
		faceRight:    owner.FacesRight,
		walkDistance: 2,
		baseHp:       25,
		hpDelta:      0,
		baseAttack:   2,
		attackDelta:  0,
		damage:       0,
		attributes:   map[string]int{},
		triggers:     map[int]ActionTrigger{},
	}

	return general
}

func NewWall(name string, owner *Player, hp int, attack int) *Minion {
	wall := &Minion{
		name:         name,
		owner:        owner,
		unitType:     "token",
		subtypes:     map[string]struct{}{"wall": struct{}{}},
		faceRight:    owner.FacesRight,
		walkDistance: 0,
		baseHp:       hp,
		hpDelta:      0,
		baseAttack:   attack,
		attackDelta:  0,
		damage:       0,
		attributes:   map[string]int{},
		triggers:     map[int]ActionTrigger{},
	}

	wall.AddActionTrigger(ActionTrigger{
		Trigger: func(action gamestate.Action, gs *gamestate.Gamestate) {
			dispelAction, ok := action.(*DispelAction)
			if ok {
				if Equal(dispelAction.Unit, wall) {
					gs.QueueAction(&RemoveUnitAction{
						Unit: wall,
					})
				}
			}
		},
		CanDispel: false,
	})

	return wall
}

func (m *Minion) GetType() string {
	return m.unitType
}

func (m *Minion) HasSubtype(subtype string) bool {
	_, ok := m.subtypes[subtype]
	return ok
}

func (m *Minion) GetName() string {
	return m.name
}

func (m *Minion) GetOwner() *Player {
	return m.owner
}

func (m *Minion) GetHp() int {
	return m.baseHp + m.hpDelta - m.damage
}

func (m *Minion) GetAttack() int {
	return m.baseAttack + m.attackDelta
}

func (m *Minion) GetWalkDistance() int {
	return m.walkDistance
}

func (m *Minion) HasAttribute(attr string) bool {
	_, ok := m.attributes[attr]
	return ok
}

func (m *Minion) IsEnemy(u Unit) bool {
	return m.GetOwner() != u.GetOwner()
}

func (m *Minion) IsNear(u Unit) bool {
	pos1 := m.GetPosition()
	pos2 := u.GetPosition()
	absDiff := pos1.Diff(pos2).Abs()
	return max(absDiff.X, absDiff.Y) == 1
}

func (m *Minion) isInlineNorth(u Unit) bool {
	pos1 := m.GetPosition()
	pos2 := u.GetPosition()
	posDiff := pos1.Diff(pos2)

	return posDiff.X == 0 && posDiff.Y > 0
}

func (m *Minion) isInlineSouth(u Unit) bool {
	pos1 := m.GetPosition()
	pos2 := u.GetPosition()
	posDiff := pos1.Diff(pos2)

	return posDiff.X == 0 && posDiff.Y < 0
}

func (m *Minion) isInlineEast(u Unit) bool {
	pos1 := m.GetPosition()
	pos2 := u.GetPosition()
	posDiff := pos1.Diff(pos2)

	return posDiff.X < 0 && posDiff.Y == 0
}

func (m *Minion) isInlineWest(u Unit) bool {
	pos1 := m.GetPosition()
	pos2 := u.GetPosition()
	posDiff := pos1.Diff(pos2)

	return posDiff.X > 0 && posDiff.Y == 0
}

func (m *Minion) IsInline(u Unit) (bool, func(Unit) bool) {
	if m.isInlineNorth(u) {
		return true, m.isInlineNorth
	} else if m.isInlineSouth(u) {
		return true, m.isInlineSouth
	} else if m.isInlineEast(u) {
		return true, m.isInlineEast
	} else if m.isInlineWest(u) {
		return true, m.isInlineWest
	} else {
		return false, nil
	}
}

func filterUnits(units []Unit, filterFunc func(Unit) bool) []Unit {
	filteredUnits := []Unit{}
	for _, unit := range units {
		if filterFunc(unit) {
			filteredUnits = append(filteredUnits, unit)
		}
	}

	return filteredUnits
}

func (m *Minion) Damage(dmg int) {
	m.damage += dmg
	if m.damage <= 0 {
		m.damage = 0
	}
}

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func abs(x int) int {
	if x < 0 {
		return x * -1
	} else {
		return x
	}
}

func (m *Minion) InRange(u Unit) bool {
	// If the unit has ranged it can attack
	if m.HasAttribute("ranged") {
		return true
	}

	pos1 := m.GetPosition()
	pos2 := u.GetPosition()
	absDiff := pos1.Diff(pos2).Abs()

	// If the unit has blast and is orthogonal to the enemey, it can attack
	if m.HasAttribute("blast") {
		isInline, _ := m.IsInline(u)
		if isInline {
			return true
		}
	}

	// If the unit is nearby, it can attack
	if max(absDiff.X, absDiff.Y) == 1 {
		return true
	}

	return false
}

func (m *Minion) GetPosition() Position {
	if m.board != nil {
		pos := m.GetBoard().GetPosition(m)
		return pos
	} else {
		return NewPosition(-1, -1)
	}
}

func (m *Minion) IsAlive() bool {
	return m.board != nil
}

func (m *Minion) FacesRight() bool {
	return m.faceRight
}

func (m *Minion) GetValidMoves() map[Position]struct{} {
	positions := map[Position]struct{}{}

	if m.board != nil {
		return m.board.GetValidMoves(m)
	}

	return positions
}

func (m *Minion) Move(pos Position) {
	if m.board != nil && pos.IsOnBoard(m.board) {
		oldPos := m.board.Units[m]
		delete(m.board.Positions, oldPos)
		m.board.Units[m] = pos
		m.board.Positions[pos] = m
	}
}

func (m *Minion) Place(owner *Player, pos Position) {
	m.owner = owner
	owner.Board.PlaceUnit(m, pos)
}

func (m *Minion) Remove() {
	if m.board != nil {
		m.board.RemoveUnit(m)
	}
	m.owner = nil
}

func (m *Minion) Dispel() {
	m.hpDelta = 0
	m.attackDelta = 0
	m.attributes = map[string]int{}

	for id, trigger := range m.triggers {
		if trigger.CanDispel {
			delete(m.triggers, id)
		}
	}

	for id, interceptor := range m.interceptors {
		if interceptor.CanDispel {
			delete(m.interceptors, id)
		}
	}
}

func (m *Minion) BuffAttack(delta int) {
	m.attackDelta += delta
}

func (m *Minion) BuffHealth(delta int) {
	m.hpDelta += delta
}

func (m *Minion) AddAttribute(attr string, value int) {
	m.attributes[attr] = value
}

func (m *Minion) GetAttributeValue(attr string) int {
	val, _ := m.attributes[attr]
	return val
}

func (m *Minion) RemoveAttribute(attr string) {
	delete(m.attributes, attr)
}

func (m *Minion) SetBoard(ub *UnitBoard) {
	m.board = ub
}

func (m *Minion) GetBoard() *UnitBoard {
	return m.board
}

func (m *Minion) AddActionTrigger(trigger ActionTrigger) int {
	triggerId := m.triggerCount
	m.triggerCount++
	m.triggers[triggerId] = trigger

	return triggerId
}

func (m *Minion) RemoveActionTrigger(triggerId int) {
	delete(m.triggers, triggerId)
}

func (m *Minion) Subscribe(gs *gamestate.Gamestate) {
	gs.Subscribe(m)
	gs.AddInterceptor(m)
}

func (m *Minion) Unsubscribe(gs *gamestate.Gamestate) {
	gs.Unsubscribe(m)
	gs.RemoveInterceptor(m)
}

func (m *Minion) Notify(action gamestate.Action, gs *gamestate.Gamestate) {
	for _, trigger := range m.triggers {
		trigger.Trigger(action, gs)
	}
}

func (m *Minion) Apply(action gamestate.Action, gs *gamestate.Gamestate) gamestate.Action {
	for _, interceptor := range m.interceptors {
		action = interceptor.Trigger(action, gs)
	}

	return action
}
