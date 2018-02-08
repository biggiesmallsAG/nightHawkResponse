export interface SearchHandler {
	search_term: string,
	search_size: number,
	case_name: string,
	search_limit: number,
	ignore_good: boolean,
	path: string
}

export interface TimelineHandler extends SearchHandler {
	endpoint: string,
	endpoint_list: Array<string>,
	start_time: string,
	end_time: string,
	start_date: string,
	end_date: string,
	time_delta: string
}