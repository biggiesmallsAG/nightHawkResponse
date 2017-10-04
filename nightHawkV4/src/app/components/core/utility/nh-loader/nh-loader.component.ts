import { Component, OnInit, OnDestroy, trigger, state, style, transition, animate } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';

import { NhLoadingService } from 'app/services/nh-loading.service';

export interface LoaderState {
	show: boolean;
}

@Component({
	selector: 'app-loading-indicator',
	template: `
		<div *ngIf="show">
			<div class="loading" [@loading]="true">
				<div class="center-align">
					<div class="btn btn-floating btn-large pulse"><i class="material-icons">query_builder</i></div>
				</div>
			</div>
		</div>
	`,
	styleUrls: ['./nh-loader.component.sass'],
	animations: [
		trigger('loading', [
			state('in', style({ opacity: 1 })),
			transition(':enter', [
				style({ opacity: 0 }),
				animate('500ms ease-out')
			]),
			transition(':leave', [
				animate('200ms ease-in', style({ opacity: 1 }))
			])
		])
	]
})
export class NhLoaderComponent implements OnInit {

	show: boolean = false;
	private subscription: Subscription;

	constructor(private loaderService: NhLoadingService) { }

	ngOnInit() {
		//observable that will hold the state of the spinner
		this.subscription = this.loaderService.loaderState
			.subscribe((state: LoaderState) => {
				this.show = state.show;
			});
	}

	ngOnDestroy() {
		this.subscription.unsubscribe();
	}

}