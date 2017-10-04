import { Injectable, Optional, Inject } from '@angular/core';
import { RuleBase, BLWLRule } from 'app/components/core/interfaces/rule.interface';

export class RuleCreator {
	constructor(public ruleBase: RuleBase) {}
	
	get terms(): Array<string> {
		return this.ruleBase.list_terms.split(',')
	}

	get comparekeys(): Array<string> {
		return this.ruleBase.compare_key.split(',')
	}

	set ruleMeta(ruleItem: BLWLRule) {
		this.ruleBase.rule_meta = ruleItem
	}

	listRule(): BLWLRule { // returns either blacklist or whitelist rule
		return <BLWLRule>{
			list_terms: this.terms,
			compare_key: this.comparekeys
		}
	}
}

@Injectable()
export class NhWatcherrulesService {
	constructor() {}

	public createRule(r : RuleBase) {
		let ruleCreate = new RuleCreator(r);

		switch (ruleCreate.ruleBase.rule_type) {
			case "blacklist":
			var blacklist_rule = ruleCreate.listRule();
			ruleCreate.ruleMeta = blacklist_rule;
			break;
			
			default:
			// code...
			break;
		}
		return ruleCreate
	}
}
