import { Component, OnInit, OnDestroy } from '@angular/core';

import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhCoreService } from 'app/services/nh-core.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

@Component({
	selector: 'app-nh-config-stats',
	templateUrl: './nh-config-stats.component.html',
	styleUrls: ['./nh-config-stats.component.sass']
})
export class NhConfigStatsComponent implements OnInit {

	private pageTitle : string = "Platform Stats";
	private pStats;
	
	constructor(
		private _pageTitleSvc:NhPageTitleService, 
		private _nhCoreSvc:NhCoreService,
		private _nHLoader:NhLoadingService) {
		
		_nhCoreSvc.GET("/platformstats")
		.subscribe(
			statsdata => this.pStats = statsdata,
			error => console.error('Error: ' + error),
			null);
	}

	ngOnInit() {
		this._pageTitleSvc.updateTitle(this.pageTitle);
		this._nHLoader.hide()
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}

	updateStats() {
		this._nHLoader.show()
		this._nhCoreSvc.GET("/platformstats")
		.subscribe(
			statsdata => {this.pStats = statsdata; this._nHLoader.hide()},
			error => console.error('Error: ' + error),
			null);		
	}
}
