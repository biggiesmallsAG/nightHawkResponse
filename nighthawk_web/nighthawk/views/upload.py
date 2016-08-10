from django.views.generic.edit import View
from django.shortcuts import render
from django.core.urlresolvers import reverse
from django.http import HttpResponse, HttpResponseRedirect, HttpRequest, JsonResponse
from django.template import loader, RequestContext
from django.contrib.auth.decorators import login_required
from django.views.decorators.csrf import ensure_csrf_cookie
from nighthawk.triageapi.dataendpoint.common import CommonAttributes
from nighthawk.triageapi.utility.validate import ValidateUserInput
from nighthawk.triageapi.delete import DeleteCaseObj
from nighthawk.forms import UploadForm, DeleteCase, DeleteEndpoint
from nighthawk import settings
import subprocess
import json

class Upload(View, CommonAttributes):
	def __init__(self):
		CommonAttributes.__init__(self)

	def get(self, request):
	    assert isinstance(request, HttpRequest)
	    return render(request, 'upload.html', 
	    	context_instance=RequestContext(request, {
	    		"request": request, 
	    		"UploadForm": UploadForm, 
	    		"DeleteEndpoint": DeleteEndpoint,
	    		"DeleteCase": DeleteCase
	    		}))

	def post(self, request):
		if request.method == 'POST':
			form = UploadForm(request.POST, request.FILES)
			files = request.FILES.getlist('upload_field')
			case_id = request.POST.get('case_number')
			concurrent = request.POST.get('concurrent')

			if ValidateUserInput(case_id).ValidateInputMixedPunctual() or case_id == '':
				if concurrent == '0':
					if len(files) > self.file_upload_max:
						return HttpResponse("Max Uploads concurrency is set to {0}. This can be reset in nightHawk.json".format(self.file_upload_max))

					for f in files:
						self.process_files(f)

					processes = []
					for f in files:
						if case_id:
							processes.append(subprocess.Popen([settings.NIGHTHAWK_GO + "/./nightHawk", "-R", "-v", "-N", "{0}".format(case_id), "-f", "{0}/{1}".format(settings.MEDIA_DIR ,f)], stderr=subprocess.STDOUT))
						else:
							processes.append(subprocess.Popen([settings.NIGHTHAWK_GO + "/./nightHawk", "-R", "-v", "-f", "{0}/{1}".format(settings.MEDIA_DIR ,f)], stderr=subprocess.STDOUT))
					
					for p in processes:
						if p.poll() is not None:
							if p.returncode == 0:
								processes.remove(p)

				else:
					for f in files:
						self.process_files(f)

					for f in files: 
						if case_id:
							subprocess.call([settings.NIGHTHAWK_GO + "/./nightHawk", "-R", "-v", "-N", "{0}".format(case_id), "-f", "{0}/{1}".format(settings.MEDIA_DIR ,f)], stderr=subprocess.STDOUT)
						else:
							subprocess.call([settings.NIGHTHAWK_GO + "/./nightHawk", "-R", "-v", "-f", "{0}/{1}".format(settings.MEDIA_DIR ,f)], stderr=subprocess.STDOUT)
								
				return HttpResponseRedirect('/upload')

			else:
				return HttpResponseRedirect('/upload')

	def process_files(self, f):
	    with open(settings.MEDIA_DIR + '/{0}'.format(f), 'wb+') as destination:
	        for chunk in f.chunks():
	            destination.write(chunk)
	

	def DeleteCaseList(self, request):
		if request.is_ajax():
			d = DeleteCaseObj()
			data = d.BuildCaseList()
			return JsonResponse(data, safe=False)

	def DeleteEndpointList(self, request):
		if request.is_ajax():
			d = DeleteCaseObj()
			data = d.BuildEndpointList()
			return JsonResponse(data, safe=False)

	def Delete(self, request):
		if request.is_ajax():
			data = json.loads(request.body)
			d = DeleteCaseObj()
			data = d.Delete(data)
			return JsonResponse(data, safe=False)
