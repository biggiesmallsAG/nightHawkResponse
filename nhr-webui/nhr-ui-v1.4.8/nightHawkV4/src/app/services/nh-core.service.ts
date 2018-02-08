import { Injectable } from '@angular/core';
import { Headers, Http, RequestOptions } from '@angular/http';
import * as API from './constants/api.constants';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/toPromise';

@Injectable()
export class NhCoreService {
	constructor(private http: Http) { }

	private getToken() {
		const token: string = localStorage.getItem('nhr-token');
		var options = new RequestOptions({
			headers: new Headers({
				'nhr-token': token,
				'Content-Type': 'application/json'
			})
		});
		return options
	}
	GET(endpoint: string) {
		const options = this.getToken();
		return this.http.get(API.apiUrl + endpoint, options)
			.map(response => {
				switch (endpoint) {
					case "/tag/show":
						return response.json()
					default:
						return response.json().data
				}
			})
			.catch(error => this.handleError(error.reason))
	}

	POST(endpoint: string, name: any) {
		const options = this.getToken();
		return this.http.post(API.apiUrl + endpoint, JSON.stringify({ name: name }), options)
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
		const options = this.getToken();
		return this.http.post(API.apiUrl + endpoint, JSON.stringify(body), options)
			.map(response => {
				switch (endpoint) {
					case "/auth/login":
						return response.json();
					default:
						return response.json().data
				}
			})
			.catch(error => this.handleError(error))
	}

	POSTUpload(endpoint: string, data: any) {
		const options = this.getToken();
		switch (endpoint) {
			case "/upload":
				let options = new RequestOptions({
					headers: new Headers({
						'nhr-token': localStorage.getItem('nhr-token')
					})
				});
				return this.http.post(API.apiUrl + endpoint, data, options)
					.map(response => response.json())
					.catch(error => this.handleError(error.reason))
		}
		return this.http.post(API.apiUrl + endpoint, data, options)
			.map(response => response.json())
			.catch(error => this.handleError(error.reason))
	}

	private handleError(error: any): Promise<any> {
		console.error('An error occurred', error);
		return Promise.reject(error.message || error.reason);
	}
}