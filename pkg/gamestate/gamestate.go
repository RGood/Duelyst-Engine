package gamestate

type Player struct {
	id          *string
	General     string
	StartingPos Position
	FacesRight  bool
}

func NewPlayer(id string, general string, pos Position, right bool) *Player {
	return &Player{
		id:          &id,
		General:     general,
		StartingPos: pos,
		FacesRight:  right,
	}
}

func (p *Player) IsAlive(gs *Gamestate) bool {
	for unit, _ := range gs.Board.Units {
		if unit.GetOwner() == p && unit.GetType() == "general" {
			return true
		}
	}

	return false
}

type Gamestate struct {
	Players      []*Player
	ActivePlayer *Player
	Board        *UnitBoard
	actions      []Action

	listeners    map[Listener]struct{}
	interceptors map[Interceptor]struct{}

	ended bool
}

func NewGamestate(players ...*Player) *Gamestate {
	gs := &Gamestate{
		Players:      players,
		ActivePlayer: players[0],
		Board:        NewUnitBoard(9, 5),
		actions:      []Action{},
		ended:        false,
		listeners:    map[Listener]struct{}{},
		interceptors: map[Interceptor]struct{}{},
	}

	for _, player := range players {
		general := NewGeneral(player.General, player)
		general.Place(gs.Board, player.StartingPos)

		gs.Subscribe(general)
	}

	return gs
}

func (gs *Gamestate) Subscribe(l Listener) {
	gs.listeners[l] = struct{}{}
}

func (gs *Gamestate) Unsubscribe(l Listener) {
	delete(gs.listeners, l)
}

func (gs *Gamestate) AddInterceptor(i Interceptor) {
	gs.interceptors[i] = struct{}{}
}

func (gs *Gamestate) RemoveInterceptor(i Interceptor) {
	delete(gs.interceptors, i)
}

func (gs *Gamestate) QueueAction(action Action) {
	gs.actions = append(gs.actions, action)
}

func (gs *Gamestate) MakeMove(action Action) *Gamestate {
	if gs.HasEnded() {
		return gs
	}

	gs.actions = append(gs.actions, action)
	for len(gs.actions) > 0 {
		activeMove := gs.actions[0]
		gs.actions = gs.actions[1:]

		for i, _ := range gs.interceptors {
			activeMove = i.Apply(activeMove, gs)
		}

		activeMove.Execute(gs)

		for listener, _ := range gs.listeners {
			listener.Notify(activeMove, gs)
		}

	}

	return gs
}

func (gs *Gamestate) EndTurn() *Gamestate {
	if gs.HasEnded() {
		return gs
	}

	apIndex := 0
	for index, player := range gs.Players {
		if player == gs.ActivePlayer {
			apIndex = index
			break
		}
	}

	apIndex++
	apIndex %= len(gs.Players)
	for !gs.Players[apIndex].IsAlive(gs) {
		apIndex++
		apIndex %= len(gs.Players)
	}

	gs.ActivePlayer = gs.Players[apIndex]

	return gs
}

func (gs *Gamestate) Winner() (bool, *Player) {
	if gs.HasEnded() {
		var winner *Player
		for _, player := range gs.Players {
			if player.IsAlive(gs) {
				winner = player
				break
			}
		}

		return true, winner
	}

	return false, nil
}

func (gs *Gamestate) HasEnded() bool {
	livePlayerCount := 0

	for _, player := range gs.Players {
		if player.IsAlive(gs) {
			livePlayerCount++
		}
	}

	return livePlayerCount < 2
}
