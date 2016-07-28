## Platform stats helper by Daniel Eden
## 26/07/2016 Update 
## daniel.eden@gmail.com

import psutil
import time
import json
import sys
import os 

sys.path.append('/Users/biggiesmalls/Documents/brstriage/nightHawkResponse/nighthawk_web')
os.environ.update(DJANGO_SETTINGS_MODULE='nighthawk.settings')
from ws4redis.publisher import RedisPublisher
from ws4redis.redis_store import RedisMessage

sleep_cycle = 2

class StatsConsumer(object):
	def __init__(self):
		self.cpu_stats = ''
		self.mem_stats = ''
		self.disk_stats = ''
		self.processes = ''

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
		except psutil.ZombieProcess:
			pass

		self.processes = core_services
			
	def GetFloatGB(self, number):
		return float((((number/1024)/1024)/1024))

class StatsBroadcaster(StatsConsumer):
	def __init__(self):
		StatsConsumer.__init__(self)
		self.redis_publisher = RedisPublisher(facility='platform', broadcast=True)

	def BroadCast(self, cpu_stats, mem_stats, disk_stats, processes):
		os_stats = {
			"cpu": cpu_stats,
			"mem": mem_stats,
			"disk": disk_stats,
			"processes": processes
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

		b.BroadCast(s.cpu_stats, s.mem_stats, s.disk_stats, s.processes)
		time.sleep(sleep_cycle)

if __name__ == '__main__':
	main()