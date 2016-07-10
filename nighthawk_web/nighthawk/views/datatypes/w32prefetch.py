from django.views.generic.edit  import View
from django.shortcuts import render
from django.core.urlresolvers import reverse
from django.http import HttpResponse, HttpResponseRedirect, HttpRequest, JsonResponse
from django.template import loader, RequestContext
from django.contrib.auth.decorators import login_required
from django.utils.decorators import method_decorator
from django.views.decorators.csrf import csrf_protect, csrf_exempt, ensure_csrf_cookie
from nighthawk.triageapi.search import QueryES
import json

class W32Prefetch(View):
	def __init__(self):
		pass

	def get(self, request, case=None, hostname=None):
		assert isinstance(request, HttpRequest)
		return render(request, 'w32prefetch.html',
	                  context_instance=RequestContext(request, {"hostname": hostname, "case": case}))

	@method_decorator(csrf_protect)
	def post(self, request, case=None, hostname=None):
		if request.method == 'POST' and request.is_ajax():
			q = QueryES()
			data = q.GetAuditData(case, hostname, "w32prefetch")

			return JsonResponse(data, safe=False)
