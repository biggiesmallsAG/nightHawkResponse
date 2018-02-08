export interface NhResponseFailure {
	reason: string,
	response: string
}

export interface NhResponseSuccess extends NhResponseFailure {
	data: Array<any>
}