package cata

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	// Synapse Springs
	core.NewEnchantEffect(4179, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		bonus := stats.Stats{}
		bonus[character.GetHighestStatType([]stats.Stat{
			stats.Strength, stats.Agility, stats.Intellect,
		})] = 480

		core.RegisterTemporaryStatsOnUseCD(character,
			"Synapse Springs",
			bonus,
			10*time.Second,
			core.SpellConfig{
				ActionID: core.ActionID{SpellID: 82174},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute,
					},
					SharedCD: core.Cooldown{
						Timer:    character.GetOffensiveTrinketCD(),
						Duration: 10 * time.Second,
					},
				},
			})
	})
}
