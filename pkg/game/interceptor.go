package game

import "github.com/RGood/game_engine/pkg/gamestate"

type UntilEndOfTurnInterceptor struct {
	validate func(gamestate.Action, *gamestate.Gamestate) bool
	execute  func(gamestate.Interceptor, gamestate.Action, *gamestate.Gamestate) gamestate.Action
}

func NewUntilEndOfTurnInterceptor(validate func(gamestate.Action, *gamestate.Gamestate) bool, execute func(gamestate.Interceptor, gamestate.Action, *gamestate.Gamestate) gamestate.Action) *UntilEndOfTurnInterceptor {
	return &UntilEndOfTurnInterceptor{
		validate: validate,
		execute:  execute,
	}
}

func (eot *UntilEndOfTurnInterceptor) Subscribe(game *gamestate.Gamestate) {
	game.AddInterceptor(eot)

	NewUntilEndOfTurnListener(func(action gamestate.Action, gamestate *gamestate.Gamestate) bool {
		_, ok := action.(*EndTurnAction)
		return ok
	}, func(listener gamestate.Listener, action gamestate.Action, gamestate *gamestate.Gamestate) {
		eot.Unsubscribe(gamestate)
	}).Subscribe(game)
}

func (eot *UntilEndOfTurnInterceptor) Unsubscribe(game *gamestate.Gamestate) {
	game.RemoveInterceptor(eot)
}

func (eot *UntilEndOfTurnInterceptor) Apply(action gamestate.Action, game *gamestate.Gamestate) gamestate.Action {
	if eot.validate(action, game) {
		return eot.execute(eot, action, game)
	}

	return action
}
