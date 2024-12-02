package monk

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

/*
Tooltip:
The Monk attunes $G himself:herself; differently depending on the weapon type.

One-handed weapons / Dual-wield one-handed weapons:
Autoattack damage increased by 40%.

Two-handed weapons:
Melee attack speed increased by 40%.
*/
func (monk *Monk) registerWayOfTheMonk() {
	mh := monk.GetMHWeapon()
	auraConfig := core.Aura{
		Label:    "Way of the Monk" + monk.Label,
		ActionID: core.ActionID{SpellID: 120277},
	}

	if mh != nil && (mh.WeaponType == proto.WeaponType_WeaponTypeStaff || mh.WeaponType == proto.WeaponType_WeaponTypePolearm) {
		auraConfig.OnGain = func(aura *core.Aura, sim *core.Simulation) {
			monk.MultiplyMeleeSpeed(sim, 1.4)
		}
		auraConfig.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
			monk.MultiplyMeleeSpeed(sim, 1/1.4)
		}
	} else {
		monk.AutoAttacks.MHConfig().DamageMultiplier *= 1.4
		monk.AutoAttacks.OHConfig().DamageMultiplier *= 1.4
	}

	core.MakePermanent(monk.RegisterAura(auraConfig))
}
