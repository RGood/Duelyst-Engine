package gamestate

type Interceptor interface {
	Apply(Action, *Gamestate) Action
	Subscribe(*Gamestate)
	Unsubscribe(*Gamestate)
}

type UntilEndOfTurnInterceptor struct {
	validate func(Action, *Gamestate) bool
	execute  func(Interceptor, Action, *Gamestate) Action
}

func NewUntilEndOfTurnInterceptor(validate func(Action, *Gamestate) bool, execute func(Interceptor, Action, *Gamestate) Action) *UntilEndOfTurnInterceptor {
	return &UntilEndOfTurnInterceptor{
		validate: validate,
		execute:  execute,
	}
}

func (eot *UntilEndOfTurnInterceptor) Subscribe(game *Gamestate) {
	game.AddInterceptor(eot)

	NewUntilEndOfTurnListener(func(action Action, gamestate *Gamestate) bool {
		_, ok := action.(*EndTurnAction)
		return ok
	}, func(listener Listener, action Action, gamestate *Gamestate) {
		eot.Unsubscribe(gamestate)
	}).Subscribe(game)
}

func (eot *UntilEndOfTurnInterceptor) Unsubscribe(game *Gamestate) {
	game.RemoveInterceptor(eot)
}

func (eot *UntilEndOfTurnInterceptor) Apply(action Action, game *Gamestate) Action {
	if eot.validate(action, game) {
		return eot.execute(eot, action, game)
	}

	return action
}
