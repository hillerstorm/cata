import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common';
import { MonkMajorGlyph, MonkMinorGlyph,MonkOptions } from '../../core/proto/monk';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import DefaultGear from './gear_sets/default.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PREPATCH_GEAR_PRESET = PresetUtils.makePresetGear('Default', DefaultGear);

export const PREPATCH_ROTATION_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const PREPATCH_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 3.07,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.47,
			[Stat.StatHitRating]: 3.42,
			[Stat.StatHasteRating]: 2.13,
			[Stat.StatMasteryRating]: 0.31,
			[Stat.StatExpertiseRating]: 3.41,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 9.58,
			[PseudoStat.PseudoStatOffHandDps]: 4.79,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 410.78,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '001100001001001010',
		glyphs: Glyphs.create({
			major1: MonkMajorGlyph.MonkMajorGlyphSpinningCraneKick,
			major2: MonkMajorGlyph.MonkMajorGlyphTouchOfKarma,
			major3: MonkMajorGlyph.MonkMajorGlyphZenMeditation,
			minor1: MonkMinorGlyph.MonkMinorGlyphBlackoutKick,
			minor2: MonkMinorGlyph.MonkMinorGlyphJab,
			minor3: MonkMinorGlyph.MonkMinorGlyphWaterRoll,
		}),
	}),
};

export const DefaultOptions = MonkOptions.create({
	classOptions: {},
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodSeafoodFeast,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};