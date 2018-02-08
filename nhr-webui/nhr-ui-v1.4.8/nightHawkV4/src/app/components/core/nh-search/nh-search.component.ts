import { Component, OnInit, OnDestroy, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhLoadingService } from 'app/services/nh-loading.service';
import { NhCoreService } from 'app/services/nh-core.service';
import { SearchHandler } from '../interfaces/search.interface';
import { NhGridHelperService } from 'app/services/nh-grid-helper.service';
import { GridOptions } from 'ag-grid'
import { MaterializeAction } from 'angular2-materialize';

@Component({
	selector: 'app-nh-search',
	templateUrl: './nh-search.component.html',
	styleUrls: ['./nh-search.component.sass']
})
export class NhSearchComponent implements OnInit {

	public searchForm: FormGroup
	public searchSubmitted: boolean = false;
	private gridInit: boolean = true;
	private gridOptions: GridOptions;
	public sourceRecord: Object;
	public auditObject: Object;
	public docId: string;
	private selectedRows;
	public searchResponse: Array<Object>;
	public emptyResponse: boolean = false;
	private columnDefs;
	public connError: Object;
	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(
		private _NhPagetitleSvc:NhPageTitleService,
		private _NhCoreSvc:NhCoreService,
		private _NhLoader:NhLoadingService,
		private _NhGridHelper:NhGridHelperService,
		private _fb:FormBuilder) { }

	ngOnInit() {
		this.searchForm = this._fb.group({
			search_term: ['', Validators.required],
			search_size: [500],
			path: ['/search']
		});
		this._NhPagetitleSvc.updateTitle('Search Framework');
		this._NhLoader.hide();
	}

	ngOnDestroy() {
		this._NhLoader.show()
	}

	searchGlobal(model: SearchHandler, isValid: boolean, event: Event) {
		this.emptyResponse = false;
		event.preventDefault();
		this.searchSubmitted = true;
		this._NhLoader.show();

		this._NhCoreSvc.POSTJSON(model.path, {
			search_term: model.search_term,
			search_size: model.search_size
		})
		.toPromise()
		.then(response => {
			this.searchResponse = response;
			try {
				if (this.searchResponse.length < 1) {
					throw Error;
				};
				this.columnDefs = this._NhGridHelper.iterate(this.searchResponse[0], "search")
				if (!this.gridInit) {
					this.gridOptions.api.setColumnDefs(this.columnDefs);
					this.gridOptions.api.setRowData(this.searchResponse);
					this.gridOptions.columnApi.autoSizeAllColumns();
				} else {
					this.gridInit = false;
					this.gridOptions = <GridOptions>{
						columnDefs: this.columnDefs,
						rowSelection: "single",
						onGridReady: () => {
							this.gridOptions.api.setRowData(this.searchResponse);
							this.gridOptions.columnApi.autoSizeAllColumns()
						}
					};
				};		
			} catch (e) {
				this.emptyResponse = true
			};
			this.searchSubmitted = false;
			this._NhLoader.hide()
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this.searchSubmitted = false;
			this._NhLoader.hide();
		})
	}

	private onSearchChange() {
		this._NhLoader.show();
		this.selectedRows = this.gridOptions.api.getSelectedRows();
		this.docId = this.selectedRows[0].id;
		const parent = this.selectedRows[0].computer_name;
		const case_name = this.selectedRows[0].case_name;
		const audit_generator = this.selectedRows[0].audit_generator;

		this.auditObject = {
			case_name: case_name,
			endpoint: parent,
			audit_type: audit_generator
		};
		this._NhCoreSvc.GET(`/show/doc/${this.docId}/${parent}`)
		.toPromise()
		.then(response => {
			this.sourceRecord = response;
			this.modalActions.emit({action:"modal",params:['open']});
			this._NhLoader.hide();
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this.searchSubmitted = false;
			this._NhLoader.hide();			
		})
	}

	destroyModalObject() {
		this.modalActions.emit({action:"modal",params:['close']});
		this.sourceRecord = '';
	}

}
