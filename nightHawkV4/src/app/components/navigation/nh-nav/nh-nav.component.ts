import { Component, EventEmitter, Output, Input, OnInit } from '@angular/core';
import {Subscription} from 'rxjs/Subscription';

import { NhPageTitleService } from '../../../services/nh-page-title.service';

@Component({
	selector: 'app-nh-nav',
	templateUrl: './nh-nav.component.html',
	styleUrls: ['./nh-nav.component.sass']
})
export class NhNavComponent implements OnInit {

	private showCases : boolean = true;
	private pageTitle : string;

	@Output() shiftCallback: EventEmitter<boolean> = new EventEmitter<boolean>();
	subscription:Subscription;

	constructor(private _pageTitleSvc:NhPageTitleService) {}

	ngOnInit() {
		this.subscription = this._pageTitleSvc.navItem$
		.subscribe(pageTitle => this.pageTitle = pageTitle);
	}

	ngDestroy() {
		this.subscription.unsubscribe();
	}

	showHideCases() {
		this.showCases = !this.showCases;
		this.shiftCallback.next(true);
	}
}
