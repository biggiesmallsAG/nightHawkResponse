import {Injectable}      from '@angular/core';
import { Headers, Http } from '@angular/http';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/toPromise';

@Injectable()
export class NhCoreService {
	public apiUrl = '/api/v1'; 
	private headers = new Headers({'Content-Type': 'application/json'});
	private mutlipart = new Headers({'Content-Type': 'multipart/form-data'})

	constructor(private http: Http) { }

	GET(endpoint: string) {
		return this.http.get(this.apiUrl + endpoint)
		.map(response => response.json().data)
		.catch(error => this.handleError(error.reason))
	}

	POST(endpoint: string, name: any) {
		return this.http.post(this.apiUrl + endpoint, JSON.stringify({name: name}), {headers: this.headers})
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
		return this.http.post(this.apiUrl + endpoint, JSON.stringify(body), {headers: this.headers})
			.map(response => response.json().data)
			.catch(error => this.handleError(error))
	}

	POSTUpload(endpoint: string, data: any) {
		return this.http.post(this.apiUrl + endpoint, data)
			.map(response => response.json())
			.catch(error => this.handleError(error.reason))
	}

	private handleError(error: any): Promise<any> {
		console.error('An error occurred', error);
		return Promise.reject(error.message || error.reason);
	}
}