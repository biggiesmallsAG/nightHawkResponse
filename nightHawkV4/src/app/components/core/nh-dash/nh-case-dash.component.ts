import { Component, Input, OnInit, OnDestroy, EventEmitter } from '@angular/core';

import { trigger, state, style, animate, transition } from '@angular/animations';
import { Router } from '@angular/router';

import { NhCoreService } from 'app/services/nh-core.service';
import { NhDataExchangeService } from 'app/services/nh-data-exchange.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

import { MaterializeAction } from 'angular2-materialize';
import { AuditHandler } from "../interfaces/audit.interface";

@Component({
	selector: 'app-nh-case-dash',
	templateUrl: './nh-case-dash.component.html',
	styleUrls: ['./nh-case-dash.component.sass'],
	animations: [
	trigger('shiftCaseTree', [
		state('1', style({height: '100%'})),
		state('0', style({display: 'none'}))
		])
	]
})

export class NhCaseDashComponent implements OnInit {
	private isVisible : boolean = true;
	private selectedCase;
	private selectedEndpoint;
	private selectedCaseDate;
	private selectedAudit;
	private sCase;
	private sCaseDate;
	private sEndpoint;
	private sAudit;
	private auditData;
	private totalHits;
	private AuditObject: AuditHandler;
	private pTree: AuditHandler;
	private connError: Object

	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(private _nHCoreSvc:NhCoreService,
		private _nHDataExch:NhDataExchangeService,
		private _nHLoader:NhLoadingService,
		private router:Router) {}

	ngOnInit() {
		this._nHCoreSvc.GET("/list/cases")
		.toPromise()
		.then(cases => {
			this.selectedCase = cases
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();
		});
		this._nHLoader.hide()
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}

	showCaseTree() {
		this.isVisible = !this.isVisible;
	}

	private getCaseEndpoint(_case: string) {
		this.sCase = _case;
		this._nHCoreSvc.GET(`/show/${_case}`)
		.toPromise()
		.then(endpoint => {
			this.selectedEndpoint = endpoint
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();			
		});
	}

	private getCaseDateFromEndpoint(_endpoint: string) {
		this.sEndpoint = _endpoint;
		this.AuditObject = {
			case_name: this.sCase,
			endpoint: this.sEndpoint
		};

		this._nHCoreSvc.GET(`/show
			/${this.AuditObject.case_name}
			/${this.AuditObject.endpoint}`)
		.toPromise()
		.then(response => {
			this.selectedCaseDate = response
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();			
		});
	}

	private getAuditDateFromEndpoint(_case_date: string) {
		this.sCaseDate = _case_date;
		this.AuditObject = {
			case_name: this.sCase,
			endpoint: this.sEndpoint,
			case_date: this.sCaseDate
		};

		this._nHCoreSvc.GET(`/show
			/${this.AuditObject.case_name}
			/${this.AuditObject.endpoint}
			/${this.AuditObject.case_date}`)
		.toPromise()
		.then(response => {
			this.selectedAudit = response
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();			
		});
	}

	private getAuditDataFromAuditType(_audittype: string) {
		for (var i = 0; i < this.selectedAudit.length; i++) {
			if (this.selectedAudit[i].key == _audittype) {
				this.totalHits = this.selectedAudit[i].doc_count
			}
		};

		this.sAudit = _audittype;
		this.AuditObject = {
			case_name: this.sCase,
			case_date: this.sCaseDate,
			endpoint: this.sEndpoint,
			audit_type: this.sAudit,
			total_hits: this.totalHits
		};

		this._nHCoreSvc.GET(`/show
			/${this.AuditObject.case_name}
			/${this.AuditObject.endpoint}
			/${this.AuditObject.case_date}
			/${this.AuditObject.audit_type}
			?from=0&size=100&sort=Record.TlnTime&order=desc`)
		.toPromise()
		.then(audit_data => {
			this.AuditObject.data = audit_data;
			this._nHDataExch.moveAuditData(this.AuditObject);
			if (this.AuditObject.audit_type === "w32processes-tree") {
				this.pTree = this.AuditObject;
				this.modalActions.emit({action:"modal",params:['open']})
			} else {
				this.router.navigate(['core/audits/view', { outlets: {
					auditoutlet: ['auditdash']
				}
			}])	
			}
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();			
		})
	}

	private destroyModalObject() {
		this.modalActions.emit({action:"modal",params:['close']});
		this.pTree = {};
	}
}
