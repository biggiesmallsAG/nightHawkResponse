import { 	Component,
			OnInit,
			Input,
			trigger, 
			style, 
			state } from '@angular/core';


import { NhPageTitleService } from 'app/services/nh-page-title.service';

@Component({
  	selector: 'nighthawk',
  	templateUrl: './nh-container.component.html',
  	styleUrls: ['./nh-container.component.sass'],
	animations: [
	trigger('shiftLeft', [
		state('1', style({marginLeft: '0'})),
		state('0', style({marginLeft: '40px'}))
		])
	]
})
export class NhContainerComponent implements OnInit {

	private isVisible : boolean = true;
	private pageTitle : string = "Platform Dashboard";
	
	constructor(private _pageTitleSvc:NhPageTitleService) {}
	ngOnInit() {
		this._pageTitleSvc.updateTitle(this.pageTitle);
	}

	pageShift(active: boolean) {
		this.isVisible = !this.isVisible;
	}

}
