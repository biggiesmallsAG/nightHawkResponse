from django.views.generic.edit  import View
from django.shortcuts import render
from django.core.urlresolvers import reverse
from django.http import HttpResponse, HttpResponseRedirect, HttpRequest, JsonResponse
from django.template import loader, RequestContext
from django.contrib.auth.decorators import login_required
from django.utils.decorators import method_decorator
from django.views.decorators.csrf import csrf_protect
from nighthawk.triageapi.search import QueryES
from nighthawk.triageapi.dataendpoint.common import CommonAttributes
from nighthawk.forms import SearchForm
import json

class Home(View, CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)
		self.welcome = "Welcome to the nightHawk Response framework. Use the node tree on the left to select an endpoint."

	def Home404(self, request):
		assert isinstance(request, HttpRequest)
		return render(request, '404.html')

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'base.html',
	                  context_instance=RequestContext(request, {"name": self.name, "version": self.nighthawk_version, "user": request.user}))

	@method_decorator(csrf_protect)
	def post(self, request):
		if request.method == 'POST' and request.is_ajax():
			return HttpResponse(self.welcome)

	def LoadCaseTree(self, request):
		if request.method == 'GET' and request.is_ajax():
			child_id = request.GET.get('id', '')
			q = QueryES()

			if child_id == "#":				
				data = q.BuildRootTree()

				return JsonResponse(data, safe=False)

			else:
				data = q.BuildEndpointAggs(child_id)
				
				return JsonResponse(data, safe=False)

	def LoadCaseTreeAudit(self, request):
		if request.method == 'GET' and request.is_ajax():
			child_id = request.GET.get('id', '')
			parent_id = request.GET.get('parent', '')
			q = QueryES()

			data = q.BuildAuditAggs(child_id, parent_id)
			
			return JsonResponse(data, safe=False)

class HomeSearch(View, CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'main_search.html',
	                  context_instance=RequestContext(request, {
	                  	"name": self.name, 
	                  	"global": self.global_search, 
	                  	"version": self.global_search_version,
	                  	"SearchForm": SearchForm
	                  	}))

	@method_decorator(csrf_protect)
	def post(self, request):
		if request.method == 'POST' and request.is_ajax():
			body = json.loads(request.body)
			data = body['data'].split('=')[1]
			q = QueryES()

			data = q.GetAuditDataMain(data)

			return JsonResponse(data, safe=False)