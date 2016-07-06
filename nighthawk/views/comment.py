# Analyst Home view
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

class Comment(View):
	def __init__(self):
		pass

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'comment.html',
	                  context_instance=RequestContext(request))

	@method_decorator(csrf_protect)
	def post(self, request):
		u = UpdateES()
		data = u.GetSessionComments()
		
		return JsonResponse(data, safe=False)

	def CommentDoc(self, request):
		if request.method == 'POST' and request.is_ajax():
			doc_data = json.loads(request.body)
			u = UpdateES()
			data = u.GetDocByComment(doc_data)

			return JsonResponse(data, safe=False)

	def DocDiaglog(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'doc_dialog.html',
	                  context_instance=RequestContext(request))