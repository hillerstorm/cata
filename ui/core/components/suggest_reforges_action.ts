import { solve } from "yalps"
import { lessEq, equalTo, greaterEq, inRange } from "yalps"
import { Model, Constraint, Coefficients, OptimizationDirection, Options, Solution } from "yalps"
import { Player } from '../player.js';
import { Sim } from '../sim.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { Stats } from '../proto_utils/stats.js';
import { TypedEvent } from '../typed_event.js';
import { Gear } from '../proto_utils/gear.js';
import { ItemSlot, Stat } from '../proto/common.js';

interface StatWeightsConfig {
	statCaps: Stats;
	preCapEPs: Stats;
}

type YalpsCoefficients = Map<string, number>;
type YalpsVariables = Map<string, YalpsCoefficients>;
type YalpsConstraints = Map<string, Constraint>;

export const sleep = async (waitTime: number) =>
  new Promise(resolve =>
    setTimeout(resolve, waitTime));

export class ReforgeOptimizer {
	protected readonly player: Player<any>;
	protected readonly sim: Sim;
	protected readonly statCaps: Stats;
	protected readonly preCapEPs: Stats;

	constructor(simUI: IndividualSimUI<any>, config: StatWeightsConfig) {
		this.player = simUI.player;
		this.sim = simUI.sim;
		this.statCaps = config.statCaps;
		this.preCapEPs = config.preCapEPs;

		simUI.addAction('Suggest Reforges', 'suggest-reforges-action', async () => {
			this.optimizeReforges();
		});
	}

