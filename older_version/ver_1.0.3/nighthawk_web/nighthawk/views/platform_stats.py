import json
from django.views.generic.edit  import View
from django.shortcuts import render
from django.core.urlresolvers import reverse
from django.http import HttpResponse, HttpResponseRedirect, HttpRequest, JsonResponse
from django.template import loader, RequestContext
from django.utils.decorators import method_decorator
from django.contrib.auth.decorators import login_required
from django.views.decorators.csrf import ensure_csrf_cookie, csrf_protect

class PlatformStats(View):
	def __init__(self):
		pass

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'platform_stats.html',
	                  context_instance=RequestContext(request))
