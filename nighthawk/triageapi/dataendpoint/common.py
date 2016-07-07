import json

class CommonAttributes():
	def __init__(self):
		with open('/opt/nighthawk/etc/nighthawk.json', 'r') as config:
			self.conf_data = json.load(config)

		self.name = 'nightHawk'
		self.nighthawk_version = 'v1.0'
		self.nighthawk_stack = 'Stack'
		self.nighthawk_stack_ver = 'v0.3'
		self.nighthawk_timeline = 'Timeline'
		self.nighthawk_timeline_ver = 'v0.2'
		self.global_search = 'Global Search'
		self.global_search_version = "v0.1"
		self.es_host = "http://{0}".format(self.conf_data['elastic']['elastic_server'])
		self.es_port = ':9200'
		self.index = '/investigations'
		self.type_hostname = '/hostname'
		self.type_audit_type = '/audit_type'
		self.file_upload_max = self.conf_data['nightHawk']['max_file_upload']