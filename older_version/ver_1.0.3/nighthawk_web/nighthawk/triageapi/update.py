from nighthawk.triageapi.dataendpoint.common import CommonAttributes
import update_control
import requests
from requests import ConnectionError
import elasticsearch
from elasticsearch_dsl import Search, Q, A
import json
from nighthawk.triageapi.utility.validate import ValidateUserInput
import datetime
from ws4redis.publisher import RedisPublisher
from ws4redis.redis_store import RedisMessage

class UpdateES(CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def UpdateDoc(self, row_data, user):
		err = {"error": "Input Validation Failed"}

		if not ValidateUserInput(row_data['comment']).ValidateInputMixed():
			return err

		query = {
			"doc": {
				"Record": {
					"Comment": {
						"Date": str(datetime.datetime.utcnow()),
						"Analyst": str(user),
						"Comment": row_data['comment']
					},
					"Tag": update_control.TagIntToStr(row_data['tag'])
				}
			}
		}

		try:
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/{0}/_update?parent={1}'.format(row_data['rowId'], row_data['parent']), data=json.dumps(query), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		try:
			q = requests.get(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + "/{0}?parent={1}".format(row_data['rowId'], row_data['parent']), auth=(self.elastic_user, self.elastic_pass), verify=False)
			case = q.json()['_source']['CaseInfo']['case_name']
		except ConnectionError as e:
			_ret = {"connection_error": e.args[0]}
			return _ret

		redis_publisher = RedisPublisher(facility='comments', broadcast=True)
		
		broadcast_comment = {
			"comment": row_data['comment'],
			"endpoint": row_data['parent'],
			"case": case,
			"analyst": str(user)
		}

		redis_publisher.publish_message(RedisMessage(json.dumps(broadcast_comment)))

		return r.json()

	def GetDocByComment(self, row_data):
		try:
			r = requests.get(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/{0}?parent={1}'.format(row_data['rowId'], row_data['parent']), auth=(self.elastic_user, self.elastic_pass), verify=False)
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
			r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			ret = {"connection_error": e.args[0]}
			return ret

		return r.json()['hits']['hits']
