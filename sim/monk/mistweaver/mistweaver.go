package mistweaver

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/monk"
)

func RegisterMistweaverMonk() {
	core.RegisterAgentFactory(
		proto.Player_MistweaverMonk{},
		proto.Spec_SpecMistweaverMonk,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewMistweaverMonk(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_MistweaverMonk) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Mistweaver Monk!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewMistweaverMonk(character *core.Character, options *proto.Player) *MistweaverMonk {
	monkOptions := options.GetMistweaverMonk()

	mw := &MistweaverMonk{
		Monk:           monk.NewMonk(character, monkOptions.Options.ClassOptions, options.TalentsString),
		StartingStance: monkOptions.Options.Stance,
	}
	mw.SetStartingStance()

	return mw
}

func (mw *MistweaverMonk) SetStartingStance() {
	if mw.StartingStance == proto.MonkStance_WiseSerpent {
		mw.Monk.Stance = monk.WiseSerpent
	} else {
		mw.Monk.Stance = monk.FierceTiger
	}
}

type MistweaverMonk struct {
	*monk.Monk
	StartingStance proto.MonkStance
}

func (mw *MistweaverMonk) GetMonk() *monk.Monk {
	return mw.Monk
}

func (mw *MistweaverMonk) Initialize() {
	mw.Monk.Initialize()
	mw.RegisterSpecializationEffects()
}

func (mw *MistweaverMonk) ApplyTalents() {
	mw.Monk.ApplyTalents()
	mw.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeLeather)
}

func (mw *MistweaverMonk) Reset(sim *core.Simulation) {
	mw.SetStartingStance()
	mw.Monk.Reset(sim)
}

func (mw *MistweaverMonk) RegisterSpecializationEffects() {
	mw.RegisterMastery()
}

func (mw *MistweaverMonk) RegisterMastery() {
}
