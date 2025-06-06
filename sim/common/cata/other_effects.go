package cata

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	core.NewItemEffect(59461, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		dummyAura := character.RegisterAura(core.Aura{
			Label:     "Raw Fury",
			ActionID:  core.ActionID{SpellID: 91832},
			Duration:  time.Second * 15,
			MaxStacks: 5,
		})

		triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Raw Fury Aura",
			ActionID:   core.ActionID{ItemID: 59461},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.5,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 5,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				dummyAura.Activate(sim)
				dummyAura.AddStack(sim)
			},
		}))

		character.ItemSwap.RegisterProc(59461, triggerAura)

		buffAura := character.NewTemporaryStatsAura("Forged Fury", core.ActionID{SpellID: 91836}, stats.Stats{stats.Strength: 1926}, time.Second*20)
		sharedCD := character.GetOffensiveTrinketCD()
		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 59461},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,
			Cast: core.CastConfig{
				SharedCD: core.Cooldown{
					Timer:    sharedCD,
					Duration: time.Second * 20,
				},
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dummyAura.Deactivate(sim)
				buffAura.Activate(sim)
			},
			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return dummyAura.GetStacks() == 5
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    trinketSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
			BuffAura: buffAura,

			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return dummyAura.GetStacks() == 5
			},
		})
	})

	core.NewItemEffect(68996, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		totalAbsorbed := 0.0
		stayWithdrawnAura := character.RegisterAura(core.Aura{
			Label:    "Stay Withdrawn",
			Duration: time.Second * 10,
			ActionID: core.ActionID{SpellID: 96993},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				tickAmount := totalAbsorbed * 0.08

				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   time.Second * 2,
					NumTicks: 5,
					OnAction: func(sim *core.Simulation) {
						character.RemoveHealth(sim, tickAmount)

						// Hacky way to force a fake stacks change log for the timeline to mimic ticks
						if sim.Log != nil {
							stacks := aura.GetStacks()
							aura.Unit.Log(sim, "%s stacks: %d --> %d", aura.ActionID, stacks, stacks)
						}
					},
					CleanUp: func(sim *core.Simulation) {
						totalAbsorbed = 0
					},
				})
			},
		})

		maxShieldStrength := 56980.0
		absorbAura := character.RegisterAura(core.Aura{
			Label:     "Stay of Execution",
			ActionID:  core.ActionID{ItemID: 68996, SpellID: 96988},
			Duration:  time.Second * 30,
			MaxStacks: int32(maxShieldStrength),
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if totalAbsorbed > 0 {
					stayWithdrawnAura.Activate(sim)
					stacks := int32(totalAbsorbed * 0.08)
					if stacks > 0 {
						stayWithdrawnAura.MaxStacks = stacks
						stayWithdrawnAura.SetStacks(sim, stacks)
					}
				}
			},
		})

		character.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, isPeriodic bool) {
			if absorbAura.IsActive() && result.Damage > 0 && totalAbsorbed < maxShieldStrength {
				remainingAbsorb := maxShieldStrength - totalAbsorbed
				absorbedDamage := min(result.Damage*0.2, remainingAbsorb)
				result.Damage -= absorbedDamage
				totalAbsorbed = min(maxShieldStrength, totalAbsorbed+absorbedDamage)
				absorbAura.SetStacks(sim, int32(totalAbsorbed))

				if sim.Log != nil {
					character.Log(sim, "Stay of Execution absorbed %.1f damage", absorbedDamage)
				}
			}
		})

		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 68996},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,
			Cast: core.CastConfig{
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 20,
				},
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				absorbAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    trinketSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeSurvival,
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 55816,
		shared.ItemVersionHeroic: 56347,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			heroic := version == shared.ItemVersionHeroic
			actionID := core.ActionID{SpellID: core.TernaryInt32(heroic, 92184, 92179)}
			armorBonus := core.GetItemEffectScaling(itemID, core.TernaryFloat64(heroic, 5.94799995422, 5.9310002327), state)

			procAura := character.NewTemporaryStatsAura(
				fmt.Sprintf("Leaden Despair Proc %s", versionLabel),
				actionID,
				stats.Stats{stats.Armor: armorBonus},
				time.Second*10)

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 30,
			}
			procAura.Icd = &icd

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     fmt.Sprintf("Leaden Despair Trigger %s", versionLabel),
				Callback: core.CallbackOnSpellHitTaken,
				ProcMask: core.ProcMaskDirect,
				ActionID: core.ActionID{ItemID: itemID},
				Harmful:  true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.35 {
						icd.Use(sim)
						procAura.Activate(sim)
					}
				},
			})
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 59514,
		shared.ItemVersionHeroic: 65110,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			heroic := version == shared.ItemVersionHeroic

			procAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     fmt.Sprintf("Heart's Revelation %s", versionLabel),
					ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 92325, 91027)},
					Duration:  time.Second * 15,
					MaxStacks: 5,
				},
				BonusPerStack: stats.Stats{stats.SpellPower: core.GetItemEffectScaling(itemID, 0.11900000274, state)},
			})

			triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Heart of Ignacious Aura %s", versionLabel),
				ActionID:   core.ActionID{ItemID: itemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskSpellDamage,
				ProcChance: 1,
				Outcome:    core.OutcomeLanded,
				ICD:        time.Second * 2,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
					procAura.AddStack(sim)
				},
			}))

			character.ItemSwap.RegisterProc(itemID, triggerAura)

			hastePerStack := core.GetItemEffectScaling(itemID, 0.49500000477, state)
			buffAura := character.RegisterAura(core.Aura{
				Label:     fmt.Sprintf("Heart's Judgement %s", versionLabel),
				ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 92328, 91041)},
				Duration:  time.Second * 20,
				MaxStacks: 5,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					deltaHasteRating := hastePerStack * float64(newStacks-oldStacks)
					character.AddStatDynamic(sim, stats.HasteRating, deltaHasteRating)
				},
			})

			sharedCD := character.GetOffensiveTrinketCD()
			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: itemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    sharedCD,
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					buffAura.Activate(sim)
					buffAura.SetStacks(sim, procAura.GetStacks())
					procAura.Deactivate(sim)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeDPS,
				ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
					return procAura.GetStacks() == 5
				},
			})

			character.AddStatProcBuff(itemID, procAura, false, core.TrinketSlots())
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 59354,
		shared.ItemVersionHeroic: 65029,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			heroic := version == shared.ItemVersionHeroic

			manaReturn := core.GetItemEffectScaling(itemID, core.TernaryFloat64(heroic, 9.90499973297, 9.89200019836), state)

			bonusPerStack := core.GetItemEffectScaling(itemID, core.TernaryFloat64(heroic, 0.1589999944, 0.15800000727), state)
			procAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     fmt.Sprintf("Inner Eye %s", versionLabel),
					ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 91320, 92329)},
					Duration:  time.Second * 15,
					MaxStacks: 5,
					Icd: &core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Second * 30,
					},
				},

				BonusPerStack: stats.Stats{stats.Spirit: bonusPerStack},
			})

			triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Jar of Ancient Remedies Aura %s", versionLabel),
				ActionID:   core.ActionID{ItemID: itemID},
				Callback:   core.CallbackOnHealDealt,
				ProcMask:   core.ProcMaskSpellHealing,
				ProcChance: 1,
				Outcome:    core.OutcomeLanded,
				ICD:        time.Second * 2,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if procAura.Icd.IsReady(sim) {
						procAura.Activate(sim)
						procAura.AddStack(sim)
					}
				},
			}))

			character.ItemSwap.RegisterProc(itemID, triggerAura)

			manaMetric := character.NewManaMetrics(core.ActionID{SpellID: core.TernaryInt32(heroic, 92331, 91322)})
			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: itemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
					return character.HasManaBar()
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					procAura.Deactivate(sim)
					// can not regain stacks for 30 seconds
					procAura.Icd.Use(sim)
					character.AddMana(sim, manaReturn, manaMetric)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeMana,
			})
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 68994,
		shared.ItemVersionHeroic: 69150,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			heroic := version == shared.ItemVersionHeroic

			bonusStats := core.GetItemEffectScaling(itemID, core.TernaryFloat64(heroic, 1.98300004005, 1.98300004005), state)

			procAuraCrit := character.NewTemporaryStatsAura(
				fmt.Sprintf("Matrix Restabilizer Crit Proc %s", versionLabel),
				core.ActionID{SpellID: core.TernaryInt32(heroic, 97140, 96978)},
				stats.Stats{stats.CritRating: bonusStats},
				time.Second*30)
			procAuraHaste := character.NewTemporaryStatsAura(
				fmt.Sprintf("Matrix Restabilizer Haste Proc %s", versionLabel),
				core.ActionID{SpellID: core.TernaryInt32(heroic, 97139, 96977)},
				stats.Stats{stats.HasteRating: bonusStats},
				time.Second*30)
			procAuraMastery := character.NewTemporaryStatsAura(
				fmt.Sprintf("Matrix Restabilizer Mastery Proc %s", versionLabel),
				core.ActionID{SpellID: core.TernaryInt32(heroic, 97141, 96979)},
				stats.Stats{stats.MasteryRating: bonusStats},
				time.Second*30)

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 105,
			}

			procAuraCrit.Icd = &icd
			procAuraHaste.Icd = &icd
			procAuraMastery.Icd = &icd

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Matrix Restabilizer Trigger %s", versionLabel),
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged,
				ProcChance: 0.2,
				ActionID:   core.ActionID{ItemID: itemID},
				Harmful:    true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if icd.IsReady(sim) {
						statType := character.GetHighestStatType([]stats.Stat{stats.CritRating, stats.HasteRating, stats.MasteryRating})
						switch statType {
						case stats.CritRating:
							procAuraCrit.Activate(sim)
						case stats.HasteRating:
							procAuraHaste.Activate(sim)
						case stats.MasteryRating:
							procAuraMastery.Activate(sim)
						default:
							panic("unexpected statType")
						}
						icd.Use(sim)
					}
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 68972,
		shared.ItemVersionHeroic: 69113,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			bonusPerStack := core.GetItemEffectScaling(itemID, 0.65799999237, state)
			buffDuration := time.Second * 15

			buffAuraCrit, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				StackingAuraLabel:    fmt.Sprintf("Blessing of the Shaper Crit %s", versionLabel),
				StackingAuraActionID: core.ActionID{SpellID: 96928},
				BonusPerStack:        stats.Stats{stats.CritRating: bonusPerStack},
				Duration:             buffDuration,
				MaxStacks:            5,
			})

			buffAuraHaste, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				StackingAuraLabel:    fmt.Sprintf("Blessing of the Shaper Haste %s", versionLabel),
				StackingAuraActionID: core.ActionID{SpellID: 96927},
				BonusPerStack:        stats.Stats{stats.HasteRating: bonusPerStack},
				Duration:             buffDuration,
				MaxStacks:            5,
			})

			buffAuraMastery, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				StackingAuraLabel:    fmt.Sprintf("Blessing of the Shaper Mastery %s", versionLabel),
				StackingAuraActionID: core.ActionID{SpellID: 96929},
				BonusPerStack:        stats.Stats{stats.MasteryRating: bonusPerStack},
				Duration:             buffDuration,
				MaxStacks:            5,
			})

			buffAura := character.RegisterAura(core.Aura{
				Label:    fmt.Sprintf("Apparatus of Khaz'goroth %s", versionLabel),
				ActionID: core.ActionID{ItemID: itemID},
				Duration: buffDuration,
			})

			titanicPower := character.RegisterAura(core.Aura{
				Label:     fmt.Sprintf("Titanic Power %s", versionLabel),
				ActionID:  core.ActionID{SpellID: 96923},
				Duration:  time.Second * 30,
				MaxStacks: 5,
			})

			offHandBlacklist := map[int32]bool{
				// DK: Threat of Thassarian crits doesn't proc
				49143: true, // Frost Strike
				49998: true, // Death Strike
				85948: true, // Festering Strike
				49020: true, // Obliterate
				45462: true, // Plague Strike
				56815: true, // Rune Strike

				// Warrior: Whirlwind off-hand crits doesn't trigger procs
				1680: true, // Whirlwind
			}

			triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Titanic Power Trigger %s", versionLabel),
				ActionID:   core.ActionID{SpellID: 96924},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 1,
				Outcome:    core.OutcomeCrit,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if buffAuraCrit.IsActive() || buffAuraHaste.IsActive() || buffAuraMastery.IsActive() {
						return
					}

					// Warrior: Raging Blow crits doesn't trigger procs (both MH and OH)
					if spell.ActionID.SpellID == 85288 {
						return
					}

					// Off-hand blacklist
					if _, blacklisted := offHandBlacklist[spell.ActionID.SpellID]; spell.ProcMask.Matches(core.ProcMaskMeleeOHSpecial) && blacklisted {
						return
					}

					titanicPower.Activate(sim)
					titanicPower.AddStack(sim)
				},
			}))

			character.ItemSwap.RegisterProc(itemID, triggerAura)

			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: itemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    character.GetOffensiveTrinketCD(),
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					statType := character.GetHighestStatType([]stats.Stat{stats.CritRating, stats.HasteRating, stats.MasteryRating})

					switch statType {
					case stats.CritRating:
						buffAuraCrit.Activate(sim)
						buffAuraCrit.SetStacks(sim, titanicPower.GetStacks())
					case stats.HasteRating:
						buffAuraHaste.Activate(sim)
						buffAuraHaste.SetStacks(sim, titanicPower.GetStacks())
					case stats.MasteryRating:
						buffAuraMastery.Activate(sim)
						buffAuraMastery.SetStacks(sim, titanicPower.GetStacks())
					default:
						panic("unexpected statType")
					}

					buffAura.Activate(sim)
					titanicPower.Deactivate(sim)
				},
				ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
					return titanicPower.IsActive()
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeDPS,
				ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
					return titanicPower.IsActive()
				},
			})
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 68995,
		shared.ItemVersionHeroic: 69167,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			heroic := version == shared.ItemVersionHeroic
			bonusPerStack := core.GetItemEffectScaling(itemID, core.TernaryFloat64(heroic, 0.09899999946, 0.10000000149), state)
			procAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     fmt.Sprintf("Accelerated %s", versionLabel),
					ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 97142, 96980)},
					Duration:  time.Second * 20,
					MaxStacks: 5,
				},
				BonusPerStack: stats.Stats{stats.CritRating: bonusPerStack},
			})

			offHandBlacklist := map[int32]bool{
				// Warrior: Slam and Whirlwind off-hand crits doesn't trigger procs
				1464: true, // Slam
				1680: true, // Whirlwind
			}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{ItemID: itemID},
				Name:       fmt.Sprintf("Vessel of Acceleration %s", versionLabel),
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrMeleeProc,
				Outcome:    core.OutcomeCrit,
				ProcChance: 1,
				Harmful:    false,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					// Off-hand blacklist
					if _, blacklisted := offHandBlacklist[spell.ActionID.SpellID]; spell.ProcMask.Matches(core.ProcMaskMeleeOHSpecial) && blacklisted {
						return
					}

					procAura.Activate(sim)
					procAura.AddStack(sim)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 68926,
		shared.ItemVersionHeroic: 69111,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			if !character.HasManaBar() {
				return
			}

			heroic := version == shared.ItemVersionHeroic
			manaMod := character.AddDynamicMod(core.SpellModConfig{
				School:   core.SpellSchoolHoly | core.SpellSchoolNature,
				IntValue: 0,
				Kind:     core.SpellMod_PowerCost_Flat,
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Millisecond * 900,
			}

			manaReturn := core.GetItemEffectScaling(itemID, core.TernaryFloat64(heroic, -0.14300000668, -0.14200000465), state)

			victoriousAura := character.GetOrRegisterAura(core.Aura{
				Label:     fmt.Sprintf("Victorious %s", versionLabel),
				ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 97120, 96907)},
				Duration:  time.Second * 20,
				MaxStacks: 10,
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.ProcMask&(core.ProcMaskSpellDamage|core.ProcMaskSpellHealing) == 0 {
						return
					}

					if !icd.IsReady(sim) {
						return
					}

					icd.Use(sim)

					aura.AddStack(sim)
				},
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Deactivate(sim)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					manaMod.Activate()
				},
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					manaMod.UpdateIntValue(int32(math.Floor(manaReturn * float64(newStacks))))
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					manaMod.Deactivate()
					manaMod.UpdateIntValue(0)
				},
			})

			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: itemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    character.GetOffensiveTrinketCD(),
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					victoriousAura.Activate(sim)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeMana,
			})
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 68981,
		shared.ItemVersionHeroic: 69138,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			heroic := version == shared.ItemVersionHeroic

			shieldStrength := core.GetItemEffectScaling(itemID, core.TernaryFloat64(heroic, 22.06299972534, 22.05800056458), state)
			actionID := core.ActionID{ItemID: itemID, SpellID: core.TernaryInt32(heroic, 97129, 96945)}
			duration := time.Second * 30

			shield := character.NewDamageAbsorptionAura(core.AbsorptionAuraConfig{
				Aura: core.Aura{
					Label:    fmt.Sprintf("Loom of Fate %s", versionLabel),
					ActionID: actionID,
					Duration: duration,
				},
				ShieldStrengthCalculator: func(unit *core.Unit) float64 {
					return shieldStrength
				},
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute,
			}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Spidersilk Spindle Trigger %s", versionLabel),
				Callback:   core.CallbackOnSpellHitTaken,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					preHitHp := character.CurrentHealth() + result.Damage
					if icd.IsReady(sim) && spell.SpellSchool == core.SpellSchoolPhysical &&
						character.CurrentHealthPercent() < 0.35 && preHitHp >= 0.35 {
						icd.Use(sim)
						shield.Activate(sim)
					}
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionNormal: 68915,
		shared.ItemVersionHeroic: 69109,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			heroic := version == shared.ItemVersionHeroic

			// Assuming full stack since sim doesn't track overhealing
			maxHeal := core.TernaryFloat64(heroic, 19283, 17095)
			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: itemID},
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskSpellHealing,
				Flags:       core.SpellFlagNoOnCastComplete,

				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    character.GetDefensiveTrinketCD(),
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 1,
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealHealing(sim, spell.Unit, maxHeal, spell.OutcomeHealingCrit)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeSurvival,
			})
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionLFR:    77979,
		shared.ItemVersionNormal: 77207,
		shared.ItemVersionHeroic: 77999,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			actionID := core.ActionID{SpellID: []int32{109721, 107994, 109724}[version]}
			baseDmg := core.GetItemEffectScaling(itemID, []float64{5.44600009918, 5.44299983978, 5.44299983978}[version], state)
			minDmg, maxDmg := core.ApplyVarianceMinMax(baseDmg, 0.40000000596)

			// 01/26/25: Testing during the first Dragon Soul PTR is consistent with these scaling
			// coefficients from simC, refer to Discord discussion from here onwards:
			// https://discord.com/channels/891730968493305867/1034670402150092841/1332728469427322982
			apMod := []float64{0.266, 0.3, 0.339}[version]

			// The AP scaling calculation can use either Melee or
			// Ranged AP depending on how the proc was triggered, so
			// keep track of it separately.
			var apSnapshot float64

			lightningStrike := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty, // ProcMask is set in the trigger aura
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(), // even ranged procs use the melee Crit multiplier as of first PTR (1.5x for Hunters)
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := sim.Roll(minDmg, maxDmg) + apMod*apSnapshot
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Vial of Shadows Trigger %s", versionLabel),
				ActionID:   core.ActionID{ItemID: itemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.45,
				ICD:        time.Second * 9,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskMelee | core.ProcMaskMeleeProc) {
						apSnapshot = spell.MeleeAttackPower()
					} else {
						apSnapshot = spell.RangedAttackPower()
					}
					lightningStrike.ProcMask = core.Ternary(spell.ProcMask.Matches(core.ProcMaskRanged), core.ProcMaskRangedProc, core.ProcMaskMeleeProc)
					lightningStrike.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionLFR:    77982,
		shared.ItemVersionNormal: 77210,
		shared.ItemVersionHeroic: 78002,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			numTargets := character.Env.GetNumTargets()

			actionID := core.ActionID{SpellID: []int32{109753, 107998, 109755}[version]}
			baseDmg := core.GetItemEffectScaling(itemID, []float64{12.25500011444, 12.24600028992, 12.24899959564}[version], state)
			minDmg, maxDmg := core.ApplyVarianceMinMax(baseDmg, 0.40000000596)
			apMod := []float64{0.598, 0.675, 0.762}[version]

			whirlingMaw := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskMeleeProc,
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					results := make([]*core.SpellResult, numTargets)

					for idx := int32(0); idx < numTargets; idx++ {
						baseDamage := sim.Roll(minDmg, maxDmg) +
							apMod*spell.MeleeAttackPower()
						results[idx] = spell.CalcDamage(sim, sim.Environment.GetTargetUnit(idx), baseDamage, spell.OutcomeMeleeSpecialCritOnly)
					}

					for idx := int32(0); idx < numTargets; idx++ {
						spell.DealDamage(sim, results[idx])
					}
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Bone-Link Fetish Trigger %s", versionLabel),
				ActionID:   core.ActionID{ItemID: itemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.15,
				ICD:        time.Second * 27,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					whirlingMaw.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionLFR:    77980,
		shared.ItemVersionNormal: 77208,
		shared.ItemVersionHeroic: 78000,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			numTargets := character.Env.GetNumTargets()

			actionID := core.ActionID{SpellID: []int32{109798, 108005, 109800}[version]}
			baseDmg := core.GetItemEffectScaling(itemID, []float64{3.81299996376, 3.81100010872, 3.81100010872}[version], state)
			minDmg, maxDmg := core.ApplyVarianceMinMax(baseDmg, 0.40000000596)
			spMod := []float64{0.277, 0.313, 0.353}[version]

			shadowboltVolley := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolShadow,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				MissileSpeed: 20,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				BonusCoefficient: spMod,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					results := make([]*core.SpellResult, numTargets)

					for idx := int32(0); idx < numTargets; idx++ {
						results[idx] = spell.CalcDamage(sim, sim.Environment.GetTargetUnit(idx), sim.Roll(minDmg, maxDmg), spell.OutcomeMagicCrit)
					}

					spell.WaitTravelTime(sim, func(sim *core.Simulation) {
						for idx := int32(0); idx < numTargets; idx++ {
							spell.DealDamage(sim, results[idx])
						}
					})
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Cunning of the Cruel Trigger %s", versionLabel),
				ActionID:   core.ActionID{ItemID: itemID},
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				ProcMask:   core.ProcMaskSpellOrSpellProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.45,
				ICD:        time.Second * 9,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shadowboltVolley.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	shared.ItemVersionMap{
		shared.ItemVersionLFR:    77983,
		shared.ItemVersionNormal: 77211,
		shared.ItemVersionHeroic: 78003,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			absorbModifier := []float64{0.43, 0.50, 0.56}[version]
			actionID := core.ActionID{ItemID: itemID, SpellID: 108008}
			duration := time.Second * 6

			shieldStrength := 0.0
			shield := character.NewDamageAbsorptionAura(core.AbsorptionAuraConfig{
				Aura: core.Aura{
					Label:    fmt.Sprintf("Indomitable %s", versionLabel),
					ActionID: actionID,
					Duration: duration,
				},
				ShieldStrengthCalculator: func(unit *core.Unit) float64 {
					return shieldStrength
				},
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute,
			}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Indomitable Pride Trigger %s", versionLabel),
				Callback:   core.CallbackOnSpellHitTaken,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					preHitHp := character.CurrentHealth() + result.Damage
					if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.50 && preHitHp >= 0.50 {
						shieldStrength = result.Damage * absorbModifier
						if shieldStrength > 1 {
							icd.Use(sim)
							shield.Activate(sim)
						}
					}
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	// Kiril, Fury of Beasts
	// Equip: Your melee and ranged attacks have a chance to trigger Fury of the Beast, granting 107 Agility and 10% increased size every 1 sec.
	// This effect stacks a maximum of 10 times and lasts 20 sec.
	// (Proc chance: 15%, 55s cooldown)
	// TODO: Verify if the aura is cancelled when swapping druid forms
	// Video from 4.3.0 showing that it doesn't: https://www.youtube.com/watch?v=A6PYbDRaH6E
	// Comment from 4.3.3 stating that it does: https://www.wowhead.com/mop-classic/item=77194/kiril-fury-of-beasts#comments:id=1639024
	shared.ItemVersionMap{
		shared.ItemVersionLFR:    78482,
		shared.ItemVersionNormal: 77194,
		shared.ItemVersionHeroic: 78473,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			bonusPerStack := core.GetItemEffectScaling(itemID, []float64{0.1099999994, 0.10899999738, 0.10899999738}[version], state)

			beastFuryAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     fmt.Sprintf("Beast Fury %s", versionLabel),
					ActionID:  core.ActionID{SpellID: []int32{109860, 108016, 109863}[version]},
					Duration:  time.Second * 20,
					MaxStacks: 10,
				},
				BonusPerStack: stats.Stats{stats.Agility: bonusPerStack},
			})

			furyOfTheBeastAura := character.RegisterAura(core.Aura{
				Label:    fmt.Sprintf("Fury of the Beast %s", versionLabel),
				ActionID: core.ActionID{SpellID: []int32{109861, 108011, 109864}[version]},
				Duration: time.Second * 20,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					beastFuryAura.Activate(sim)
					core.StartPeriodicAction(sim, core.PeriodicActionOptions{
						Period:   time.Second,
						NumTicks: 10,
						OnAction: func(sim *core.Simulation) {
							beastFuryAura.AddStack(sim)
						},
					})
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Fury of the Beast Trigger %s", versionLabel),
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.15,
				ICD:        time.Second * 55,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					furyOfTheBeastAura.Activate(sim)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	// These spells ignore the slot the weapon is in.
	// Any other ability should only trigger the proc if the weapon is in the right slot.
	ignoresSlot := map[int32]bool{
		23881: true, // Bloodthirst
		6544:  true, // Heroic Leap
	}

	// Souldrinker
	// Equip: Your melee attacks have a chance to drain your target's health, damaging the target for an amount equal to 1.3%/1.5%/1.7% of your maximum health and healing you for twice that amount.
	// (Proc chance: 15%)
	shared.ItemVersionMap{
		shared.ItemVersionLFR:    78488,
		shared.ItemVersionNormal: 77193,
		shared.ItemVersionHeroic: 78479,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: []int32{109828, 108022, 109831}[version]}
			label := fmt.Sprintf("Drain Life Trigger %s", versionLabel)
			hpModifier := []float64{0.013, 0.015, 0.017}[version]
			meleeWeaponSlots := core.AllWeaponSlots()

			var damageDealt float64
			drainLifeHeal := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID.WithTag(2),
				SpellSchool: core.SpellSchoolShadow,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealHealing(sim, target, damageDealt*2, spell.OutcomeAlwaysHit)
				},
			})

			drainLife := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID.WithTag(1),
				SpellSchool: core.SpellSchoolShadow,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := character.MaxHealth() * hpModifier
					damageDealt = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit).Damage

					drainLifeHeal.Cast(sim, &character.Unit)
				},
			})

			makeProcTrigger := func(character *core.Character, isMH bool) {
				itemSlot := core.Ternary(isMH, meleeWeaponSlots[:1], meleeWeaponSlots[1:])
				procMask := core.Ternary(isMH, core.ProcMaskMeleeMH|core.ProcMaskMeleeProc, core.ProcMaskMeleeOH)

				aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
					Name:     fmt.Sprintf("%s %s", label, core.Ternary(isMH, "MH", "OH")),
					ActionID: core.ActionID{ItemID: itemID},
					ProcMask: core.ProcMaskMelee,
					Outcome:  core.OutcomeLanded,
					Callback: core.CallbackOnSpellHitDealt,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if _, ignore := ignoresSlot[spell.ActionID.SpellID]; !spell.ProcMask.Matches(procMask) && !ignore {
							return
						}

						if sim.Proc(0.15, label) {
							drainLife.Cast(sim, result.Target)
						}
					},
				})

				character.ItemSwap.RegisterProcWithSlots(itemID, aura, itemSlot)
			}

			if character.ItemSwap.CouldHaveItemEquippedInSlot(itemID, proto.ItemSlot_ItemSlotMainHand) {
				makeProcTrigger(character, true)
			}
			if character.ItemSwap.CouldHaveItemEquippedInSlot(itemID, proto.ItemSlot_ItemSlotOffHand) {
				makeProcTrigger(character, false)
			}
		})
	})

	// No'Kaled, the Elements of Death
	// Equip: Your melee attacks have a chance to blast your enemy with Fire, Shadow, or Frost, dealing 6781/7654/8640 to 10171/11481/12960 damage.
	// (Proc chance: 7%)
	shared.ItemVersionMap{
		shared.ItemVersionLFR:    78481,
		shared.ItemVersionNormal: 77188,
		shared.ItemVersionHeroic: 78472,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			defaultProcMask := core.ProcMaskMeleeProc
			procMask := character.GetDynamicProcMaskForWeaponEffect(itemID)

			baseDamage := core.GetItemEffectScaling(itemID, []float64{9.78800010681, 9.78299999237, 9.73799991608}[version], state)
			minDamage, maxDamage := core.ApplyVarianceMinMax(baseDamage, 0.40000000596)

			registerSpell := func(actionID core.ActionID, spellSchool core.SpellSchool) *core.Spell {
				return character.RegisterSpell(core.SpellConfig{
					ActionID:    actionID,
					SpellSchool: spellSchool,
					ProcMask:    core.ProcMaskEmpty,
					Flags:       core.SpellFlagPassiveSpell,
					MaxRange:    45,

					DamageMultiplier: 1,
					CritMultiplier:   character.DefaultCritMultiplier(),
					ThreatMultiplier: 1,

					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						baseDamage := sim.RollWithLabel(minDamage, maxDamage, "No'Kaled, the Elements of Death")
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
					},
				})
			}

			flameblast := registerSpell(
				core.ActionID{SpellID: []int32{109871, 107785, 109872}[version]},
				core.SpellSchoolFire)

			iceblast := registerSpell(
				core.ActionID{SpellID: []int32{109869, 107789, 109870}[version]},
				core.SpellSchoolFrost)

			shadowblast := registerSpell(
				core.ActionID{SpellID: []int32{109867, 107787, 109868}[version]},
				core.SpellSchoolShadow)

			spells := []*core.Spell{flameblast, iceblast, shadowblast}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     fmt.Sprintf("No'Kaled Trigger %s", versionLabel),
				ActionID: core.ActionID{ItemID: itemID},
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}

					if _, ignore := ignoresSlot[spell.ActionID.SpellID]; !spell.ProcMask.Matches(*procMask|defaultProcMask) && !ignore {
						return
					}

					if sim.Proc(0.07, "No'Kaled, the Elements of Death") {
						spell := spells[int(sim.RollWithLabel(0, float64(len(spells)), "No'Kaled spell to cast"))]
						spell.Cast(sim, result.Target)
					}
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	// Rathrak, the Poisonous Mind
	// Equip: Your harmful spellcasts have a chance to poison all enemies near your target for 7715/8710/9830 nature damage over 10 sec.
	// (Proc chance: 15%, 17s cooldown)
	shared.ItemVersionMap{
		shared.ItemVersionLFR:    78484,
		shared.ItemVersionNormal: 77195,
		shared.ItemVersionHeroic: 78475,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			tickDamage := core.GetItemEffectScaling(itemID, []float64{1.78199994564, 1.78100001812, 1.78100001812}[version], state)

			blastOfCorruption := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: []int32{109851, 107831, 109854}[version]},
				SpellSchool: core.SpellSchoolNature,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				Dot: core.DotConfig{
					IsAOE: true,
					Aura: core.Aura{
						Label: fmt.Sprintf("Blast of Corruption %s", versionLabel),
					},
					NumberOfTicks:       5,
					TickLength:          time.Second * 2,
					AffectedByCastSpeed: false,

					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						for _, aoeTarget := range sim.Encounter.TargetUnits {
							result := dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, tickDamage, dot.Spell.OutcomeMagicCritNoHitCounter)

							if result.DidCrit() {
								dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
							} else {
								dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
							}
						}
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.AOEDot().Apply(sim)
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Rathrak Trigger %s", versionLabel),
				ActionID:   core.ActionID{ItemID: itemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskSpellOrSpellProc,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.15,
				ICD:        time.Second * 17,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					blastOfCorruption.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	// Vishanka, Jaws of the Earth
	// Equip: Your ranged attacks have a chance to deal 7040/7950/8970 damage over 2 sec.
	// (Proc chance: 15%, 17s cooldown)
	// Time between ticks: 200ms
	shared.ItemVersionMap{
		shared.ItemVersionLFR:    78480,
		shared.ItemVersionNormal: 78359,
		shared.ItemVersionHeroic: 78471,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			tickDamage := []float64{760, 859, 970}[version]

			speakingOfRage := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: []int32{109856, 107821, 109858}[version]},
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: fmt.Sprintf("Speaking of Rage %s", versionLabel),
					},
					NumberOfTicks:       10,
					TickLength:          time.Millisecond * 200,
					AffectedByCastSpeed: false,

					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						result := dot.Spell.CalcAndDealPeriodicDamage(sim, target, tickDamage, dot.Spell.OutcomeRangedCritOnlyNoHitCounter)

						if result.DidCrit() {
							dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
						} else {
							dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
						}
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.CritMultiplier(1, 0),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).Apply(sim)
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       fmt.Sprintf("Vishanka Trigger %s", versionLabel),
				ActionID:   core.ActionID{ItemID: itemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskRanged | core.ProcMaskRangedProc,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.15,
				ICD:        time.Second * 17,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					speakingOfRage.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})

	// Ti'tahk, the Steps of Time
	// Equip: Your spells have a chance to grant you 1708/1928/2176 haste rating for 10 sec and 342/386/435 haste rating to up to 3 allies within 20 yards.
	// (Proc chance: 15%, 50s cooldown)
	// The buff has two effects, one for the caster and one shared.
	// * The first effect is 1366/1542/1741 haste rating on the caster.
	// * The second effect is 342/386/435 haste rating on the caster and up to 3 allies within 20 yards.
	// E.g. for the LFR version it's 1366 + 342 = 1708 haste rating for the caster with a shared ID so we just combine them.
	// TODO: Add the shared buff as an optional misc raid buff?
	//       Could be annoying with 3 different versions, uptime etc.
	shared.ItemVersionMap{
		shared.ItemVersionLFR:    78486,
		shared.ItemVersionNormal: 77190,
		shared.ItemVersionHeroic: 78477,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		titahkLabel := fmt.Sprintf("Ti'tahk, the Steps of Time %s", versionLabel)
		titahkAuraID := []int32{109842, 107804, 109844}[version]

		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			titahkBonusBase := core.GetItemEffectScaling(itemID, 1.57700002193, state)
			titahkBonusShared := core.GetItemEffectScaling(itemID, []float64{0.39500001073, 0.39500001073, 0.3939999938}[version], state)
			titahkBonus := titahkBonusBase + titahkBonusShared
			procAura := character.NewTemporaryStatsAura(
				titahkLabel,
				core.ActionID{SpellID: titahkAuraID},
				stats.Stats{stats.HasteRating: titahkBonus},
				time.Second*10)

			handler := func(triggerAura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOrMeleeProc | core.ProcMaskRangedOrRangedProc) {
					return
				}
				if triggerAura.Icd.IsReady(sim) && sim.Proc(0.15, fmt.Sprintf("%s Trigger", titahkLabel)) {
					procAura.Activate(sim)
					triggerAura.Icd.Use(sim)
				}
			}

			triggerAura := character.RegisterAura(core.Aura{
				Label:    fmt.Sprintf("%s Trigger", titahkLabel),
				ActionID: core.ActionID{ItemID: itemID},
				Duration: core.NeverExpires,
				Icd: &core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 50,
				},
				OnSpellHitDealt:       handler,
				OnPeriodicDamageDealt: handler,
				OnHealDealt:           handler,
				OnPeriodicHealDealt:   handler,
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					handler(aura, sim, spell, nil)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})
}
