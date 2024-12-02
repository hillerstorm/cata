package monk

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (monk *Monk) ApplyTalents() {
	// Level 45
	monk.registerAscension()
	monk.registerChiBrew()
}

func (monk *Monk) registerAscension() {
	if !monk.Talents.Ascension {
		return
	}

	core.MakePermanent(monk.GetOrRegisterAura(core.Aura{
		Label:    "Ascension" + monk.Label,
		ActionID: core.ActionID{SpellID: 115396},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: Increase max mana by 15% when in Wise Serpent stance
			monk.MultiplyEnergyRegenSpeed(sim, 1.15)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			monk.MultiplyEnergyRegenSpeed(sim, 1.0/1.15)
		},
	}))
}

func (monk *Monk) registerChiBrew() {
	if !monk.Talents.ChiBrew {
		return
	}

	actionID := core.ActionID{SpellID: 115399}
	chiMetrics := monk.NewChiMetrics(actionID)

	var chiBrewAura *core.Aura
	chiBrewAura = monk.RegisterAura(core.Aura{
		Label:     "Chi Brew" + monk.Label,
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 2,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if monk.chiBrewRecharge != nil {
				monk.chiBrewRecharge.Cancel(sim)
			}

			aura.Activate(sim)
			aura.SetStacks(sim, 2)
		},

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if !aura.IsActive() {
				return
			}

			if newStacks < oldStacks {
				nextRecharge := &core.PendingAction{
					NextActionAt: sim.CurrentTime + time.Second*45,
					OnAction: func(sim *core.Simulation) {
						aura.Activate(sim)
						aura.AddStack(sim)
					},
				}

				if monk.chiBrewRecharge != nil {
					// If we have an existing stack recharging, set this new one as current when it's done.
					// This way we can always check next recharge time from the APL.
					oldAction := monk.chiBrewRecharge.OnAction
					monk.chiBrewRecharge.OnAction = func(sim *core.Simulation) {
						monk.chiBrewRecharge = nextRecharge
						oldAction(sim)
					}
				} else {
					monk.chiBrewRecharge = nextRecharge
				}

				sim.AddPendingAction(nextRecharge)
			} else if newStacks > oldStacks {
				if newStacks == 2 || monk.chiBrewRecharge.IsConsumed() {
					monk.chiBrewRecharge = nil
				}
			}
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellChiBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return chiBrewAura.GetStacks() >= 1
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: Add 5 Elusive Brew stacks for Brewmasters
			// TODO: Add 2 Mana Tea stacks for Mistweavers

			monk.AddChi(sim, 2, chiMetrics)
			monk.AddBrewStacks(sim, 2)

			chiBrewAura.RemoveStack(sim)
		},
	})
}
