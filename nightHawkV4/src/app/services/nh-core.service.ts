import { Injectable } from '@angular/core';
import { Headers, Http, RequestOptions } from '@angular/http';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/toPromise';

@Injectable()
export class NhCoreService {
	private apiUrl = '/api/v1';
	private token: string = localStorage.getItem('nhr-token');
	private options = new RequestOptions({
		headers: new Headers({
			'nhr-token': this.token,
			'Content-Type': 'application/json'
		})
	});
	constructor(private http: Http) {}

	GET(endpoint: string) {
		return this.http.get(this.apiUrl + endpoint, this.options)
			.map(response => response.json().data)
			.catch(error => this.handleError(error.reason))
	}

	POST(endpoint: string, name: any) {
		return this.http.post(this.apiUrl + endpoint, JSON.stringify({ name: name }), this.options)
			.map(response => {
				switch (endpoint) {
					case "/config":
						return response.json()
					default:
						return response.json().data
				}
			})
			.catch(error => this.handleError(error.reason))
	}

	POSTJSON(endpoint: string, body: any) {
		return this.http.post(this.apiUrl + endpoint, JSON.stringify(body), this.options)
			.map(response => response.json().data)
			.catch(error => this.handleError(error))
	}

	POSTUpload(endpoint: string, data: any) {
		return this.http.post(this.apiUrl + endpoint, data, this.options)
			.map(response => response.json())
			.catch(error => this.handleError(error.reason))
	}

	private handleError(error: any): Promise<any> {
		console.error('An error occurred', error);
		return Promise.reject(error.message || error.reason);
	}
}