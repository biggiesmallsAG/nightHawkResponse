import json
from django.views.generic.edit  import View
from django.shortcuts import render
from django.core.urlresolvers import reverse
from django.http import HttpResponse, HttpResponseRedirect, HttpRequest, JsonResponse
from django.template import loader, RequestContext
from django.utils.decorators import method_decorator
from django.contrib.auth.decorators import login_required
from django.views.decorators.csrf import ensure_csrf_cookie, csrf_protect

from nighthawk.triageapi.update import UpdateES
from nighthawk.forms import UpdateDoc as UpdateDocForm

class UpdateDoc(View):
	def __init__(self):
		pass

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'update_doc.html',
	                  context_instance=RequestContext(request, {"UpdateDoc": UpdateDocForm}))

	@method_decorator(csrf_protect)
	def post(self, request):
		row_data = json.loads(request.body)
		u = UpdateES()

		data = u.UpdateDoc(row_data, request.user)
		
		return JsonResponse(data, safe=False)
