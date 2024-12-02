package monk

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

/*
Tooltip:
Channels Jade lightning at the target, causing ${6*($m1+*$ap*0.386)} Nature damage over 6 sec. When dealing damage, you have a 30% chance to generate 1 Chi.

If the enemy attacks you within melee range while victim to Crackling Jade Lightning, they are knocked back a short distance. This effect has an 8 sec cooldown.

TODO: Check if it does a one-time hit check or per tick
TODO: Spell or melee hit / crit
TODO: Courageous Primal Diamond should make all ticks ignore mana cost
*/
func (monk *Monk) registerCracklingJadeLightning() {
	actionID := core.ActionID{SpellID: 117952}
	energyMetrics := monk.NewEnergyMetrics(actionID)
	manaMetrics := monk.NewManaMetrics(actionID)
	chiMetrics := monk.NewChiMetrics(core.ActionID{SpellID: 123333})
	avgScaling := monk.ClassSpellScaling * 0.1800000072

	canTick := func(sim *core.Simulation, spell *core.Spell) (bool, func()) {
		isWiseSerpent := monk.StanceMatches(WiseSerpent)
		currentResource := core.TernaryFloat64(isWiseSerpent, monk.CurrentMana(), monk.CurrentEnergy())
		baseCost := core.TernaryFloat64(isWiseSerpent, 0.0157*monk.BaseMana, 20.0)
		cost := spell.ApplyCostModifiers(baseCost)

		if currentResource >= cost {
			return true, func() {
				if isWiseSerpent {
					monk.SpendMana(sim, cost, manaMetrics)
				} else {
					monk.SpendEnergy(sim, cost, energyMetrics)
				}
			}
		} else {
			return false, nil
		}
	}

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolNature,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MonkSpellCracklingJadeLightning,
		MaxRange:       40,

		ManaCost: core.ManaCostOptions{
			BaseCost: core.TernaryFloat64(monk.StanceMatches(WiseSerpent), 0.0157, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Crackling Jade Lightning" + monk.Label,
			},
			NumberOfTicks:       6,
			TickLength:          time.Second,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if ok, spendCost := canTick(sim, dot.Spell); ok {
					spendCost()
					baseDamage := avgScaling + dot.Spell.MeleeAttackPower()*0.386
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicCrit)

					if sim.Proc(0.3, "Crackling Jade Lightning") {
						monk.AddChi(sim, dot.Spell, 1, chiMetrics)
					}
				} else {
					monk.AutoAttacks.EnableMeleeSwing(sim)
					monk.ExtendGCDUntil(sim, sim.CurrentTime+monk.ChannelClipDelay)

					// Deactivating within OnTick causes a panic since tickAction gets set to nil in the default OnExpire
					sim.AddPendingAction(&core.PendingAction{
						NextActionAt: sim.CurrentTime,
						OnAction: func(sim *core.Simulation) {
							dot.Deactivate(sim)
						},
					})
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultSpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)

			if result.Landed() {
				dot := spell.Dot(target)
				dot.Apply(sim)
				expiresAt := dot.ExpiresAt()
				monk.AutoAttacks.StopMeleeUntil(sim, expiresAt, false)
				monk.ExtendGCDUntil(sim, expiresAt+monk.ReactionTime)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
