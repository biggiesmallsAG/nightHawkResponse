import { Component, OnInit, OnDestroy, ViewEncapsulation, EventEmitter } from '@angular/core';

import { NhDataExchangeService } from 'app/services/nh-data-exchange.service';
import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhGridHelperService } from 'app/services/nh-grid-helper.service';
import { NhCoreService } from 'app/services/nh-core.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

import { Subscription } from 'rxjs/Subscription';
import { AuditHandler } from "../interfaces/audit.interface";
import { GridOptions } from 'ag-grid';
import { MaterializeAction } from 'angular2-materialize';

import { CustomDateComponent } from 'app/components/core/nh-audits/custom-renderers/customdate.component';
import { Router } from '@angular/router';

@Component({
	selector: 'app-nh-audit-outlet',
	templateUrl: './nh-audit-outlet.component.html',
	styleUrls: ['./nh-audit-outlet.component.sass'],
	encapsulation: ViewEncapsulation.None
})
export class NhAuditOutletComponent implements OnInit {

	public auditData: any = [];
	subscription1: Subscription;
	private auditObject: AuditHandler;
	public gridOptions: GridOptions;
	public colApi;
	public sourceRecord: Object;
	private selectedRows;
	public showGrid: boolean;
	public docId: string;
	public frameworkComponents;
	private columnDefs = [];
	private colAuditType = [];
	private colCaseInfo = [];
	private colRecord = [];
	private colHeaders = [];
	private pageTitle: string = "";
	private gridInit: boolean = true;
	private auditChange: boolean = false;
	public connError: Object;
	private AO;
	private lastRow = 100;
	private sortModel: string = "";
	private sortOrder: string = "";

	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(
		private _nHDataExch:NhDataExchangeService,
		private _nHPageTitleSvc:NhPageTitleService,
		private _nHGridHelper:NhGridHelperService,
		private _nHCoreSvc:NhCoreService,
		private _nHLoader:NhLoadingService,
		private _router:Router) {
		this.showGrid = true;
	}

	ngOnInit() {
		this.subscription1 = this._nHDataExch.dO$
		.subscribe(_auditdata => {
			this.auditData = [];
			this.auditObject = _auditdata;
			this.frameworkComponents = { agDateInput: CustomDateComponent };
			try {
				for (var i in this.auditObject.data) {
					this.auditData.push(this.auditObject.data[i]);
				};
	
				var _case = this.auditObject.case_name;
				var _audit = this.auditObject.audit_type;
	
				this.pageTitle = `${_case} (${_audit})`;
				this._nHPageTitleSvc.updateTitle(this.pageTitle);
	
				if (!this.gridInit) {
					this.auditChange = true;
					this.createColumnDefs();
					this.gridOptions.columnDefs = this.columnDefs;
					this.gridOptions.api.setDatasource(this.SourceGridData());
				} else {
					// Init ColDefs
					if (this.auditData.length > 0) {
						this.createColumnDefs();
						this.gridOptions = <GridOptions>{
							columnDefs: this.columnDefs,
							rowModelType: "infinite",
							rowSelection: "single",
							maxBlocksInCache: 5,
							enableServerSideFilter: true,
							enableServerSideSorting: true,
							datasource: this.SourceGridData(),
							onGridReady: (params) => {
								this.colApi = params.columnApi
							}
						};
					};
				};
			}
			catch (error) {
				this._router.navigate([''])
			}
		});
		this._nHLoader.hide()
	}

	ngOnDestroy() {
		this.subscription1.unsubscribe();
		this._nHLoader.show()
	}

	autoSizeAll() {
		var allColumnIds = [];
		this.colApi.getAllColumns().forEach(function(column) {
		  allColumnIds.push(column.colId);
		});
		this.colApi.autoSizeColumns(allColumnIds);
	}

	private onAuditSelectionChanged() {
		this._nHLoader.show();
		this.selectedRows = this.gridOptions.api.getSelectedRows();
		this.docId = this.selectedRows[0]._id;
		this._nHCoreSvc.GET(`/show/doc/${this.docId}/${this.auditObject.endpoint}`)
		.toPromise()
		.then(response => {
			this.sourceRecord = response;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();
		},
		error => {
			this.connError = error;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();
		})
	}

	destroyModalObject() {
		this.modalActions.emit({action:"modal",params:['close']});
		this.sourceRecord = '';
	}

