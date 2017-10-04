import elasticsearch
from elasticsearch_dsl import Search, Q, A

def GetGeneratorQuery(case, endpoint_id, start, length, str_query, sort, order):
	s = Search()
	s = s[int(start):int(length)+int(start)]
	s = s.fields([	"Record.Path",
				    "Record.Url",
				    "Record.SourceUrl",
				    "Record.TlnTime",
				    "Record.File.Accessed",
				    "Record.File.Modified",
				    "Record.File.Changed",
				    "Record.FilenameChanged",
				    "Record.FilenameModified",
				    "Record.FilenameAccessed",
				    "AuditType.Generator"
				])

	order_dict = {
		"0": "Record.TlnTime",
		"1": "AuditType.Generator",
		"2": "Record.File.Modified",
		"3": "Record.File.Accessed",
		"4": "Record.File.Changed",
		"5": "Record.Path"
	}

	if str_query == "":
		_sort = {
			"{0}".format(order_dict[str(sort)]): {
				"order": order
			}
		}

		t = Q('query_string', default_field="Record.TlnTime", query="*") & Q('match', ComputerName=endpoint_id) & ~Q('match', AuditType__Generator="w32processes-memory") & ~Q('match', AuditType__Generator="w32useraccounts")
		query = s.query(t).filter('term', CaseInfo__case_name=case).sort(_sort)

	else:
		_sort = {
			"{0}".format(order_dict[str(sort)]): {
				"order": order
			}
		}
		
		t = Q('query_string', default_field="Record.TlnTime", query="*") & Q('match', ComputerName=endpoint_id) & ~Q('match', AuditType__Generator="w32processes-memory") & ~Q('match', AuditType__Generator="w32useraccounts") & Q('query_string', fields=[
					"Record.Path",
				    "Record.Url",
				    "Record.SourceUrl",
				    "AuditType.Generator"], query="{0}*".format(str_query))
		query = s.query(t).filter('term', CaseInfo__case_name=case).sort(_sort)

	return query.to_dict()