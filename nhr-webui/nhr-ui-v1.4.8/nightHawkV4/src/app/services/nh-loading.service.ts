import { Injectable } from '@angular/core';
import { Observable, Subject, BehaviorSubject } from 'rxjs';

import { LoaderState } from 'app/components/core/utility/nh-loader/nh-loader.component';

@Injectable()
export class NhLoadingService {

	private loaderSubject = new Subject<LoaderState>();
	loaderState = this.loaderSubject
	.asObservable();

	private counter: number = 0;
	
	constructor() { }

	//emits true to show spinner on the state observable in the component class
	public show(): void {
		this.counter++;
		this.loaderSubject.next(<LoaderState>{show: true});
	}

	//emits false to show spinner on the state observable in the component class
	public hide(): void {
		if(this.counter >= 0 ){
			this.counter--;
		}
		if(this.counter <= 0){
			this.loaderSubject.next(<LoaderState>{show: false});
		}
	}

	//emits false to show spinner on the state observable in the component class
	public forceHide(): void {
		this.counter = 0;
		this.loaderSubject.next(<LoaderState>{show: false});
	}

	public getCount(): number {
		return this.counter;
	}

}
