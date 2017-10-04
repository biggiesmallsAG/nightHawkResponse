#!/usr/bin/env python
## Platform stats helper by Daniel Eden
## 28/08/2016 Update 
## Version: 0.2
## daniel.eden@gmail.com

import psutil
import time
import json
import sys
import os 
import requests
from requests import ConnectionError
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)
import elasticsearch
from elasticsearch_dsl import Search, A

sys.path.append('/opt/nighthawk/web')
os.environ.update(DJANGO_SETTINGS_MODULE='nighthawk.settings')
from nighthawk.triageapi.dataendpoint.common import CommonAttributes
from ws4redis.publisher import RedisPublisher
from ws4redis.redis_store import RedisMessage

sleep_cycle = 10

class StatsConsumer(object, CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)
		self.cpu_stats = ''
		self.mem_stats = ''
		self.disk_stats = ''
		self.processes = ''
		self.es_stats = ''

	def GenerateCPU(self):
		self.cpu_stats = psutil.cpu_times_percent(percpu=True)

	def GenerateMem(self):
		mem_stats = psutil.virtual_memory()
		self.mem_stats = {
			"available": "{0}GB".format(self.GetFloatGB(mem_stats.available)),
			"free": "{0}GB".format(self.GetFloatGB(mem_stats.free)),
			"total": "{0}GB".format(self.GetFloatGB(mem_stats.total)),
			"used": "{0}GB".format(self.GetFloatGB(mem_stats.used))
		}

	def GenerateDisk(self):
		disk_stats = psutil.disk_partitions()
		disk_array = []

		for x in disk_stats:
			y = psutil.disk_usage(x.mountpoint)
			disk_array.append({
				"device": x.device,
				"mountpoint": x.mountpoint,
				"usage": {
					"total": "{0}GB".format(self.GetFloatGB(y.total)),
					"free": "{0}GB".format(self.GetFloatGB(y.free))
				}
			})

		self.disk_stats = disk_array

	def GenerateProcess(self):
		core_pids = ['nginx', 'uwsgi', 'java', 'nightHawk']
		core_services = []
		
		try:
			for x in psutil.pids():
				p = psutil.Process(x)
				if p.name() in core_pids:
					core_services.append({
							"service": {
								"name": p.name(),
								"cmdline": p.cmdline(),
								"status": p.status()
							} 
						})
		except:
			pass

		self.processes = core_services
			
	def GetEsStats(self):
		req = True
		while req:
			try:
				try:
					r = requests.get(self.es_host + ":" + self.es_port + '/_cluster/stats', verify=False, auth=(self.elastic_user, self.elastic_pass))
				except ConnectionError as e:
					return {"connection_error": e.args[0]}

				data = r.json()
				s = Search()

				aggs_cases = A('terms', field='CaseInfo.case_name', size=0)
				s.aggs.bucket('cases', aggs_cases)
				query = s.query()

				try:
					r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
				except ConnectionError as e:
					return {"connection_error": e.args[0]}

				aggs = r.json()
				cases = aggs['aggregations']['cases']['buckets']
				
				s = Search()

				aggs_endpoints = A('terms', field="ComputerName.raw", size=0)
				s.aggs.bucket('endpoints', aggs_endpoints)
				query = s.query()

				try:
					r = requests.post(self.es_host + ":" + self.es_port + self.index + self.type_audit_type + '/_search', data=json.dumps(query.to_dict()), auth=(self.elastic_user, self.elastic_pass), verify=False)
				except ConnectionError as e:
					return {"connection_error": e.args[0]}

				aggs = r.json()
				endpoints = aggs['aggregations']['endpoints']['buckets']

				self.es_stats = {
					"docs": data['indices']['docs']['count'],
					"status": data['status'],
					"cluster_name": data['cluster_name'],
					"indices": data['indices']['count'],
					"nodes": data['nodes'],
					"cases": len(cases),
					"endpoints": len(endpoints)
				}

				req = False
			except:
				pass
				
	def GetFloatGB(self, number):
		return float((((number/1024)/1024)/1024))

class StatsBroadcaster(StatsConsumer):
	def __init__(self):
		StatsConsumer.__init__(self)
		self.redis_publisher = RedisPublisher(facility='platform', broadcast=True)

	def BroadCast(self, cpu_stats, mem_stats, disk_stats, processes, es_stats):
		os_stats = {
			"cpu": cpu_stats,
			"mem": mem_stats,
			"disk": disk_stats,
			"processes": processes,
			"es_stats": es_stats
		}

		self.redis_publisher.publish_message(RedisMessage(json.dumps(os_stats)))

def main():
	
	s = StatsConsumer()
	b = StatsBroadcaster()

	while True:
		s.GenerateCPU()
		s.GenerateMem()
		s.GenerateDisk()
		s.GenerateProcess()
		s.GetEsStats()

		b.BroadCast(s.cpu_stats, s.mem_stats, s.disk_stats, s.processes, s.es_stats)
		time.sleep(sleep_cycle)

if __name__ == '__main__':
	main()