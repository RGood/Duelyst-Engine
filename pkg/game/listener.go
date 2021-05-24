package game

import "github.com/RGood/game_engine/pkg/gamestate"

type ExecuteOnceListener struct {
	validate func(gamestate.Action, *gamestate.Gamestate) bool
	execute  func(gamestate.Listener, gamestate.Action, *gamestate.Gamestate)
}

func NewExecuteOnceListener(val func(gamestate.Action, *gamestate.Gamestate) bool, ex func(gamestate.Listener, gamestate.Action, *gamestate.Gamestate)) *ExecuteOnceListener {
	return &ExecuteOnceListener{
		validate: val,
		execute:  ex,
	}
}

func (eol *ExecuteOnceListener) Subscribe(game *gamestate.Gamestate) {
	game.Subscribe(eol)
}

func (eol *ExecuteOnceListener) Unsubscribe(game *gamestate.Gamestate) {
	game.Unsubscribe(eol)
}

func (eol *ExecuteOnceListener) Notify(action gamestate.Action, game *gamestate.Gamestate) {
	if eol.validate(action, game) {
		eol.execute(eol, action, game)
		game.Unsubscribe(eol)
	}
}

type UntilEndOfTurnListener struct {
	validate func(gamestate.Action, *gamestate.Gamestate) bool
	execute  func(gamestate.Listener, gamestate.Action, *gamestate.Gamestate)
}

func NewUntilEndOfTurnListener(validate func(gamestate.Action, *gamestate.Gamestate) bool, execute func(gamestate.Listener, gamestate.Action, *gamestate.Gamestate)) *UntilEndOfTurnListener {
	return &UntilEndOfTurnListener{
		validate: validate,
		execute:  execute,
	}
}

func (eot *UntilEndOfTurnListener) Subscribe(game *gamestate.Gamestate) {
	game.Subscribe(eot)
}

func (eot *UntilEndOfTurnListener) Unsubscribe(game *gamestate.Gamestate) {
	game.Unsubscribe(eot)
}

func (eot *UntilEndOfTurnListener) Notify(action gamestate.Action, game *gamestate.Gamestate) {
	if eot.validate(action, game) {
		eot.execute(eot, action, game)
	}

	_, ok := action.(*EndTurnAction)
	if ok {
		game.Unsubscribe(eot)
	}
}
