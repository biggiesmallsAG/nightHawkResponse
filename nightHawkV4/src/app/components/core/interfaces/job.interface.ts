export interface JobHandler {
	timestamp: string,
	worker: string,
	log_level: string,
	body: JobBody
}

export interface JobBody {
	audit_file: Array<string>,
	case_id: string,
	in_progress: boolean,
	is_complete: boolean,
	uid: string,
	user_id: string
}
