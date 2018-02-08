import { Component, OnInit, OnDestroy } from '@angular/core';

import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

@Component({
	selector: 'app-nh-config',
	templateUrl: './nh-config.component.html',
	styleUrls: ['./nh-config.component.sass']
})
export class NhConfigComponent implements OnInit {

	private pageTitle : string = "Configuration";
	
	constructor(
		private _pageTitleSvc:NhPageTitleService,
		private _nHLoader:NhLoadingService) {}

	ngOnInit() {
		this._pageTitleSvc.updateTitle(this.pageTitle);
		this._nHLoader.hide()
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}
}
