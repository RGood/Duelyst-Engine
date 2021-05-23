package gamestate

type Artifact struct {
	Name      string
	Cost      int
	Charges   int
	Owner     *Player
	intercept func(*Artifact, Action, *Gamestate) Action
	notify    func(*Artifact, Action, *Gamestate)
	onEquip   func(*Artifact, *Gamestate)
	onUnequip func(*Artifact, *Gamestate)
}

func NewArtifact(name string, cost int) *Artifact {
	return &Artifact{
		Name:      name,
		Cost:      cost,
		Charges:   3,
		Owner:     nil,
		intercept: nil,
		notify:    nil,
		onEquip:   nil,
		onUnequip: nil,
	}
}

func (artifact *Artifact) OnEquip(effect func(*Artifact, *Gamestate)) *Artifact {
	artifact.onEquip = effect

	return artifact
}

func (artifact *Artifact) OnUnEquip(effect func(*Artifact, *Gamestate)) *Artifact {
	artifact.onUnequip = effect

	return artifact
}

func (artifact *Artifact) OnNotify(effect func(*Artifact, Action, *Gamestate)) *Artifact {
	artifact.notify = effect

	return artifact
}

func (artifact *Artifact) OnIntercept(effect func(*Artifact, Action, *Gamestate) Action) *Artifact {
	artifact.intercept = effect

	return artifact
}

func (artifact *Artifact) Equip(owner *Player, gamestate *Gamestate) {
	artifact.Owner = owner
	artifact.AddIntercept(gamestate)
	artifact.Subscribe(gamestate)

	if artifact.onEquip != nil {
		artifact.onEquip(artifact, gamestate)
	}
}

func (artifact *Artifact) Remove(gamestate *Gamestate) {
	if artifact.onUnequip != nil {
		artifact.onUnequip(artifact, gamestate)
	}

	artifact.RemoveIntercept(gamestate)
	artifact.Unsubscribe(gamestate)
	artifact.Owner = nil
}

func (artifact *Artifact) Subscribe(gamestate *Gamestate) {
	gamestate.Subscribe(artifact)
}

func (artifact *Artifact) Unsubscribe(gamestate *Gamestate) {
	gamestate.Unsubscribe(artifact)
}

func (artifact *Artifact) AddIntercept(gamestate *Gamestate) {
	gamestate.AddInterceptor(artifact)
}

func (artifact *Artifact) RemoveIntercept(gamestate *Gamestate) {
	gamestate.RemoveInterceptor(artifact)
}

func (artifact *Artifact) Apply(action Action, gamestate *Gamestate) Action {
	if artifact.intercept != nil {
		return artifact.intercept(artifact, action, gamestate)
	}

	return action
}

func (artifact *Artifact) Notify(action Action, gamestate *Gamestate) {
	if artifact.notify != nil {
		artifact.notify(artifact, action, gamestate)
	}

	damageAction, ok := action.(*DamageAction)

	if ok && damageAction.Damage > 0 && damageAction.Unit.GetOwner() == artifact.Owner && damageAction.Unit.GetType() == "general" {
		artifact.Charges--

		if artifact.Charges == 0 {
			gamestate.MakeMove(&RemoveArtifactAction{
				Artifact: artifact,
			})
		}
	}

}
