package monk

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	SpellFlagBuilder = core.SpellFlagAgentReserved2
	SpellFlagSpender = core.SpellFlagAgentReserved3
)

type OnChiSpent func(sim *core.Simulation, chiSpent int32)
type OnNewBrewStacks func(sim *core.Simulation, stacksToAdd int32)

type Monk struct {
	core.Character

	ClassSpellScaling float64

	Talents           *proto.MonkTalents
	Options           *proto.MonkOptions
	BrewmasterOptions *proto.BrewmasterMonk_Options
	MistweaverOptions *proto.MistweaverMonk_Options
	WindwalkerOptions *proto.WindwalkerMonk_Options

	Stance Stance

	StanceOfTheFierceTigerAura *core.Aura
	StanceOfTheSturdyOxAura    *core.Aura
	StanceOfTheWiseSerpentAura *core.Aura
	StanceOfTheFierceTiger     *core.Spell
	StanceOfTheSturdyOx        *core.Spell
	StanceOfTheWiseSerpent     *core.Spell

	ComboBreakerBlackoutKickAura *core.Aura
	ComboBreakerTigerPalmAura    *core.Aura

	Jab          *core.Spell
	TigerPalm    *core.Spell
	BlackoutKick *core.Spell

	onChiSpent      OnChiSpent
	onNewBrewStacks OnNewBrewStacks
	chiBrewRecharge *core.PendingAction
}

func (monk *Monk) SpendChi(sim *core.Simulation, chiToSpend int32, metrics *core.ResourceMetrics) {
	monk.SpendPartialComboPoints(sim, chiToSpend, metrics)
	if monk.onChiSpent != nil {
		monk.onChiSpent(sim, chiToSpend)
	}
}

func (monk *Monk) RegisterOnChiSpent(onChiSpent func(*core.Simulation, int32)) {
	monk.onChiSpent = onChiSpent
}

func (monk *Monk) AddBrewStacks(sim *core.Simulation, stacksToAdd int32) {
	if monk.onNewBrewStacks != nil {
		monk.onNewBrewStacks(sim, stacksToAdd)
	}
}
func (monk *Monk) RegisterOnNewBrewStacks(onNewBrewStacks func(*core.Simulation, int32)) {
	monk.onNewBrewStacks = onNewBrewStacks
}

func (monk *Monk) GetCharacter() *core.Character {
	return &monk.Character
}

func (monk *Monk) GetMonk() *Monk {
	return monk
}

func (monk *Monk) AddRaidBuffs(_ *proto.RaidBuffs)   {}
func (monk *Monk) AddPartyBuffs(_ *proto.PartyBuffs) {}

func (monk *Monk) HasPrimeGlyph(glyph proto.MonkPrimeGlyph) bool {
	return monk.HasGlyph(int32(glyph))
}
func (monk *Monk) HasMajorGlyph(glyph proto.MonkMajorGlyph) bool {
	return monk.HasGlyph(int32(glyph))
}
func (monk *Monk) HasMinorGlyph(glyph proto.MonkMinorGlyph) bool {
	return monk.HasGlyph(int32(glyph))
}

func (monk *Monk) Initialize() {
	monk.AutoAttacks.MHConfig().CritMultiplier = monk.MeleeCritMultiplier()
	monk.AutoAttacks.OHConfig().CritMultiplier = monk.MeleeCritMultiplier()

	monk.registerStances()
	monk.applyGlyphs()
	monk.registerSpells()
	monk.registerWayOfTheMonk()
}

func (monk *Monk) registerSpells() {
	monk.registerJab()
	monk.registerTigerPalm()
	monk.registerBlackoutKick()
}

func (monk *Monk) Reset(sim *core.Simulation) {
	switch monk.Stance {
	case SturdyOx:
		monk.StanceOfTheSturdyOxAura.Activate(sim)
	case WiseSerpent:
		monk.StanceOfTheWiseSerpentAura.Activate(sim)
	case FierceTiger:
		monk.StanceOfTheFierceTigerAura.Activate(sim)
	}
}

func (monk *Monk) MeleeCritMultiplier() float64 {
	return monk.Character.MeleeCritMultiplier(1, 0)
}
func (monk *Monk) SpellCritMultiplier() float64 {
	return monk.Character.SpellCritMultiplier(1, 0)
}

func NewMonk(character *core.Character, options *proto.MonkOptions, talents string) *Monk {
	monk := &Monk{
		Character:         *character,
		Talents:           &proto.MonkTalents{},
		Options:           options,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassMonk),
	}

	core.FillTalentsProto(monk.Talents.ProtoReflect(), talents, [3]int{5, 0, 0})

	monk.PseudoStats.CanParry = true

	maxChi := 4

	if monk.Talents.Ascension {
		maxChi += 1
	}

	monk.EnableEnergyBar(100, int32(maxChi), proto.Class_ClassMonk)

	monk.EnableAutoAttacks(monk, core.AutoAttackOptions{
		MainHand:       monk.WeaponFromMainHand(0),
		OffHand:        monk.WeaponFromOffHand(0),
		AutoSwingMelee: true,
	})

	monk.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	monk.AddStatDependency(stats.Agility, stats.AttackPower, 2)
	monk.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])

	return monk
}

type MonkAgent interface {
	GetMonk() *Monk
}

const (
	MonkSpellFlagNone int64 = 0
	MonkSpellJab      int64 = 1 << iota
	MonkSpellTigerPalm
	MonkSpellBlackoutKick

	// Talents
	MonkSpellChiBrew

	// Windwalker
	MonkSpellTigereyeBrew
	MonkSpellRisingSunKick
	MonkSpellEnergizingBrew
	MonkSpellTigerStrikes

	MonkSpellLast
	MonkSpellsAll = MonkSpellLast<<1 - 1
)
