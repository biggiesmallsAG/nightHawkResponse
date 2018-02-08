import { Component, Output, EventEmitter, OnInit, OnDestroy	} from '@angular/core';

import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';
import { NhPageTitleService } from 'app/services/nh-page-title.service';
import { NhCoreService } from 'app/services/nh-core.service';
import { NhLoadingService } from 'app/services/nh-loading.service';

import { MaterializeAction } from 'angular2-materialize';

import { UploadHandler } from '../interfaces/upload.interface';
import { DeleteCase, DeleteEndpoint } from '../interfaces/deletecase.interface';

@Component({
	selector: 'app-nh-upload',
	templateUrl: './nh-upload.component.html',
	styleUrls: ['./nh-upload.component.sass']
})
export class NhUploadComponent implements OnInit {

	public pageTitle : string = "Upload & Delete";
	public selectedCases;
	public selectedEndpoints;
	public uploadResponse;
	public deleteResponse;
	public deleteCaseForm: FormGroup;
	public deleteEndpointForm: FormGroup;
	public _deleteEndpointFromCase: boolean = false;
	public deleteObject = {
		delete_case: false, case: "",
		delete_endpoint: false, endpoint: ""
	};
	public uploadForm: FormGroup;
	public submitted: boolean;

	modalActions = new EventEmitter<string|MaterializeAction>();

	constructor(private _pageTitleSvc:NhPageTitleService,
		private _nhCoreSvc:NhCoreService,
		private _nHLoader:NhLoadingService,
		private _fb:FormBuilder) {

	}
	
	ngOnInit() {
		this._pageTitleSvc.updateTitle(this.pageTitle);

		this.uploadForm = this._fb.group({
			files: ['', <any>Validators.required],
			case_name: ['', <any>Validators.required],
			path: ['/upload']
		});

		this.deleteCaseForm = this._fb.group({
			case_name: ['', Validators.required],
			path: ['/delete/case']
		});
		this.deleteEndpointForm = this._fb.group({
			endpoint: ['', Validators.required],
			path: ['/delete/endpoint']
		});

		this._nhCoreSvc.GET("/list/cases")
		.subscribe(cases => this.selectedCases = cases);

		this._nhCoreSvc.GET("/list/endpoints")
		.subscribe(endpoints => this.selectedEndpoints = endpoints);
		this._nHLoader.hide()
	}

	ngOnDestroy() {
		this._nHLoader.show()
	}

	setCase(obj, model: string) {
		switch (model) {
			case "case":
			if (obj != "") {
				this.deleteObject.delete_case = !this.deleteObject.delete_case;
				this.deleteObject.case = obj
			} else {
				this.deleteObject.delete_case = false;
				this.deleteObject.case = "";
			}
			break
			case "endpoint":
			if (obj != "") {
				this.deleteObject.delete_endpoint = !this.deleteObject.delete_endpoint;
				this.deleteObject.endpoint = obj
			} else {
				this.deleteObject.delete_endpoint = false;
				this.deleteObject.endpoint = "";
			}
			break
		}
	}

	deleteCase(model: DeleteCase, isValid: boolean, event: Event) {
		event.preventDefault();
		this._nHLoader.show();

		this._nhCoreSvc.POSTJSON(model.path, {
			case_name: model.case_name
		})
		.toPromise()
		.then(response => {
			this.deleteResponse = response;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();
		})
	}

	deleteEndpoint(model: DeleteEndpoint, isValid: boolean, event: Event) {
		event.preventDefault();
		this._nHLoader.show();
		this._nhCoreSvc.POSTJSON(model.path, {
			endpoint: model.endpoint
		})
		.toPromise()
		.then(response => {
			this.deleteResponse = response;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide()
		})
	}

	deleteEndpointFromCase() {
		this._nHLoader.show();
		this._nhCoreSvc.GET(`/delete/${this.deleteObject.case}/${this.deleteObject.endpoint}`)
		.toPromise()
		.then(response => {
			this.deleteResponse = response;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();
		})
	}

	save(model: UploadHandler, isValid: boolean, event: Event) {
		event.preventDefault();
		this.submitted = true;
		this._nHLoader.show();
		let files: FileList = model.files;
		let formData: FormData = new FormData();

		for(var i = 0; i < files.length; i++){
			formData.append(files[i].name, files[i]);
		}
		formData.append("case_name", model.case_name);
	
		this._nhCoreSvc.POSTUpload(model.path, formData)
		.subscribe(response => {
			this.uploadResponse = response;
			this.modalActions.emit({action:"modal",params:['open']});
			this._nHLoader.hide();
		})
	}

	onChange(event: any) {
		let fileList: FileList = event.target.files;
		if (fileList.length > 0) {
			this.uploadForm.patchValue({files: fileList})
		}
	}
}
