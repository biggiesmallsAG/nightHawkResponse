import { Component, Input } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';

import { NhValidatorService } from 'app/services/nh-validator.service';

@Component({
	selector: 'nh-form-control',
	template: `<div class="alert-danger" *ngIf="errorMessage !== null">
	<strong class="strong-color">Error:</strong> {{errorMessage}}</div>`,
	styleUrls: ['./nh-form-control.component.sass']
})
export class NhFormControlComponent {
	@Input() control: FormControl
	constructor() { }
	
	get errorMessage() {
		for (let propertyName in this.control.errors) {
			if (this.control.errors.hasOwnProperty(propertyName) && this.control.touched) {
				return NhValidatorService.getValidatorErrorMessage(propertyName, this.control.errors[propertyName])	
			}
		}

		return null;
	}
}
