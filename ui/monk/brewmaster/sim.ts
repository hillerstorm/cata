import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, HandType, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import { Sim } from '../../core/sim';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecBrewmasterMonk, {
	cssClass: 'brewmaster-monk-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Monk),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
		Stat.StatExpertiseRating,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
		PseudoStat.PseudoStatPhysicalHitPercent,
		PseudoStat.PseudoStatSpellHitPercent,
	],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatAgility, Stat.StatStrength, Stat.StatAttackPower, Stat.StatMasteryRating, Stat.StatExpertiseRating],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.PREPATCH_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.PREPATCH_EP_PRESET.epWeights,
		// Stat caps for reforge optimizer
		statCaps: (() => {
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
			return expCap;
		})(),
		softCapBreakpoints: (() => {
			const meleeHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [8, 27],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0, 0],
			});

			return [meleeHitSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			arcaneBrilliance: true,
			bloodlust: true,
			markOfTheWild: true,
			icyTalons: true,
			moonkinForm: true,
			leaderOfThePack: true,
			powerWordFortitude: true,
			strengthOfEarthTotem: true,
			trueshotAura: true,
			wrathOfAirTotem: true,
			demonicPact: true,
			blessingOfKings: true,
			blessingOfMight: true,
			legacyOfTheEmperor: true,
			legacyOfTheWhiteTiger: true,
			communion: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			mangle: true,
			sunderArmor: true,
			shadowAndFlame: true,
			earthAndMoon: true,
			bloodFrenzy: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.CritBuff,
		BuffDebuffInputs.MajorArmorDebuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InFrontOfTarget,
			OtherInputs.InputDelay,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.PREPATCH_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DefaultTalents],
		// Preset rotations that the user can quickly select.
		rotations: [],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PREPATCH_GEAR_PRESET],
	},

	autoRotation: (_: Player<Spec.SpecBrewmasterMonk>): APLRotation => {
		return APLRotation.create();
	},

	raidSimPresets: [
		{
			spec: Spec.SpecBrewmasterMonk,
			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceAlliancePandaren,
				[Faction.Horde]: Race.RaceHordePandaren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.PREPATCH_GEAR_PRESET.gear,
					2: Presets.PREPATCH_GEAR_PRESET.gear,
					3: Presets.PREPATCH_GEAR_PRESET.gear,
					4: Presets.PREPATCH_GEAR_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.PREPATCH_GEAR_PRESET.gear,
					2: Presets.PREPATCH_GEAR_PRESET.gear,
					3: Presets.PREPATCH_GEAR_PRESET.gear,
					4: Presets.PREPATCH_GEAR_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

const getActiveEPWeight = (player: Player<Spec.SpecBrewmasterMonk>, sim: Sim): Stats => {
	if (sim.getUseCustomEPValues()) {
		return player.getEpWeights();
	} else {
		return Presets.PREPATCH_EP_PRESET.epWeights;
	}
}

export class BrewmasterMonkSimUI extends IndividualSimUI<Spec.SpecBrewmasterMonk> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecBrewmasterMonk>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateSoftCaps: (softCaps: StatCap[]) => {
					// Dynamic adjustments to the static Hit soft cap EP
					const meleeSoftCap = softCaps.find(v => v.unitStat.equalsPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent));
					if (meleeSoftCap) {
						const activeEPWeight = getActiveEPWeight(player, this.sim);
						const initialEP = activeEPWeight.getPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent);
						const mhWep = player.getEquippedItem(ItemSlot.ItemSlotMainHand);
						const ohWep = player.getEquippedItem(ItemSlot.ItemSlotOffHand);
						if (mhWep?.item.handType === HandType.HandTypeTwoHand || !ohWep) {
							meleeSoftCap.breakpoints = [meleeSoftCap.breakpoints[0]];
							meleeSoftCap.postCapEPs = [0];
						} else if (ohWep) {
							meleeSoftCap.postCapEPs = [initialEP/2, 0]
						}
					}

					return softCaps
				},
				getEPDefaults: (player: Player<Spec.SpecBrewmasterMonk>) => {
					return getActiveEPWeight(player, this.sim);
				},
			});
		});
	}
}