import { Component, OnInit } from '@angular/core';

import { NhPageTitleService } from 'app/services/nh-page-title.service';

@Component({
	selector: 'app-nh-audits',
	templateUrl: './nh-audits.component.html',
	styleUrls: ['./nh-audits.component.sass']
})
export class NhAuditsComponent implements OnInit {

	private pageTitle : string = "Case & Audits";

	constructor(private _pageTitleSvc:NhPageTitleService) {}
	ngOnInit() {
		this._pageTitleSvc.updateTitle(this.pageTitle);
	}

}
