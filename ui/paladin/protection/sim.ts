import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation } from '../../core/proto/apl.js';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, TristateEffect } from '../../core/proto/common.js';
import { PaladinMajorGlyph, PaladinSeal, ProtectionPaladin_Rotation as ProtectionPaladinRotation } from '../../core/proto/paladin.js';
import * as AplUtils from '../../core/proto_utils/apl_utils.js';
import { Stats } from '../../core/proto_utils/stats.js';
import { TypedEvent } from '../../core/typed_event.js';
import * as PaladinInputs from '../inputs.js';
// import * as ProtInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecProtectionPaladin, {
	cssClass: 'protection-paladin-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Paladin),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatSpellHit,
		Stat.StatMeleeCrit,
		Stat.StatExpertise,
		Stat.StatMeleeHaste,
		Stat.StatSpellPower,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatBlock,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatResilience,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatMastery,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatBlock,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatResilience,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatMastery,
	],
	// modifyDisplayStats: (player: Player<Spec.SpecProtectionPaladin>) => {
	// 	let stats = new Stats();

	// 	TypedEvent.freezeAllAndDo(() => {
	// 		if (
	// 			player.getMajorGlyphs().includes(PaladinMajorGlyph.GlyphOfSealOfVengeance) &&
	// 			player.getSpecOptions().classOptions?.seal == PaladinSeal.Vengeance
	// 		) {
	// 			stats = stats.addStat(Stat.StatExpertise, 10 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
	// 		}
	// 	});

	// 	return {
	// 		talents: stats,
	// 	};
	// },
	defaults: {
		// Default equipped gear.
		gear: Presets.PRERAID_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		// Values for now are pre-Cata initial WAG
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.GenericAoeTalents.data,
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
			communion: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({
			vampiricTouch: true,
		}),
		debuffs: Debuffs.create({
			ebonPlaguebringer: true,
			shadowAndFlame: true,
			bloodFrenzy: true,
			mangle: true,
			exposeArmor: true,
			sunderArmor: true,
			vindication: true,
			thunderClap: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [PaladinInputs.AuraSelection(), PaladinInputs.StartingSealSelection()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InspirationUptime,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.GenericAoeTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_DEFAULT],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
			Presets.T11_PRESET,
			Presets.T11CTC_PRESET,
		],
	},

	autoRotation: (_player: Player<Spec.SpecProtectionPaladin>): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (_player: Player<Spec.SpecProtectionPaladin>, simple: ProtectionPaladinRotation, cooldowns: Cooldowns): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecProtectionPaladin,
			talents: Presets.GenericAoeTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceBloodElf,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.PRERAID_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.PRERAID_PRESET.gear,
				},
			},
		},
	],
});

export class ProtectionPaladinSimUI extends IndividualSimUI<Spec.SpecProtectionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
