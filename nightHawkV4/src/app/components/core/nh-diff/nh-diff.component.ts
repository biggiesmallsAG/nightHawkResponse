import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { NhCoreService } from 'app/services/nh-core.service';
import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

@Component({
	selector: 'app-nh-diff',
	templateUrl: './nh-diff.component.html',
	styleUrls: ['./nh-diff.component.sass']
})
export class NhDiffComponent implements OnInit {

	constructor(
		private _nHCoreSvc:NhCoreService,
		private _nHLoaderSvc:NhLoadingService,
		private _nHPageTitleSvc:NhPageTitleService) { }

	ngOnInit() {

		this._nHPageTitleSvc.updateTitle("Diffing Framework");
		this._nHLoaderSvc.hide();
	}

	ngOnDestroy() {
		this._nHLoaderSvc.show();
	}
}
