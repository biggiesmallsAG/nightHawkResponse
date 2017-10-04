import json
import datetime
from django.views.generic.edit  import View
from django.shortcuts import render
from django.core.urlresolvers import reverse
from django.http import HttpResponse, HttpResponseRedirect, HttpRequest, JsonResponse
from django.template import loader, RequestContext
from django.utils.decorators import method_decorator
from django.contrib.auth.decorators import login_required
from django.views.decorators.csrf import ensure_csrf_cookie, csrf_protect
from django.contrib.auth.models import User
from nighthawk.forms import TasksForm
from nighthawk.models import TaskModel
from ws4redis.publisher import RedisPublisher
from ws4redis.redis_store import RedisMessage

class Tasks(View):
	def __init__(self):
		self.redis_publish = RedisPublisher(facility='tasks', broadcast=True)
		self.task_comment = ''

	def get(self, request):
		assert isinstance(request, HttpRequest)
		return render(request, 'tasks.html',
		              context_instance=RequestContext(request, {
		              		"TaskForm": TasksForm
              	}))

	def post(self, request):
		if request.is_ajax():
			objectmodel = json.loads(request.body)
			for x in objectmodel:
				if not x['task_inactive']:
					task_model = TaskModel(task_date=str(datetime.datetime.utcnow()), task_analyst=str(request.user), 
						task_assignee=x['task_assignee'], task_endpoints=x['endpoints'], task_name=x['task_name'],
						task_urgency=x['urgency'], task_inactive=x['task_inactive'])

					task_model.save()
					self.task_comment = {
										"user": x['task_assignee'],
										"task": x['task_name'],
										"assigned_by": str(request.user),
										"urgency": x['urgency']
									}

					self.redis_publish.publish_message(RedisMessage(json.dumps(self.task_comment)))
					return HttpResponse(task_model.id)

				task_model = TaskModel.objects.all()
				for y in task_model:
					if y.id == x['id']:
						y.task_inactive = True
						y.save()
						return HttpResponse(y.id)


	def GetAnalyst(self, request):
		if request.is_ajax():
			users = User.objects.all()
			user_obj = {
				"analyst": str(request.user),
				"assignee": [x.username for x in users]
			}

			return JsonResponse(user_obj, safe=False)

	def GetActiveTasks(self, request):
		if request.is_ajax():
			tasks = []
			for x in TaskModel.objects.all():
				if x.task_inactive:
					tasks.append({
							"datetime": x.task_date,
							"id": x.id,
							"task_name": x.task_name,
							"endpoints": x.task_endpoints,
							"task_assignee": x.task_assignee,
							"task_inactive": x.task_inactive, 
							"task_done": True,
							"task_analyst": x.task_analyst,
							"urgency": x.task_urgency
						})

				else:
					tasks.append({
							"datetime": x.task_date,
							"id": x.id,
							"task_name": x.task_name,
							"endpoints": x.task_endpoints,
							"task_assignee": x.task_assignee,
							"task_inactive": x.task_inactive, 
							"task_done": False,
							"task_analyst": x.task_analyst,
							"urgency": x.task_urgency
						})	

			return JsonResponse(tasks, safe=False)
