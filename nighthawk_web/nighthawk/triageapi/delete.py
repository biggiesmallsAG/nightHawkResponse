from nighthawk.triageapi.dataendpoint.common import CommonAttributes
import delete_queries
import elasticsearch
from elasticsearch_dsl import Search, Q, A
import requests
from requests import ConnectionError
import json

class DeleteCaseObj(CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def BuildCaseList(self):
		s = Search()
		t = Q('query_string', query="*")
		aggs_casenum = A('terms', field="CaseInfo.case_name", size=0)

		s.aggs.bucket('casenum', aggs_casenum)
		query = s.query(t)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = [{
			"case": ""
		}]

		for x in r.json()['aggregations']['casenum']['buckets']:
			data.append({
				"case": x['key']
			})

		return data

	def BuildEndpointList(self):
		s = Search()
		aggs_generator = A('terms', field='ComputerName.raw', size=0)

		s.aggs.bucket('endpoints', aggs_generator)
		query = s.query()

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = [{
			"endpoint": ""
		}]

		for x in r.json()['aggregations']['endpoints']['buckets']:
			data.append({
					"endpoint" : x['key']
				})

		return data

	def Delete(self, del_obj):
		s = Search()
		s = s[0:10]

		query = delete_queries.GetGeneratorQuery(del_obj)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query), auth=(self.elastic_user, self.elastic_pass), verify=False)
			parent = r.json()['hits']['hits'][0]['_parent']
			r = requests.delete(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_query?parent={0}'.format(parent), data=json.dumps(query), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		return r.json()
