#!/usr/bin/env python
## ElasticSearch Reindexer by Daniel Eden
## 29/08/2016 Update
## - Fixed SSL based communications 
## daniel.eden@gmail.com

import requests
from requests import ConnectionError
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)
import json
import re
import sys
sys.path.append('/opt/nighthawk/web')
from nighthawk.triageapi.dataendpoint.common import CommonAttributes

class SearchQuery(CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def CheckAliases(self):
		print "[+] Obtaining latest index alias to determine index number"
		try:
			r = requests.get(self.es_host + ":" + self.es_port + '/_aliases', auth=(self.elastic_user, self.elastic_pass), verify=False)
		except ConnectionError as e:
			print '[!] Error connecting to {0}{1}'.format(self.es_host, self.es_port)

		aliases = r.json()
		current_index = ''
		for x in aliases:
			try:
				for k, v in aliases[x]['aliases'].iteritems():
					if k == 'investigations':
						print '[+] Got latest investigations alias: {0}'.format(x)
						current_index = x
						return current_index
			except KeyError:
				pass

	def GetMappingAndCreateIndex(self, current_index):
		rex = re.compile("([a-zA-Z]+)([0-9]+)")
		m = rex.match(current_index)

		index_num = int(m.group(2)) + 1

		print "[+] New index will be {0}{1}".format(self.index, index_num)

		try:
			print '[-] Sending mapping to new index'
			r = requests.put("{0}:{1}{2}{3}".format(self.es_host, self.es_port, self.index, index_num), data=json.dumps(self.mapping_file), auth=(self.elastic_user, self.elastic_pass), verify=False)
			try:
				if r.json()['acknowledged']:
					print '[+] Returned successfully, index created.'
					return 0, int(m.group(2))
			except KeyError:
				print '[!] Error: {0}'.format(r.json()['error']['root_cause'][0]['reason'])
				return 1

		except ConnectionError as e:
			print '[!] Error connecting to {0}{1}'.format(self.es_host, self.es_port)

	def RemoveOldAlias(self, op_code, index_num):
		if op_code == 0:
			print '[-] Removing old alias'
			try:
				
				remove_alias = {
					"actions": [
						{
							"remove": {
								"index": "investigations{0}".format(index_num), "alias": "investigations"
							}
						}
					]
				}

				r = requests.post(self.es_host + ":" + self.es_port + '/_aliases', data=json.dumps(remove_alias), auth=(self.elastic_user, self.elastic_pass), verify=False)
				try:
					if r.json()['acknowledged']:
						print '[+] Returned successfully, alias removed.'
						return 0
				except KeyError:
					print '[!] Error: {0}'.format(r.json()['error']['root_cause'][0]['reason'])
					return 1
			except ConnectionError as e:
				print '[!] Error connecting to {0}{1}'.format(self.es_host, self.es_port)
		else:
			print '[!] Returned op_code 1, error in index creation and mapping. Exiting now'
			sys.exit(1)

	def ReindexData(self, op_code, index, index_num):
		if op_code == 0:
			print '[-] Reindexing data from {0} to investigations{1}'.format(index, (index_num+1))
			try:

				reindex = {
					"source": {
						"index": "investigations{0}".format(index_num)
					},
					"dest": {
						"index": "investigations{0}".format(index_num + 1)
					} 
				}

				print '[-] Large datasets will take a while, sit back and grab a coke....'
				r = requests.post(self.es_host + ":" + self.es_port + '/_reindex', data=json.dumps(reindex), auth=(self.elastic_user, self.elastic_pass), verify=False)

				try:
					if r.json()['created']:
						print '[+] Returned successfully, indexed {0} documents.'.format(r.json()['created'])
						return 0
				except KeyError:
					print '[!] Error: {0}'.format(r.json()['error']['root_cause'][0]['reason'])
					return 1
			except ConnectionError as e:
				print '[!] Error connecting to {0}{1}'.format(self.es_host, self.es_port)
		else:
			print '[!] Returned op_code 1, error in index creation and mapping. Exiting now'
			sys.exit(1)

	def Version(self):
		print "-- Reindexing automation by Daniel Eden (nightHawk Response team)."
		print "-- Version 1.0.3. 29/08/2016\n"

def main():

	s = SearchQuery()
	s.Version()
	index = s.CheckAliases()
	op_code, index_num = s.GetMappingAndCreateIndex(index)
	op_code = s.RemoveOldAlias(op_code, index_num)

	s.ReindexData(op_code, index, index_num)

if __name__ == '__main__':
	main()