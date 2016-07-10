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

class W32RAWFiles(View):
	def __init__(self):
		pass

	def get(self, request, case=None, hostname=None):
		assert isinstance(request, HttpRequest)
		return render(request, 'w32rawfiles.html',
	                  context_instance=RequestContext(request, {"hostname": hostname, "case": case}))

	@method_decorator(csrf_protect)
	def post(self, request, case=None, hostname=None):
		if request.method == 'POST' and request.is_ajax():
			start = request.POST.get("start", '')
			length = request.POST.get("length", '')
			search_val = request.POST.get("search[value]")
			sort = request.POST.get("order[0][column]", '')
			order = request.POST.get("order[0][dir]", '')

			q = QueryES()
			data = q.GetAuditData(case, hostname, "w32rawfiles", start=start, length=length, str_query=search_val, sort=sort, order=order)

			return JsonResponse(data, safe=False)
