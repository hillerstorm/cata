package monk

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

/*
// 116645 - Teachings of the Monastery (Mistweaver)
Tooltip:
Attack with the palm of your hand, dealing $?s116645[${6*$<low>} to ${6*$<high>}][${3*$<low>} to ${3*$<high>}] damage.

Also grants you Tiger Power, causing your attacks to ignore $125359m1% of enemies' armor for $125359d.
*/
func (monk *Monk) registerTigerPalm() {
	actionID := core.ActionID{SpellID: 100787}
	chiMetrics := monk.NewChiMetrics(actionID)

	tigerPowerDebuff := monk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Tiger Power" + target.Label,
			ActionID: core.ActionID{SpellID: 125359},
			Duration: time.Second * 20,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				monk.AttackTables[aura.Unit.UnitIndex].ArmorIgnoreFactor += 0.3
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				monk.AttackTables[aura.Unit.UnitIndex].ArmorIgnoreFactor -= 0.3
			},
		})
	})

	monk.TigerPalm = monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellTigerPalm,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 3.0,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultMeleeCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return monk.ComboPoints() >= 1 || monk.ComboBreakerTigerPalmAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				if monk.ComboBreakerTigerPalmAura.IsActive() {
					monk.SpendChi(sim, 0, chiMetrics)
				} else {
					monk.SpendChi(sim, 1, chiMetrics)
				}

				tigerPowerDebuff.Get(result.Target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
