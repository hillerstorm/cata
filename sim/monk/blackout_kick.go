package monk

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

/*
// 116645 - Teachings of the Monastery (Mistweaver)
// 115070 - Stance of the Wise Serpent (Mistweaver)
// 128595 - Combat Conditioning (Windwalker)
// 117967 - Brewmaster Training (Brewmaster)
// 115307 - Shuffle (Brewmaster)
// 127722 - Serpent's Zeal (Mistweaver)
Tooltip:
Kick with a blast of Chi energy, dealing ${7.12*$<low>} to ${7.12*$<high>} Physical damage

// Mistweaver in Wise Serpent stance
$?s116645&a115070[ to your target and ${3.56*$<low>} to ${3.56*$<high>} to up to $116645s4 additional nearby targets][].

// Windwalker
$?s128595[ If behind the target, you deal an additional $m2% damage over $128531d. If in front of the target, you are instantly healed for $m2% of the damage done.][]

// Brewmaster
$?s117967[Also causes you to gain Shuffle, increasing your parry chance by $115307s1% and your Stagger amount by an additional $115307s2% for $115307d.][]

// Mistweaver
$?s116645[Also empowers you with Serpent's Zeal, causing you and your summoned Jade Serpent Statue to heal nearby injured targets equal to $127722m1% of your auto-attack damage.][]
*/
func (monk *Monk) registerBlackoutKick() {
	actionID := core.ActionID{SpellID: 100784}.WithTag(1)
	chiMetrics := monk.NewChiMetrics(actionID)

	monk.BlackoutKick = monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellBlackoutKick,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 7.12,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultMeleeCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return monk.ComboPoints() >= 2 || monk.ComboBreakerBlackoutKickAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() && !monk.ComboBreakerBlackoutKickAura.IsActive() {
				monk.SpendChi(sim, 2, chiMetrics)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
