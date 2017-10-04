export interface AuditHandler {
	case_name?: string,
	case_date?: string,
	endpoint?: string,
	audit_type?: string,
	total_hits?: any,
	data?: any
}
