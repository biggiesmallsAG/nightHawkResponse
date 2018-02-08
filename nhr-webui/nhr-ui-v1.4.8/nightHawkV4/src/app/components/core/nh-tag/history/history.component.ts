import { Component, OnInit, Input } from '@angular/core';
import { GridOptions } from 'ag-grid';
import { NhCoreService } from 'app/services/nh-core.service';
import { NhGridHelperService } from 'app/services/nh-grid-helper.service';
import { AuditHandler } from 'app/components/core/interfaces/audit.interface';

@Component({
  selector: 'app-tag-history',
  templateUrl: './history.component.html',
  styleUrls: ['./history.component.sass']
})
export class HistoryComponent implements OnInit {
  @Input() tagMeta: AuditHandler;
  @Input() docId: string;
  public gridOptions: GridOptions;
  private columnDefs: Object[] = [];
  public tagData: Object[] = [];
  constructor(
    private _NhCoreSvc: NhCoreService,
    private _NhGridHelperSvc: NhGridHelperService
  ) { }

  ngOnInit() {
    this._NhCoreSvc.GET(`/tag/show/${this.tagMeta.case_name}/${this.tagMeta.endpoint}/${this.tagMeta.audit_type}/${this.docId}`)
      .toPromise()
      .then(
      (response) => {
        this.tagData = response
      }
      )
    this.columnDefs = [
      {
        headerName: "Timestamp",
        field: "Timestamp"
      },
      {
        headerName: "Tag name",
        field: "TagName"
      },
      {
        headerName: "Created By",
        field: "CreatedBy"
      },
      {
        headerName: "Audit Type",
        field: "Audit"
      }
    ];

    this.gridOptions = <GridOptions>{
      columnDefs: this.columnDefs,
      onGridReady: () => {
        if (this.tagData != null) {
          this.gridOptions.api.setRowData(this.tagData)
        }
      }
    }; 
  }

}
