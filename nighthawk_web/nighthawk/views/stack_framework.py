from django.views.generic.edit import View
from django.shortcuts import render
from django.core.urlresolvers import reverse
from django.http import HttpResponse, HttpResponseRedirect, HttpRequest, JsonResponse
from django.template import loader, RequestContext
from django.utils.decorators import method_decorator
from django.views.decorators.csrf import csrf_protect, csrf_exempt, ensure_csrf_cookie
from nighthawk.triageapi.dataendpoint.common import CommonAttributes
from nighthawk.triageapi.stack import StackES

import json

class StackView(View, CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'stack.html', 
	    	context_instance=RequestContext(request, {"name": self.name, "stack": self.nighthawk_stack, "stack_ver": self.nighthawk_stack_ver}))

	@method_decorator(csrf_protect)
	def post(self, request):
		stack_data = json.loads(request.body)
		s = StackES()

		data = s.GetAuditData(stack_data)
		
		return JsonResponse(data, safe=False)

	def LoadStackTree(self, request):
		if request.method == 'GET' and request.is_ajax():
			child_id = request.GET.get('id', '')
			q = StackES()

			if child_id == "#":				
				data = q.BuildRootTree()

				return JsonResponse(data, safe=False)

			else:
				data = q.BuildAuditAggs(child_id)
				
				return JsonResponse(data, safe=False)


class StackResponse(View, CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'stack_view.html', 
	    	context_instance=RequestContext(request, {"stack": self.nighthawk_stack, "stack_ver": self.nighthawk_stack_ver}))

	@method_decorator(csrf_protect)
	def post(self, request):
		stack_data = json.loads(request.body)
		s = StackES()

		data = s.GetAuditData(stack_data)
		
		return JsonResponse(data, safe=False)
