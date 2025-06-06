package mop

import (
	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
)

func RegisterAllEnchants() {

	// Enchants

	// Permanently enchant a melee weapon to sometimes increase haste by 450 for 12s when healing or dealing
	// spell or melee damage.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Enchant Weapon - Hurricane",
		EnchantID: 4083,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Permanently enchant a weapon to sometimes increase Spirit by 200 for 15s when healing or dealing damage
	// with spells.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Enchant Weapon - Heartsong",
		EnchantID: 4084,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Permanently enchant a weapon to sometimes increase Intellect by 500 for 12s when dealing damage or healing
	// with spells.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Enchant Weapon - Power Torrent",
		EnchantID: 4097,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Permanently enchant a weapon to sometimes increase dodge by 600 and movement speed by 15% for 10s when
	// striking in melee, stacking with passive movement speed effects.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Enchant Weapon - Windwalk",
		EnchantID: 4098,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Permanently enchant a weapon to sometimes increase attack power by 1000 for 12s when striking in melee.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Enchant Weapon - Landslide",
		EnchantID: 4099,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// |CFF00F0FFThis recipe automatically improves when you reach 550 skill in Tailoring.|R
	//
	// Embroiders a subtle pattern of light into your cloak, giving you a chance to increase your Intellect by
	// 580 for 15s when casting a spell.
	//
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Lightweave Embroidery (Rank 2)",
		EnchantID: 4115,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellDamageProc,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// |CFF00F0FFThis recipe automatically improves when you reach 550 skill in Tailoring.|R
	//
	// Embroiders a magical pattern into your cloak, giving you a chance to increase your Spirit by 580 for 15s
	// when you cast a spell.
	//
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Darkglow Embroidery (Rank 2)",
		EnchantID: 4116,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// |CFF00F0FFThis recipe automatically improves when you reach 550 skill in Tailoring.|R
	//
	// Embroiders a magical pattern into your cloak, causing your damaging melee and ranged attacks to sometimes
	// increase your attack power by 1000 for 15s.
	//
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Swordguard Embroidery (Rank 2)",
		EnchantID: 4118,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Attaches a permanent scope to a bow or gun that sometimes increases ranged attack power by 800 for 10s
	// when dealing damage with ranged attacks.
	//
	// Attaching this scope to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Gnomish X-Ray Scope",
		EnchantID: 4175,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Attaches a permanent device to a bow or gun that sometimes launches a rabid critter, dealing additional
	// damage to the target and increasing your Agility by 300 for 10s.
	//
	// Attaching this device to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Flintlocke's Woodchucker",
		EnchantID: 4267,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Permanently attaches Lord Blastington's special scope to a ranged weapon, sometimes increasing Agility
	// by 1800 for 10s when dealing damage with ranged attacks.
	//
	// Attaching this scope to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Lord Blastington's Scope of Doom",
		EnchantID: 4699,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Permanently attaches a mirrored scope to a ranged weapon, sometimes increases critical strike by 900 for
	// 10s when dealing damage with ranged attacks.
	//
	// Attaching this scope to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Mirror Scope",
		EnchantID: 4700,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Embroiders a subtle pattern of light into your cloak, giving you a chance to increase your Intellect by
	// 2000 for 15s when casting a spell.
	//
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Lightweave Embroidery (Rank 3)",
		EnchantID: 4892,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellDamageProc,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Embroiders a magical pattern into your cloak, giving you a chance to increase your Spirit by 3000 for
	// 15s when you cast a spell.
	//
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Darkglow Embroidery (Rank 3)",
		EnchantID: 4893,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})

	// Embroiders a magical pattern into your cloak, causing your damaging melee and ranged attacks to sometimes
	// increase your attack power by 4000 for 15s.
	//
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Swordguard Embroidery (Rank 3)",
		EnchantID: 4894,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
}
