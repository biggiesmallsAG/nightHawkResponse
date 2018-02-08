import { Component, EventEmitter, OnInit, OnDestroy } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhCoreService } from 'app/services/nh-core.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

import { SysConfig } from 'app/services/interfaces/sysconfig.interface';
import { MaterializeAction } from 'angular2-materialize';

@Component({
	selector: 'app-nh-config-sys',
	templateUrl: './nh-config-sys.component.html',
	styleUrls: ['./nh-config-sys.component.sass']
})
export class NhConfigSysComponent implements OnInit {

	private pageTitle : string = "System Configuration";
	private isDisabled: boolean = false;
	private updateResponse;
	private configForm: FormGroup;
	modalActions = new EventEmitter<string|MaterializeAction>();
	configData: SysConfig;

	constructor(
		private _pageTitleSvc:NhPageTitleService, 
		private _nhCoreSvc:NhCoreService,
		private _nHLoader:NhLoadingService,
		private _fb:FormBuilder) {
	}

	ngOnInit() {
		
		this._nhCoreSvc.GET("/config")
		.subscribe(
			confdata => {
				this.configData = confdata;
				this.configForm = this._fb.group({
					nighthawk: this._fb.group({
						ip_addr: [this.configData.nighthawk.ip_addr],
						max_procs: [this.configData.nighthawk.max_procs, <any>Validators.required],
						max_goroutine: [this.configData.nighthawk.max_goroutine, <any>Validators.required],
						bulk_post_size: [this.configData.nighthawk.bulk_post_size, <any>Validators.required],
						opcontrol: [this.configData.nighthawk.opcontrol, <any>Validators.required],
						sessiondir_size: [this.configData.nighthawk.sessiondir_size, <any>Validators.required],
						check_hash: [this.configData.nighthawk.check_hash, <any>Validators.required],
						check_stack: [this.configData.nighthawk.check_stack, <any>Validators.required],
						verbose: [this.configData.nighthawk.verbose, <any>Validators.required],
						verbose_level: [this.configData.nighthawk.verbose_level, <any>Validators.required]		
					}),
					elastic: this._fb.group({
						elastic_server: [this.configData.elastic.elastic_server, <any>Validators.required],
						elastic_port: [this.configData.elastic.elastic_port, <any>Validators.required],
						elastic_user: [this.configData.elastic.elastic_user],
						elastic_pass: [this.configData.elastic.elastic_pass],
						elastic_ssl: [this.configData.elastic.elastic_ssl, <any>Validators.required],
						elastic_index: [this.configData.elastic.elastic_index, <any>Validators.required]
					})
				}) 
			},
			error => console.error('Error: ' + error),
			null);
		this._pageTitleSvc.updateTitle(this.pageTitle);
		this._nHLoader.hide()
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}

	updateConfig(model: SysConfig, isValid: boolean, event: Event) {
		event.preventDefault();

		this._nhCoreSvc.POST("/config", model)
		.subscribe(updateData => {
			this.updateResponse = updateData;
			this.modalActions.emit({action:"modal",params:['open']});
		});
	}

}
