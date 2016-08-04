from nighthawk.triageapi.dataendpoint.common import CommonAttributes
import search_queries
import elasticsearch
from elasticsearch_dsl import Search, Q, A
import requests
from requests import ConnectionError
import json

class QueryES(CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def BuildRootTree(self):
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
			"id": "current_inv", "parent": "#", "text": "Current Investigations", "type": "root"
		}, {
			"id": "comp_inv", "parent": "#", "text": "Completed Investigations", "type": "root"
		}]

		for x in r.json()['aggregations']['casenum']['buckets']:
			data.append({
				"id": x['key'], "parent": "current_inv", "text": x['key'], "children": True, "type": "case"
			})

		return data

	def BuildEndpointAggs(self, child_id):
		s = Search()
		s = s[0:1000]
		t = Q('has_child', type='audit_type', query=Q('query_string', default_field="CaseInfo.case_name", query=child_id))
		query = s.query(t)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_hostname + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = []

		for x in r.json()['hits']['hits']:
			data.append({
					"id" : x['_id'], "parent": child_id, "text": x['_id'].upper(), "children": True, "type": "endpoint"
				})

		return data

	def BuildAuditAggs(self, child_id, parent_id):
		s = Search()
		s = s[0]
		t = Q('query_string', default_field="CaseInfo.case_name", query=parent_id) & Q('match', ComputerName=child_id)
		aggs_generator = A('terms', field='AuditType.Generator', size=0)

		s.aggs.bucket('datatypes', aggs_generator)
		query = s.query(t)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = []
		exclude = ['w32processes-memory', 'stateagentinspector', 'w32disks']

		for y in r.json()['aggregations']['datatypes']['buckets']:
			if not y['key'] in exclude:
				data.append({
						"id": y['key'], "parent": child_id, "text": y['key'], "type": "audit", "a_attr": {"href": "#" + y['key'] + '/' + parent_id + "/" + child_id }
					})

		return data

	def GetAuditData(self, case, child_id, data_type, start=None, length=None, str_query=None, sort=None, order=None):
		q = ['w32registryraw', 'filedownloadhistory', 'urlhistory', 'timeline', 'w32apifiles', 'w32rawfiles', 'w32eventlogs']

		if data_type in q:
			query = search_queries.GetGeneratorQuery(data_type, str_query, case, child_id, start, length, sort, order)
		else:
			s = Search()
			s = s[0:1000]
			t = Q('query_string', default_field="ComputerName.raw", query=child_id) & Q('query_string', default_field="CaseInfo.case_name", query=case)
			query = s.query(t).filter('term', AuditType__Generator=data_type)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = []

		try:
			for x in r.json()['hits']['hits']:
				data.append(x)
		except KeyError:
			return data

		return data

	def GetAuditDataMain(self, data):
		s = Search()
		s = s[0:1000]
		s = s.highlight('*')
		s = s.highlight_options(require_field_match=False)
		t = Q('query_string', query=data) & ~Q('query_string', default_field="AuditType.Generator", query="stateagentinspector") & ~Q('query_string', default_field="AuditType.Generator", query="w32processes-tree")

		query = s.query(t)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = []

		try:
			for x in r.json()['hits']['hits']:
				for y, v in x['highlight'].iteritems():
					data.append({
							"doc_id": x['_id'],
							"endpoint": x['_parent'],
							"audittype": x['_source']['AuditType']['Generator'],
							"field": y,
							"response": v
						})
		except KeyError:
			pass

		return data