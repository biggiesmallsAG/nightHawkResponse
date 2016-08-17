import elasticsearch
from elasticsearch_dsl import Search, Q, A

def GetGeneratorQuery(del_obj):
	s = Search()
	
	if not del_obj['singular']:	
		if del_obj['case']:
			t = Q('query_string', default_field="CaseInfo.case_name", query=del_obj['case'])

		if del_obj['endpoint']:
			t = Q('query_string', default_field="ComputerName.raw", query=del_obj['endpoint'])

		query = s.query(t)
		return query.to_dict()

	t = Q('query_string', default_field="CaseInfo.case_name", query=del_obj['case']) & Q('query_string', default_field="ComputerName.raw", query=del_obj['endpoint'])
	query = s.query(t)
	return query.to_dict()