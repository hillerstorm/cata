import * as PresetUtils from '../../core/preset_utils.js';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common.js';
import { PaladinMajorGlyph, PaladinSeal, RetributionPaladin_Options as RetributionPaladinOptions } from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P4_Gear from './gear_sets/p4.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so it's good to
// keep them in a separate file.

export const P4_GEAR_PRESET = PresetUtils.makePresetGear('P4', P4_Gear);

export const APL_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P4_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P4',
	Stats.fromMap(
		{
			[Stat.StatHitRating]: 1.5,
			[Stat.StatExpertiseRating]: 1.3,
			[Stat.StatStrength]: 1.0,
			[Stat.StatHasteRating]: 0.81,
			[Stat.StatMasteryRating]: 0.8,
			[Stat.StatCritRating]: 0.65,
			[Stat.StatAttackPower]: 0.44,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 2.18,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '221223',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfTemplarsVerdict,
			major2: PaladinMajorGlyph.GlyphOfDoubleJeopardy,
		}),
	}),
};

export const P4_BUILD_PRESET = PresetUtils.makePresetBuild('P4', {
	gear: P4_GEAR_PRESET,
	epWeights: P4_EP_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
});

export const DefaultOptions = RetributionPaladinOptions.create({
	classOptions: {
		seal: PaladinSeal.Truth,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58088, // Flask of Titanic Strength
	foodId: 62670, // Beer-Basted Crocolisk
	potId: 58146, // Golemblood Potion
	prepotId: 58146, // Golemblood Potion
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};
