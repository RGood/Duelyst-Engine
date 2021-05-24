package game

import "github.com/RGood/game_engine/pkg/gamestate"

type Artifact struct {
	Name      string
	Cost      int
	Charges   int
	Owner     *Player
	intercept func(*Artifact, gamestate.Action, *gamestate.Gamestate) gamestate.Action
	notify    func(*Artifact, gamestate.Action, *gamestate.Gamestate)
	onEquip   func(*Artifact, *gamestate.Gamestate)
	onUnequip func(*Artifact, *gamestate.Gamestate)
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

func (artifact *Artifact) OnEquip(effect func(*Artifact, *gamestate.Gamestate)) *Artifact {
	artifact.onEquip = effect

	return artifact
}

func (artifact *Artifact) OnUnEquip(effect func(*Artifact, *gamestate.Gamestate)) *Artifact {
	artifact.onUnequip = effect

	return artifact
}

func (artifact *Artifact) OnNotify(effect func(*Artifact, gamestate.Action, *gamestate.Gamestate)) *Artifact {
	artifact.notify = effect

	return artifact
}

func (artifact *Artifact) OnIntercept(effect func(*Artifact, gamestate.Action, *gamestate.Gamestate) gamestate.Action) *Artifact {
	artifact.intercept = effect

	return artifact
}

func (artifact *Artifact) Equip(owner *Player, gamestate *gamestate.Gamestate) {
	artifact.Owner = owner
	artifact.AddIntercept(gamestate)
	artifact.Subscribe(gamestate)

	if artifact.onEquip != nil {
		artifact.onEquip(artifact, gamestate)
	}
}

func (artifact *Artifact) Remove(gamestate *gamestate.Gamestate) {
	if artifact.onUnequip != nil {
		artifact.onUnequip(artifact, gamestate)
	}

	artifact.RemoveIntercept(gamestate)
	artifact.Unsubscribe(gamestate)
	artifact.Owner = nil
}

func (artifact *Artifact) Subscribe(gamestate *gamestate.Gamestate) {
	gamestate.Subscribe(artifact)
}

func (artifact *Artifact) Unsubscribe(gamestate *gamestate.Gamestate) {
	gamestate.Unsubscribe(artifact)
}

func (artifact *Artifact) AddIntercept(gamestate *gamestate.Gamestate) {
	gamestate.AddInterceptor(artifact)
}

func (artifact *Artifact) RemoveIntercept(gamestate *gamestate.Gamestate) {
	gamestate.RemoveInterceptor(artifact)
}

func (artifact *Artifact) Apply(action gamestate.Action, gamestate *gamestate.Gamestate) gamestate.Action {
	if artifact.intercept != nil {
		return artifact.intercept(artifact, action, gamestate)
	}

	return action
}

func (artifact *Artifact) Notify(action gamestate.Action, gamestate *gamestate.Gamestate) {
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
