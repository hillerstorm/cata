package monk

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

/*
Tooltip:
You Jab the target, dealing ${1.5*$<low>} to ${1.5*$<high>} damage and generating $s2 Chi.
*/
func (monk *Monk) registerJab() {
	monk.Jab = monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 100780},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellJab,
		MaxRange:       core.MaxMeleeRange,

		EnergyCost: core.EnergyCostOptions{
			Cost:          core.TernaryFloat64(monk.StanceMatches(WiseSerpent), 0, 40),
			Refund:        0.8,
			RefundMetrics: monk.EnergyRefundMetrics,
		},
		ManaCost: core.ManaCostOptions{
			BaseCost: core.TernaryFloat64(monk.StanceMatches(WiseSerpent), 0.08, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.5,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				chi := 1

				if monk.StanceMatches(FierceTiger) {
					chi += 1
				}

				monk.AddChi(sim, int32(chi), spell.ComboPointMetrics())
			}
		},
	})
}
