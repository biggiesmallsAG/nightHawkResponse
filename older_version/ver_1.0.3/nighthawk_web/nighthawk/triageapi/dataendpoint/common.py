import json

class CommonAttributes():
	def __init__(self):
		with open('/opt/nighthawk/etc/nighthawk.json', 'r') as config:
			self.conf_data = json.load(config)

		with open('/opt/nighthawk/lib/elastic/ElasticMapping.json', 'r') as mapping:
			self.mapping_file = json.load(mapping)

		self.name = 'nightHawk'
		self.nighthawk_version = 'v1.0.3'
		self.nighthawk_stack = 'Stack'
		self.nighthawk_stack_ver = 'v0.7'
		self.nighthawk_timeline = 'Timeline'
		self.nighthawk_timeline_ver = 'v0.3'
		self.global_search = 'Global Search'
		self.global_search_version = "v0.2"

		if self.conf_data['elastic']['elastic_ssl']:
			self.es_host = "https://{0}".format(self.conf_data['elastic']['elastic_server'])
			self.es_port = str(self.conf_data['elastic']['elastic_port'])
		else:
			self.es_host = "http://{0}".format(self.conf_data['elastic']['elastic_server'])
			self.es_port = str(self.conf_data['elastic']['elastic_port'])

		self.elastic_user = self.conf_data['elastic']['elastic_user']
		self.elastic_pass = self.conf_data['elastic']['elastic_pass']
		self.index = '/investigations'
		self.type_hostname = '/hostname'
		self.type_audit_type = '/audit_type'
		self.file_upload_max = self.conf_data['nightHawk']['max_file_upload']