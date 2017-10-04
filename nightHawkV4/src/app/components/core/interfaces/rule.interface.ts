export interface RuleBase extends BLWLRule {
	rule_type: string,
	rule_name: string,
	realert_duration: number,
	realert_timelength: string,
	rule_meta: BLWLRule
}

export interface BLWLRule {
	list_terms: any,
	compare_key: any
}