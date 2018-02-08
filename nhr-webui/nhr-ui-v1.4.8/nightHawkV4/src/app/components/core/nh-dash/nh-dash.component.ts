import { Component, OnInit, OnDestroy, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { NhCoreService } from 'app/services/nh-core.service';
import { NhGridHelperService } from 'app/services/nh-grid-helper.service';
import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhLoadingService } from 'app/services/nh-loading.service';
import { GridOptions } from 'ag-grid'
import { StackHandler } from '../interfaces/stack.interface';

import { MaterializeAction } from 'angular2-materialize';

@Component({
	selector: 'app-nh-dash',
	templateUrl: './nh-dash.component.html',
	styleUrls: ['./nh-dash.component.sass']
})
export class NhDashComponent implements OnInit {

	private columDefs: Array<Object> = [
	{headerName: "Alert", sort: "desc", field: "alert_time", cellRenderer: this.cellRenderer},
	{headerName: "Rule Name", field: "rule_name", cellRenderer: this.cellRenderer},
	{headerName: "Matches", field: "match_body.num_matches", cellRenderer: this.cellRenderer},
	{headerName: "Match Doc", field: "match_body._id", cellRenderer: this.cellRenderer},
	{headerName: "Case Name", field: "match_body.CaseInfo.case_name", cellRenderer: this.cellRenderer},
	{headerName: "Endpoint", field: "match_body.CaseInfo.computer_name", cellRenderer: this.cellRenderer},
	{headerName: "Case Date", field: "match_body.CaseInfo.case_date", cellRenderer: this.cellRenderer},
	{headerName: "Case Analyst", field: "match_body.CaseInfo.case_analyst", cellRenderer: this.cellRenderer}
	];
	private gridOptions: GridOptions;
	private selectedRow;
	public watcherResults;
	public resultDoc: Object;
	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(
		private _nHLoader:NhLoadingService,
		private _nHCoreSvc:NhCoreService,
		private _nHPageTitleSvc:NhPageTitleService,
		private _nHGridHelperSvc:NhGridHelperService) { }

	ngOnInit() {
		this._nHPageTitleSvc.updateTitle("Main Dashboard");
		this._nHCoreSvc.GET("/watcher/results")
		.toPromise()
		.then(response => {
			this.watcherResults = response;
			this.gridOptions = <GridOptions>{
				rowHeight: 45,
				headerHeight: 45,
				columnDefs: this.columDefs,
				rowData: this.watcherResults,
				rowSelection: "single",
				headerCellRenderer: this.headerRenderer,
				onGridReady: () => {
					this.gridOptions.columnApi.autoSizeAllColumns()
				}
			};
			this._nHLoader.hide();
		});
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}

	private getWatcherMatchById() {
		this._nHLoader.show();
		this.selectedRow = this.gridOptions.api.getSelectedRows();
		const _id = this.selectedRow[0].match_body._id
		this._nHCoreSvc.GET(`/watcher/results/${_id}`)
		.toPromise()
		.then(response => {
			this.resultDoc = response;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();
		})
	}

	destroyModalObject() {
		this.modalActions.emit({action:"modal",params:['close']});
		this.resultDoc = "";
	}

	private headerRenderer = (params) => {
		return '<div class="btn-flat greytext">' + params.value + '</div>'
	}

	private cellRenderer(params) {
		switch (params.column.colId) {
			case "alert_time":
			return '<div class="btn-flat nhblue"> \
			'+ params.value +'</div>'
			case "rule_name":
			return '<div class="btn-flat nhorange"> \
			'+ params.value +'</div>'
			case "match_body.num_matches":
			return '<div class="btn-flat limetext"> \
			'+ params.value +'</div>'
			case "match_body._id":
			return '<div class="btn-flat red"> \
			'+ params.value +'</div>'
			case "match_body.CaseInfo.case_name":
			return '<div class="btn-flat greytext"> \
			'+ params.value +'</div>'
			case "match_body.CaseInfo.computer_name":
			return '<div class="btn-flat greytext"> \
			'+ params.value +'</div>'
			case "match_body.CaseInfo.case_date":
			return '<div class="btn-flat nhblue"> \
			'+ params.value +'</div>'
			case "match_body.CaseInfo.case_analyst":
			if (typeof(params.value) == 'undefined') {
				return ''
			} else {
				return '<div class="btn-flat greytext"> \
				'+ params.value +'</div>'
			}
		}
	}
}
