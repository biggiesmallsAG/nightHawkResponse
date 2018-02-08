import { 	Component, Input, OnInit  } from '@angular/core';
import { trigger, state, style, animate, transition } from '@angular/animations';
			
@Component({
	selector: 'app-nh-tree',
	templateUrl: './nh-tree.component.html',
	styleUrls: ['./nh-tree.component.sass'],
	animations: [
	trigger('hideCases', [
		state('1', style({width: '0px'})),
		state('0', style({width: '40px'}))
		])
	]
})
export class NhTreeComponent implements OnInit {
	@Input() isVisible : boolean = true;
	constructor() { }

	ngOnInit() {
	}

}
