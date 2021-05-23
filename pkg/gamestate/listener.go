package gamestate

type Listener interface {
	Notify(Action, *Gamestate)
	Subscribe(*Gamestate)
	Unsubscribe(*Gamestate)
}

type ExecuteOnceListener struct {
	validate func(Action, *Gamestate) bool
	execute  func(Listener, Action, *Gamestate)
}

func NewExecuteOnceListener(val func(Action, *Gamestate) bool, ex func(Listener, Action, *Gamestate)) *ExecuteOnceListener {
	return &ExecuteOnceListener{
		validate: val,
		execute:  ex,
	}
}

func (eol *ExecuteOnceListener) Subscribe(game *Gamestate) {
	game.Subscribe(eol)
}

func (eol *ExecuteOnceListener) Unsubscribe(game *Gamestate) {
	game.Unsubscribe(eol)
}

func (eol *ExecuteOnceListener) Notify(action Action, game *Gamestate) {
	if eol.validate(action, game) {
		eol.execute(eol, action, game)
		game.Unsubscribe(eol)
	}
}

type UntilEndOfTurnListener struct {
	validate func(Action, *Gamestate) bool
	execute  func(Listener, Action, *Gamestate)
}

func NewUntilEndOfTurnListener(validate func(Action, *Gamestate) bool, execute func(Listener, Action, *Gamestate)) *UntilEndOfTurnListener {
	return &UntilEndOfTurnListener{
		validate: validate,
		execute:  execute,
	}
}

func (eot *UntilEndOfTurnListener) Subscribe(game *Gamestate) {
	game.Subscribe(eot)
}

func (eot *UntilEndOfTurnListener) Unsubscribe(game *Gamestate) {
	game.Unsubscribe(eot)
}

func (eot *UntilEndOfTurnListener) Notify(action Action, game *Gamestate) {
	if eot.validate(action, game) {
		eot.execute(eot, action, game)
	}

	_, ok := action.(*EndTurnAction)
	if ok {
		game.Unsubscribe(eot)
	}
}
