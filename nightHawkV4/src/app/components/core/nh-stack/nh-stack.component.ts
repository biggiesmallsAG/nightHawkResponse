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
	selector: 'app-nh-stack',
	templateUrl: './nh-stack.component.html',
	styleUrls: ['./nh-stack.component.sass']
})
export class NhStackComponent implements OnInit {

	private searchSubmitted: boolean = false;
	private filterGood: boolean = false;
	private gridInit: boolean = true;
	private ctxgridInit: boolean = true;
	private stackResponse: Array<Object>;
	private contextResponse: Array<Object>;
	private ctxgridOptions: GridOptions;
	private gridOptions: GridOptions;
	private contextItem: string = '';
	private sourceRecord: Object;
	private connError: Object;
	private ctxcolDefs;
	private columnDefs;
	private selectedRows: Array<any>;
	private ctxselectedRows: Array<any>;
	private stackableObjects: Array<Object> = [
	{name: "Services", value: "service"},
	{name: "Prefetch", value: "prefetch"},
	{name: "Scheduled Tasks", value: "task"},
	{name: "DNS", value: "dns/a"},
	{name: "Local Ports", value: "locallistenport"},
	{name: "Persistence", value: "persistence"},
	{name: "Run Keys", value:"runkey"}
	];

	private stackForm: FormGroup;
	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(
		private _NhCoreSvc:NhCoreService,
		private _NhGridHelper:NhGridHelperService,
		private _NhPageTitleSvc:NhPageTitleService,
		private _NhLoader:NhLoadingService,
		private _fb:FormBuilder) {}

	ngOnInit() {
		this.stackForm = this._fb.group({
			stack_type: ['', Validators.required],
			search_limit: [''],
			sort_desc: [false],
			ignore_good_service: [''],
			path: ['/stacking/']
		});
		this._NhPageTitleSvc.updateTitle("Stacking Framework");
		this._NhLoader.hide()
	}

	ngOnDestroy() {
		this._NhLoader.show();
	}

	private createColumnDefs() {
		this.columnDefs = [];
		this.columnDefs = this._NhGridHelper.iterate(this.stackResponse[0], "stack");
	};

	private stackType(type: string) {
		type == "service" ? this.filterGood = !this.filterGood : this.filterGood = false
	}

	private onEndpointSelectionChanged(model: StackHandler) {
		this.ctxselectedRows = this.ctxgridOptions.api.getSelectedRows();

		this._NhCoreSvc.POSTJSON("/stacking/context/endpoint", {
			type: model.stack_type,
			context_item: this.contextItem,
			endpoint: this.ctxselectedRows[0].endpoint_name
		})
		.toPromise()
		.then(response => {
			this.sourceRecord = response[0];
			this.modalActions.emit({action:"modal",params:['open']});
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._NhLoader.hide();
		})
	}

	private onStackSelectionChanged(model: StackHandler) {
		this.selectedRows = this.gridOptions.api.getSelectedRows();

		switch (model.stack_type) {
			case "service":
			this.contextItem = this.selectedRows[0].ServicePath
			break;
			case "prefetch":
			this.contextItem = this.selectedRows[0].AppFullPath
			break;
			case "persistence":
			this.contextItem = this.selectedRows[0].Path
			break;
			case "runkey":
			this.contextItem = this.selectedRows[0].RegKey
			break;
			case "task":
			this.contextItem = this.selectedRows[0].TaskName
			break;
			case "dns/a":
			this.contextItem = this.selectedRows[0].IpAddress
			break;
			default:
			break;
		}
		this._NhCoreSvc.POSTJSON("/stacking/context", {type: model.stack_type, context_item: this.contextItem})
		.toPromise()
		.then(response => {
			this.contextResponse = response;
			this.ctxcolDefs = this._NhGridHelper.iterate(this.contextResponse[0], "stack")
			if (!this.ctxgridInit) {
				this.ctxgridOptions.api.setColumnDefs(this.ctxcolDefs);
				this.ctxgridOptions.api.setRowData(this.contextResponse);
				this.ctxgridOptions.columnApi.autoSizeAllColumns();
			} else {
				this.ctxgridInit = false;
				this.ctxgridOptions = <GridOptions>{
					columnDefs: this.ctxcolDefs,
					rowSelection: "single",
					onGridReady: () => {
						this.ctxgridOptions.api.setRowData(this.contextResponse);
						this.ctxgridOptions.columnApi.autoSizeAllColumns()
					}
				};
			};		
		})
	}

	private getStack(model: StackHandler, isValid: boolean, event: Event) {
		event.preventDefault();
		this._NhLoader.show();
		this.searchSubmitted = true;
		this.selectedRows = null;

		this._NhCoreSvc.POSTJSON(model.path + model.stack_type, model)
		.toPromise()
		.then(response => {
			this.stackResponse = response;
			this.searchSubmitted = false;
			this.createColumnDefs();

			if (!this.gridInit) {
				this.gridOptions.api.setColumnDefs(this.columnDefs);
				this.gridOptions.api.setRowData(this.stackResponse);
				this.gridOptions.columnApi.autoSizeAllColumns();
				this._NhLoader.hide();
			} else {
				
				this.gridInit = false;
				this.gridOptions = <GridOptions>{
					columnDefs: this.columnDefs,
					rowSelection: "single",
					onGridReady: () => {
						this.gridOptions.api.setRowData(this.stackResponse);
						this.gridOptions.columnApi.autoSizeAllColumns()
					}
				};
				this._NhLoader.hide();
			};						
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._NhLoader.hide();
		});

	}

	private destroyModalObject() {
		this.modalActions.emit({action:"modal",params:['close']});
		this.sourceRecord = '';
	}
}