	private createColumnDefs() {
		this.columnDefs = [];
		this.columnDefs = this._nHGridHelper.iterate(this.auditData[0]);
	};

	private SourceGridData() {
		// We must be able to tell the function how many rows it should be expecting
		// and make a call to nHCoreSvc to update the size and sort if needed.

		this.lastRow = this.auditObject.total_hits;
		return {
			rowCount: null,
			getRows: params => {
				if (this.gridInit) {
					this.gridInit = false;
					params.successCallback(this.auditData, this.lastRow) // lastRow must be set via a calculation of the total rows avaiable.
				} else if (this.auditChange) {
					this.auditChange = false;
					params.successCallback(this.auditData, this.lastRow) // check from Observable for change in audit
				} else {
					try {
						var sourceRex = new RegExp("_source\.(.*)");
						var rowHeader = sourceRex.exec(params.sortModel[0].colId);
						
						switch (rowHeader[1]) {
							// default fallthrough for dates/integers
							case "Record.TlnTime":
							case "Record.SizeInBytes":
							case "Record.ReportedSizeInBytes":
							case "Record.TimesExecuted":
							case "Record.Created":
							case "Record.LastRun":
							case "Record.LastVisitDate":
							case "Record.BytesDownloaded":
							case "Record.LastModifiedDate":
							case "Record.MaxBytes":
							case "Record.CacheHitCount":
							case "Record.LastCheckedDate":
							case "Record.VisitCount":
							case "Record.Modified":
							case "Record.ReportedLengthInBytes":
							case "Record.RegModified":
							case "Record.FileCreated":
							case "Record.FileModified":
							case "Record.FileAccessed":
							case "Record.FileChanged":
							case "Record.Registry.JobCreated":
							case "Record.Registry.TlnTime":
							case "Record.Registry.Modified":
							case "Record.Registry.ReportedLengthInBytes":
							case "Record.File.JobCreated":
							case "Record.File.TlnTime":
							case "Record.File.SizeInBytes":
							case "Record.File.Created":
							case "Record.File.Modified":
							case "Record.File.Accessed":
							case "Record.File.Changed":
							case "Record.PeInfo.PETimeStamp":
							case "Record.Pid":
							case "Record.Index":
							case "Record.EID":
							case "Record.GenTime":
							case "Record.WriteTime":
							case "Record.LocalPort":
							case "Record.RemotePort":
							case "Record.CreationDate":
							case "Record.MostRecentRunTime":
							case "Record.NextRunTime":
							case "Record.DataLength":
							this.sortModel = rowHeader[1];
							break;

							default:
							this.sortModel = rowHeader[1] + '.keyword';
							break;
						};
						this.sortOrder = params.sortModel[0].sort;
					} catch (e) {
						this.sortModel = "Record.TlnTime";
						this.sortOrder = "desc";
					};
					if (Object.keys(params.filterModel).length != 0) { // we need to get the filterModel if it exists and POST it.
						for (var k in params.filterModel) {
							var colId = sourceRex.exec(Object.keys(params.filterModel)[0])
							var filterOn = params.filterModel[k];
							// using POSTUpload to get full response data struct back other than pure data.
							this._nHCoreSvc.POSTUpload(`/show
								/${this.auditObject.case_name}
								/${this.auditObject.endpoint}
								/${this.auditObject.case_date}
								/${this.auditObject.audit_type}
								?from=${params.startRow}&size=100&sort=${this.sortModel}&order=${this.sortOrder}`, {
									colId: colId[1],
									filterOn: filterOn
								})
							.toPromise()
							.then(_auditData => {
								params.successCallback(_auditData.data, _auditData.total)
							}, 
							error => {
								this.connError = error;
								this.modalActions.emit({action:"modal",params:['open']});
								this._nHLoader.hide();
							})
						}
					} else {
						this._nHCoreSvc.GET(`/show
							/${this.auditObject.case_name}
							/${this.auditObject.endpoint}
							/${this.auditObject.case_date}
							/${this.auditObject.audit_type}
							?from=${params.startRow}&size=100&sort=${this.sortModel}&order=${this.sortOrder}`)
						.toPromise()
						.then(_auditData => {
							this.auditData = _auditData;
							params.successCallback(_auditData, this.lastRow)
						}, 
						error => {
							this.connError = error;
							this.modalActions.emit({action:"modal",params:['open']});
							this._nHLoader.hide();
						});
					}
				}
			}
		}
	}
}
