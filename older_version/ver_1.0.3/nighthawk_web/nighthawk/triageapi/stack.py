from nighthawk.triageapi.dataendpoint.common import CommonAttributes
import search_queries
import elasticsearch
from elasticsearch_dsl import Search, Q, A
import requests
from requests import ConnectionError
import json
import time
import stack_queries

class StackES(CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def BuildRootTree(self):
		s = Search()
		t = Q('has_parent', type='hostname', query=Q('query_string', query="*"))
		aggs = A('terms', field='AuditType.Generator', size=16)

		s.aggs.bucket('datatypes', aggs)
		query = s.query(t)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = [{
			"id": "stackable", "parent": "#", "text": "Stackable Data", "type": "root"
		}]

		i = ['w32services', 'w32tasks', 'w32scripting-persistence', 'w32prefetch', 'w32network-dns', 'urlhistory', 'filedownloadhistory']

		for x in r.json()['aggregations']['datatypes']['buckets']:
			if x['key'] not in i:
				pass
			else:
				data.append({
					"id" : x['key'], "parent": "stackable", "text": x['key'], "children": True, "type": "stack"
				})

		return data

	def BuildAuditAggs(self, child_id):
		s = Search()
		s = s[0]
		t = Q('has_parent', type='hostname', query=Q('query_string', query='*'))
		aggs = A('terms', field='ComputerName.raw', size=0)


		s.aggs.bucket('endpoints', aggs)
		query = s.query(t).filter('term', AuditType__Generator=child_id)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = []

		for x in r.json()['aggregations']['endpoints']['buckets']:
			data.append({
				"id": x['key'].upper(), "parent": child_id, "text": x['key'].upper(), "type": "endpoint", "a_attr": {"href": "#" + x['key'].upper() + "/" + child_id}
			})

		return data

	def GetAuditData(self, stack_data):
		query = stack_queries.GetAuditGenerator(stack_data)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret
		
		for k, v in stack_data.iteritems():
			audittype = k
			endpoints = len(v)

		data = []
		exclude = ['', 'URL ?', '1URL ?', 'gURL ?', 'wURL ?', 'sURL ?', 'wwURL ?', 'bURL ?', '2URL ?', 'auURL ?', 'dURL ?', 'geoURL ?', 'maxURL ?', 'teURL ?', 'cURL ?', 'wwwURL ?']

		for x in r.json()['aggregations']['generator']['buckets']:
			if x['key'] not in exclude:			
				data.append({
						"attribute": x['key'], "endpoints": [y['key'].upper() for y in x['endpoint']['buckets']], "endpoint_count": len([y['key'].upper() for y in x['endpoint']['buckets']]), "doc_count": x['doc_count'], "audittype": audittype
					})

		return data

