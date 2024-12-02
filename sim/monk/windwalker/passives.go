package windwalker

import (
	"time"

	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/monk"
)

func (ww *WindwalkerMonk) registerPassives() {
	ww.registerCombatConditioning()
	ww.registerComboBreaker()
	ww.registerTigerStrikes()
}

func (ww *WindwalkerMonk) registerCombatConditioning() {
	if !ww.HasMinorGlyph(proto.MonkMinorGlyph_MonkMinorGlyphBlackoutKick) && ww.PseudoStats.InFrontOfTarget {
		return
	}

	// TODO: This should be able to crit...
	// TODO: The ignite effect should also tick every second instead of every 2 seconds
	cata.RegisterIgniteEffect(&ww.Unit, cata.IgniteConfig{
		ActionID:           core.ActionID{SpellID: 100784}.WithTag(2), // actual 128531
		DotAuraLabel:       "Blackout Kick (DoT)" + ww.Label,
		DisableCastMetrics: true,
		IncludeAuraDelay:   true,
		SpellSchool:        core.SpellSchoolPhysical,

		ProcTrigger: core.ProcTrigger{
			Name:           "Combat Conditioning" + ww.Label,
			Callback:       core.CallbackOnSpellHitDealt,
			ClassSpellMask: monk.MonkSpellBlackoutKick,
			Outcome:        core.OutcomeLanded,
		},

		DamageCalculator: func(result *core.SpellResult) float64 {
			return result.Damage * 0.2
		},
	})
}

func (ww *WindwalkerMonk) registerComboBreaker() {
	ww.ComboBreakerBlackoutKickAura = ww.RegisterAura(core.Aura{
		Label:    "Combo Breaker: Blackout Kick" + ww.Label,
		ActionID: core.ActionID{SpellID: 116768},
		Duration: time.Second * 20,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&monk.MonkSpellBlackoutKick == 0 || !result.Landed() {
				return
			}

			ww.ComboBreakerBlackoutKickAura.Deactivate(sim)
		},
	})

	ww.ComboBreakerTigerPalmAura = ww.RegisterAura(core.Aura{
		Label:    "Combo Breaker: Tiger Palm" + ww.Label,
		ActionID: core.ActionID{SpellID: 118864},
		Duration: time.Second * 20,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&monk.MonkSpellTigerPalm == 0 || !result.Landed() {
				return
			}

			ww.ComboBreakerTigerPalmAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&ww.Unit, core.ProcTrigger{
		Name:           "Combo Breaker: Blackout Kick Trigger" + ww.Label,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: monk.MonkSpellJab,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.12,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ww.ComboBreakerBlackoutKickAura.Activate(sim)
		},
	})

	core.MakeProcTriggerAura(&ww.Unit, core.ProcTrigger{
		Name:           "Combo Breaker: Tiger Palm Trigger" + ww.Label,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: monk.MonkSpellJab,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.12,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ww.ComboBreakerTigerPalmAura.Activate(sim)
		},
	})
}

func (ww *WindwalkerMonk) registerTigerStrikes() {
	tigerStrikesMHID := core.ActionID{SpellID: 120274}
	tigerStrikesOHID := core.ActionID{SpellID: 120278}

	var tigerStrikesMHSpell *core.Spell
	var tigerStrikesOHSpell *core.Spell
	tigerStrikesBuff := ww.RegisterAura(core.Aura{
		Label:     "Tiger Strikes" + ww.Label,
		ActionID:  core.ActionID{SpellID: 120273},
		Duration:  time.Second * 15,
		MaxStacks: 4,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			mhConfig := *ww.AutoAttacks.MHConfig()
			mhConfig.ActionID = tigerStrikesMHID
			mhConfig.ClassSpellMask = monk.MonkSpellTigerStrikes
			mhConfig.Flags |= core.SpellFlagPassiveSpell
			mhConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhiteNoGlance)
			}
			tigerStrikesMHSpell = ww.GetOrRegisterSpell(mhConfig)

			if ww.HasOHWeapon() {
				ohConfig := *ww.AutoAttacks.OHConfig()
				ohConfig.ActionID = tigerStrikesOHID
				ohConfig.ClassSpellMask = monk.MonkSpellTigerStrikes
				ohConfig.Flags |= core.SpellFlagPassiveSpell
				ohConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhiteNoGlance)
				}
				tigerStrikesOHSpell = ww.GetOrRegisterSpell(ohConfig)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ww.MultiplyMeleeSpeed(sim, 1.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ww.MultiplyMeleeSpeed(sim, 1/1.5)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !aura.IsActive() || !result.Landed() || spell.ClassSpellMask&monk.MonkSpellTigerStrikes != 0 {
				return
			}

			var spellToCast *core.Spell
			if spell == ww.AutoAttacks.MHAuto() {
				spellToCast = tigerStrikesMHSpell
			} else if spell == ww.AutoAttacks.OHAuto() {
				spellToCast = tigerStrikesOHSpell
			}

			if spellToCast != nil {
				aura.RemoveStack(sim)
				delaySeconds := sim.RollWithLabel(0.8, 1.2, "Tiger Strikes Delay")
				sim.AddPendingAction(&core.PendingAction{
					NextActionAt: sim.CurrentTime + core.DurationFromSeconds(delaySeconds),
					Priority:     core.ActionPriorityAuto,
					OnAction: func(sim *core.Simulation) {
						spellToCast.Cast(sim, result.Target)
					},
				})
			}
		},
	})

	core.MakeProcTriggerAura(&ww.Unit, core.ProcTrigger{
		Name:       "Tiger Strikes Buff Trigger" + ww.Label,
		ActionID:   core.ActionID{SpellID: 120272},
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskWhiteHit,
		ProcChance: 1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&monk.MonkSpellTigerStrikes != 0 {
				return
			}

			if sim.Proc(0.08, "Tiger Strikes") {
				tigerStrikesBuff.Activate(sim)
				tigerStrikesBuff.SetStacks(sim, 4)
			}
		},
	})
}
