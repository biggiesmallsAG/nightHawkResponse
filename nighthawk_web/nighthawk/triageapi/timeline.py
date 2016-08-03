from nighthawk.triageapi.dataendpoint.common import CommonAttributes
import search_queries
import elasticsearch
from elasticsearch_dsl import Search, Q, A
import requests
from requests import ConnectionError
import json
import timeline_queries

class TimeLineES(CommonAttributes):
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
			"id": "timeline", "parent": "#", "text": "Timeline", "type": "root"
		}]

		for x in r.json()['aggregations']['casenum']['buckets']:
			data.append({
				"id" : x['key'], "parent": "timeline", "text": x['key'], "children": True, "type": "case"
			})

		return data

	def BuildAuditAggs(self, child_id):
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
					"id" : x['_id'], "parent": child_id, "text": x['_id'].upper(), "type": "endpoint"
				})

		return data

	def GetAuditData(self, case, endpont_id, start=None, length=None, str_query=None, sort=None, order=None):
		query = timeline_queries.GetGeneratorQuery(case, endpont_id, start, length, str_query, sort, order)

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		data = []

		for x in r.json()['hits']['hits']:
			generator = x['fields']['AuditType.Generator'][0]
			
			if generator == 'w32scripting-persistence':
				data.append({
						"time": x['fields']['Record.TlnTime'], 
						"path": x['fields']['Record.Path'], 
						"generator": x['fields']['AuditType.Generator'],
						"file_accessed": x['fields']['Record.File.Accessed'],
						"file_modified": x['fields']['Record.File.Modified'],
						"file_changed": x['fields']['Record.File.Changed']
					})

			elif generator == 'w32rawfiles':
				data.append({
						"time": x['fields']['Record.TlnTime'], 
						"path": x['fields']['Record.Path'], 
						"generator": x['fields']['AuditType.Generator'],
						"file_accessed": x['fields']['Record.FilenameAccessed'],
						"file_modified": x['fields']['Record.FilenameModified'],
						"file_changed": x['fields']['Record.FilenameChanged']
					})

			elif generator == 'urlhistory':
				data.append({
						"time": x['fields']['Record.TlnTime'], 
						"path": x['fields']['Record.Url'], 
						"generator": x['fields']['AuditType.Generator'],
						"file_accessed": "",
						"file_modified": "",
						"file_changed": ""
					})

			elif generator == 'filedownloadhistory':
				data.append({
						"time": x['fields']['Record.TlnTime'], 
						"path": x['fields']['Record.SourceUrl'], 
						"generator": x['fields']['AuditType.Generator'],
						"file_accessed": "",
						"file_modified": "",
						"file_changed": ""
					})
			
			elif generator == 'w32registryraw':
				data.append({
						"time": x['fields']['Record.TlnTime'], 
						"path": x['fields']['Record.Path'], 
						"generator": x['fields']['AuditType.Generator'],
						"file_accessed": "",
						"file_modified": "",
						"file_changed": ""
					})

		ret = {
			"data": data
		}

		return ret

