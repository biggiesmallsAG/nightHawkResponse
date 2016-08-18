from django import forms
from django.utils import html

class UploadForm(forms.Form):

	CONCURRENT_CHOICES = (
		('0', 'Yes'),
		('1', 'No')
		)

	upload_field = forms.FileField(label="Select .zip file", widget=forms.ClearableFileInput(attrs={'multiple': True}))
	case_number = forms.CharField(label="Case Number", widget=forms.TextInput(attrs={'class': 'form-control', 'style': 'width: 70%','placeholder': 'ex. XXX-123'}))
	concurrent = forms.ChoiceField(label="Concurrent (Multi Processing)?", choices=CONCURRENT_CHOICES)

class UpdateDoc(forms.Form):

	LENGTH_CHOICES = (
		('0', 'True'),
		('1', 'False'),
		)

	TAG_CHOICES = (
		('0', 'Benign'),
		('1', 'Follow Up'),
		('2', 'Malicious'),
		('3', 'For Review')
		)
	form_name = "UpdateDocForm"
	update_comment = forms.CharField(label="Comment", widget=forms.TextInput(attrs={'class': 'form-control', 'style': 'width: 70%', 'ng-model': 'comment', 'ng-required': 'true'}))
	update_fp = forms.ChoiceField(label="False Positive", choices=LENGTH_CHOICES, required=False, widget=forms.Select(attrs={'class': 'form-control', 'style': 'width: 30%', 'ng-model': 'falsepositive', 'ng-required': 'true'}))
	update_tag = forms.ChoiceField(label="Tag", choices=TAG_CHOICES, required=False, widget=forms.Select(attrs={'class': 'form-control', 'style': 'width: 30%', 'ng-model': 'tag', 'ng-required': 'true'}))

class SearchForm(forms.Form):
	search = forms.CharField(label="", widget=forms.TextInput(attrs={'class': 'form-control', 'style': 'width: 30%','placeholder': 'ex. .exe OR .dll'}))

class DeleteCase(forms.Form):

	delete_case = forms.ChoiceField(label="Delete Case", widget=forms.Select(attrs={
		'ng-model': 'case', 
		'ng-change': 'update()', 
		'ng-options': "casex.case for casex in cases| orderBy:'name'"
		}))

class DeleteEndpoint(forms.Form):

	delete_endpoint = forms.ChoiceField(label="Delete Endpoint", widget=forms.Select(attrs={
		'ng-model': 'endpoint', 
		'ng-change': 'update()', 
		'ng-options': "endpointx.endpoint for endpointx in endpoints| orderBy:'name'"
		}))

class TasksForm(forms.Form):
	URGENCY_CHOICES = (
		('urgent', 'Urgent'),
		('one_hour', 'One Hour'),
		('five_hours', 'Five Hours')
		)

	task_name = forms.CharField(label="Task", widget=forms.TextInput(attrs={
			"class": "form-control",
			"placeholder": "Task name?",
			"ng-model": "task_name",
			"ng-required": "true"
		}))

	task_endpoint = forms.CharField(label="Endpoint", widget=forms.TextInput(attrs={
			"class": "form-control",
			"placeholder": "Endpoint(s)",
			"ng-model": "endpoints",
			"ng-required": "true"
		}))

	task_urgency = forms.ChoiceField(label="Urgency", widget=forms.RadioSelect(attrs={
			"ng-model": "urgency",
			"ng-required": "true"
		}), choices=URGENCY_CHOICES)

	task_assignee = forms.ChoiceField(label="Assigned To", widget=forms.Select(attrs={
			"ng-options": "x for x in assignee.assignee",
			"ng-model": "assigned",
			"ng-change": "update()",
			"ng-required": "true"
		}))