	async optimizeReforges() {
		console.log("Starting Reforge optimization...");

		// First, clear all existing Reforges
		console.log("Clearing existing Reforges...");
		const baseGear = this.player.getGear().withoutReforges(this.player.canDualWield2H());
		const baseStats = await this.updateGear(baseGear);

		// Compute effective stat caps for just the Reforge contribution
		const reforgeCaps = baseStats.computeStatCapsDelta(this.statCaps);
		console.log("Stat caps for Reforge contribution:");
		console.log(reforgeCaps);

		// Set up YALPS model
		const variables = this.buildYalpsVariables(baseGear);
		const constraints = this.buildYalpsConstraints(baseGear);

		// Solve in multiple passes to enforce caps
		await this.solveModel(baseGear, reforgeCaps, variables, constraints);
	}

	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		return Stats.fromProto(this.player.getCurrentStats().finalStats);
	}

	buildYalpsVariables(gear: Gear): YalpsVariables {
		const variables = new Map<string, YalpsCoefficients>();

		for (const slot of gear.getItemSlots()) {
			const item = gear.getEquippedItem(slot);

			if (!item) {
				continue;
			}

			for (const reforgeData of this.player.getAvailableReforgings(item)) {
				const variableKey = `${slot}_${reforgeData.id}`;
				const coefficients = new Map<string, number>();
				coefficients.set(ItemSlot[slot], 1);

				for (const fromStat of reforgeData.fromStat) {
					coefficients.set(Stat[fromStat], reforgeData.fromAmount);
				}
				
				for (const toStat of reforgeData.toStat) {
					coefficients.set(Stat[toStat], reforgeData.toAmount);
				}

				variables.set(variableKey, coefficients);
			}
		}

		return variables;
	}

	buildYalpsConstraints(gear: Gear): YalpsConstraints {
		const constraints = new Map<string, Constraint>();

		for (const slot of gear.getItemSlots()) {
			constraints.set(ItemSlot[slot], lessEq(1));
		}

		return constraints;
	}

	async solveModel(gear: Gear, reforgeCaps: Stats, variables: YalpsVariables, constraints: YalpsConstraints) {
		// Calculate EP scores for each Reforge option
		const updatedVariables = this.updateReforgeScores(variables, constraints);
		console.log("Optimization variables and constraints for this iteration:");
		console.log(updatedVariables);
		console.log(constraints);

		// Set up and solve YALPS model
		const model: Model = {
			direction: "maximize",
			objective: "score",
			constraints: constraints,
			variables: updatedVariables,
			binaries: true
		};
		const solution = solve(model);
		console.log("LP solution for this iteration:");
		console.log(solution);

		// Apply the current solution
		await this.applyLPSolution(gear, solution);

		// Check if any unconstrained stats exceeded their specified cap.
		// If so, add these stats to the constraint list and re-run the solver.
		// If no unconstrained caps were exceeded, then we're done.
		const [anyCapsExceeded, updatedConstraints] = this.checkCaps(solution, reforgeCaps, updatedVariables, constraints);

		if (!anyCapsExceeded) {
			console.log("Reforge optimization has converged!");
		} else {
			console.log("One or more stat caps were exceeded, starting constrained iteration...");
			await sleep(100);
			await this.solveModel(gear, reforgeCaps, updatedVariables, updatedConstraints);
		}
	}

	updateReforgeScores(variables: YalpsVariables, constraints: YalpsConstraints): YalpsVariables {
		const updatedVariables = new Map<string, YalpsCoefficients>();

		for (const [variableKey, coefficients] of variables.entries()) {
			let score = 0;
			const updatedCoefficients = new Map<string, number>();

			for (const [coefficientKey, value] of coefficients.entries()) {
				updatedCoefficients.set(coefficientKey, value);

				// Determine whether the key corresponds to a stat change.	
				// If so, check whether the stat has already been constrained to be capped in a previous iteration.
				// Apply stored EP only for unconstrained stats.
				if (coefficientKey.includes('Stat') && !constraints.has(coefficientKey)) {
					const statKey = (Stat as any)[coefficientKey] as Stat;
					score += this.preCapEPs.getStat(statKey) * value;
				}
			}

			updatedCoefficients.set("score", score);
			updatedVariables.set(variableKey, updatedCoefficients);
		}

		return updatedVariables;
	}

	async applyLPSolution(gear: Gear, solution: Solution) {
		let updatedGear = gear.withoutReforges(this.player.canDualWield2H());

		for (const [variableKey, _coefficient] of solution.variables) {
			const splitKey = variableKey.split("_");
			const slot = parseInt(splitKey[0]) as ItemSlot;
			const reforgeId = parseInt(splitKey[1]);
			const equippedItem = gear.getEquippedItem(slot);

			if (equippedItem) {
				updatedGear = updatedGear.withEquippedItem(slot, equippedItem.withReforge(this.sim.db.getReforgeById(reforgeId)!), this.player.canDualWield2H());
			}
		}

		await this.updateGear(updatedGear);
	}

	checkCaps(solution: Solution, reforgeCaps: Stats, variables: YalpsVariables, constraints: YalpsConstraints): [boolean, YalpsConstraints] {
		// First add up the total stat changes from the solution
		let reforgeStatContribution = new Stats();

		for (const [variableKey, _coefficient] of solution.variables) {
			for (const [coefficientKey, value] of variables.get(variableKey)!.entries()) {
				if (coefficientKey.includes('Stat')) {
					const statKey = (Stat as any)[coefficientKey] as Stat;
					reforgeStatContribution = reforgeStatContribution.addStat(statKey, value);
				}
			}
		}

		console.log("Total stat contribution from Reforging:");
		console.log(reforgeStatContribution);

		// Then check whether any unconstrained stats exceed their cap
		let anyCapsExceeded = false;
		const updatedConstraints = new Map<string, Constraint>(constraints);

		for (const [statKey, value] of reforgeStatContribution.asArray().entries()) {
			const cap = reforgeCaps.getStat(statKey);
			const statName = Stat[statKey];
			
			if ((cap != 0) && (value > cap) && !constraints.has(statName)) {
				updatedConstraints.set(statName, greaterEq(cap));
				anyCapsExceeded = true;
				console.log("Cap exceeded for: %s", statName); 
			}
		}

		return [anyCapsExceeded, updatedConstraints];
	}
}