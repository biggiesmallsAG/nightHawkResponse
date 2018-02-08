import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { NhCoreService } from 'app/services/nh-core.service';
import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhLoadingService } from 'app/services/nh-loading.service';
import { NhWatcherrulesService } from 'app/services/nh-watcherrules.service';

import { MaterializeAction } from 'angular2-materialize';
import { RuleBase } from '../interfaces/rule.interface';

@Component({
	selector: 'app-nh-watcher',
	templateUrl: './nh-watcher.component.html',
	styleUrls: ['./nh-watcher.component.sass']
})
export class NhWatcherComponent implements OnInit {

	public ruleTypes: Array<Object> = [
		{type: "blacklist", display: "Blacklist"},
		{type: "whitelist", display: "Whitelist"}
	];
	private realertDuration: Array<Object> = [
		{length: "minutes", display: "Minutes"},
		{length: "hours", display: "Hours"},
		{length: "days", display: "Days"}
	];
	public ruleForm: FormGroup;
	private ruleCreator;
	public inView: string;
	private duration: string;

	constructor(
		private _nHLoaderSvc:NhLoadingService,
		private _nHPageTitleSvc:NhPageTitleService,
		private _nHWatcherRulesSvc:NhWatcherrulesService,
		private _nHCoreSvc:NhCoreService,
		private _fb:FormBuilder) { }

	ngOnInit() {

		this.ruleForm = this._fb.group({
			rule_type: ['', Validators.required],
			rule_name: ['', Validators.required], //build custom validator that checks existance of named rule
			compare_key: [''],
			list_terms: [''],
			realert_duration: [''],
			realert_timelength: ['']
		});

		this._nHPageTitleSvc.updateTitle("Watcher Framework");
		this._nHLoaderSvc.hide()
	}

	ngOnDestroy() {
		this._nHLoaderSvc.show();
	}

	ruleTypeInView(rule_type: string) {
		this.inView = rule_type;
	}

	private getRealertDuration(duration: string) {
		this.duration = duration;
	}

	submitRule(model: RuleBase, isValid: boolean, event: Event) {
		event.preventDefault();
		this.ruleCreator = this._nHWatcherRulesSvc.createRule(model)
		console.log(this.ruleCreator)
		this._nHCoreSvc.POSTJSON('/watcher/generate/rule', this.ruleCreator)
		.toPromise()
		.then()
	}
}
