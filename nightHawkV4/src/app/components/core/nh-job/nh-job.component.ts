import { Component, OnInit, OnDestroy, EventEmitter } from '@angular/core';

import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhWebsocketService } from 'app/services/nh-websocket.service';
import { NhCoreService } from 'app/services/nh-core.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

import { MaterializeAction } from 'angular2-materialize';
import { JobHandler } from '../interfaces/job.interface';

@Component({
	selector: 'app-nh-job',
	templateUrl: './nh-job.component.html',
	styleUrls: ['./nh-job.component.sass']
})
export class NhJobComponent implements OnInit {

	private pageTitle : string = "Job Tracking";
	private isDisabled: boolean = false;
	private job;
	private jobarray = [];
	private completed;

	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(
		private _pageTitleSvc:NhPageTitleService, 
		private _nhWebSocketSvc:NhWebsocketService,
		private _nhCoreSvc:NhCoreService,
		private _nHLoader:NhLoadingService) {
		
		_nhWebSocketSvc.connect('ws://localhost:8080/api/v1/subscribe/uploadjobs')
		.subscribe(response => { 
			this.job = JSON.parse(response.data);  
			this.manageJobs(this.job);
		});

	}

	ngOnInit() {
		this._pageTitleSvc.updateTitle(this.pageTitle);
		this._nhCoreSvc.GET("/list/completedjobs")
			.subscribe(response => {this.completed = response; this._nHLoader.hide()})
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}

	manageJobs(job: Object) {
		let _job = <JobHandler>job;

		if (this.jobarray.length < 1 || !(_job.body.is_complete)) {
			this.jobarray.push(_job);
		} else {
			for (var i = 0; i < this.jobarray.length; i++) {
				if (_job.body.uid === this.jobarray[i].body.uid) {
					if (!_job.body.in_progress) {
						this.jobarray.splice(i)
					}
				} else {
					if (_job.body.in_progress) {
						this.jobarray.push(_job)
					}
				}
			}
		}
	}
}
