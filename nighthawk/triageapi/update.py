from nighthawk.triageapi.dataendpoint.common import CommonAttributes
import update_control
import requests
from requests import ConnectionError
import elasticsearch
from elasticsearch_dsl import Search, Q, A
import json

class UpdateES(CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def UpdateDoc(self, row_data):
		query = {
			"doc": {
				"Record": {
					"Comment": {
						"Date": row_data['date'],
						"Analyst": row_data['analyst'],
						"Comment": row_data['comment']
					},
					"Tag": update_control.TagIntToStr(row_data['tag'])
				}
			}
		}

		try:
			r = requests.post(self.es_host + self.es_port + self.index + self.type_audit_type + '/{0}/_update?parent={1}'.format(row_data['rowId'], row_data['parent']), data=json.dumps(query))
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		return r.json()

	def GetDocByComment(self, row_data):
		try:
			r = requests.get(self.es_host + self.es_port + self.index + self.type_audit_type + '/{0}?parent={1}'.format(row_data['rowId'], row_data['parent']))
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		return r.json()

	def GetSessionComments(self):
		s = Search()
		s = s[0:1000]
		t = ~Q('query_string', default_field="AuditType.Generator", query='w32tasks') & Q('query_string', default_field='Record.Comment.Date', query="2*")
		query = s.query(t)

		try:
			r = requests.post(self.es_host + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()))
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		return r.json()['hits']['hits']
