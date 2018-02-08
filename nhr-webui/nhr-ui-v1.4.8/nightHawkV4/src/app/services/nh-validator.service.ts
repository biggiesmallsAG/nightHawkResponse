import { Injectable } from '@angular/core';

@Injectable()
export class NhValidatorService {

	constructor() { }
	
	static getValidatorErrorMessage(validatorName: string, validatorValue?: any) {
		let config = {
			'required': 'Required',
			'missingCodeName': 'Codename is required',
			'invalidTimestampStruct': 'Timestamp must be in format "2017-03-03T13:12:00Z"',
			'minlength': `Minimum length ${validatorValue.requiredLength}, actual length: ${validatorValue.actualLength}`,
			'maxlength': `Maximum length ${validatorValue.requiredLength}, actual length: ${validatorValue.actualLength}`
		};

		return config[validatorName]
	}

	static timestampValidity(control) {
		if (control.value.match(/^(\d{4}\-\d{2}\-\d{2}T\d{2}:\d{2}:\d{2}Z|\s*)$/)) {
			return null;
		} else {
			return {
				'invalidTimestampStruct': true
			}
		}
	}

	static apiKeyLength(length: number) {
		return (control) => {
			switch (true) {
				case (control.value.length === 0):
				return null;

				case (control.value.length > length):
				return {
					'maxlength': {
						requiredLength: length,
						actualLength: control.value.length
					}
				}
				case (control.value.length < length):
				return {
					'minlength': {
						requiredLength: length,
						actualLength: control.value.length
					}
				}
				default:
				break;
			}
		}
	}
}
