import { Component, OnInit, OnDestroy, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { NhCoreService } from 'app/services/nh-core.service';
import { NhGridHelperService } from 'app/services/nh-grid-helper.service';
import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhLoadingService } from 'app/services/nh-loading.service';
import { NhValidatorService } from 'app/services/nh-validator.service';
import { GridOptions } from 'ag-grid'

import { TimelineHandler } from '../interfaces/search.interface';

import { MaterializeAction } from 'angular2-materialize';

@Component({
	selector: 'app-nh-timeline',
	templateUrl: './nh-timeline.component.html',
	styleUrls: ['./nh-timeline.component.sass']
})
export class NhTimelineComponent implements OnInit {

	private timelineForm: FormGroup;
	private selectedCase: Array<Object>;
	private selectedEndpoint: Array<Object>;
	private gridOptions: GridOptions;
	private gridInit: boolean = true;
	private emptyResponse: boolean = false;
	private columnDefs;
	private timelineResponse;
	private connError: Object;

	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(
		private _nHLoader:NhLoadingService, 
		private _nHCoreSvc:NhCoreService,
		private _nHPageTitleSvc:NhPageTitleService,
		private _nHGridHelper:NhGridHelperService,
		private _fb:FormBuilder) { }

	ngOnInit() {
		this._nHCoreSvc.GET("/list/cases")
		.toPromise()
		.then(cases => {
			this.selectedCase = cases
		}, 
		error => {
			this.connError = error;
		});

		this.timelineForm = this._fb.group({
			endpoint: ['', Validators.required],
			case: ['', Validators.required],
			start_time: ['', NhValidatorService.timestampValidity],
			end_time: ['', NhValidatorService.timestampValidity],
			time_delta: [''],
			search_limit: [500],
			ignore_good: [false],
			path: ['/search/timeline']
		});

		this._nHPageTitleSvc.updateTitle("Timeline Framework")
		this._nHLoader.hide()
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}

	private getCaseEndpoint(_case: string) {
		this._nHCoreSvc.GET(`/show/${_case}`)
		.toPromise()
		.then(endpoint => {
			this.selectedEndpoint = endpoint
		}, 
		error => {
			this.connError = error
		});
	}

	private parsedate(date, time: string) {
		return new Date(date + ' ' + time)
	}

	private searchTimeline(model: TimelineHandler, isValid: boolean, event: Event) {
		event.preventDefault()
		this._nHLoader.show();
		
		var endpoint_list = [];

		endpoint_list.push(model.endpoint);
		model.endpoint_list = endpoint_list 

		this._nHCoreSvc.POSTJSON(model.path, model)
		.toPromise()
		.then(response => {
			this.timelineResponse = response;
			this.columnDefs = this._nHGridHelper.iterate(this.timelineResponse[0], "timeline");
			try {
				if (this.timelineResponse.length < 1) {
					throw Error;
				} else {
					this.emptyResponse = false;
				}
				if (!this.gridInit) {
					this.gridOptions.api.setColumnDefs(this.columnDefs);
					this.gridOptions.api.setRowData(this.timelineResponse);
					this.gridOptions.columnApi.autoSizeAllColumns();
				} else {
					
					this.gridInit = false;
					this.gridOptions = <GridOptions>{
						columnDefs: this.columnDefs,
						rowSelection: "single",
						onGridReady: () => {
							this.gridOptions.api.setRowData(this.timelineResponse);
							this.gridOptions.columnApi.autoSizeAllColumns()
						}
					};
				};		
			} catch (e) {
				this.emptyResponse = true
			}
			this._nHLoader.hide();
		}, 
		error => {
			this.connError = error
		});
	}
}
