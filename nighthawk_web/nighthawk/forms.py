from django import forms
from django.utils import html

class UploadForm(forms.Form):

	CONCURRENT_CHOICES = (
		('0', 'Yes'),
		('1', 'No')
		)

	upload_field = forms.FileField(label="Select .zip file", widget=forms.ClearableFileInput(attrs={'multiple': True}))
	case_number = forms.CharField(label="Case Number", widget=forms.TextInput(attrs={'class': 'form-control', 'style': 'width: 20%','placeholder': 'ex. XXX-123'}))
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
	update_date = forms.DateField(label="Date", widget=forms.TextInput(attrs={'class': 'form-control datepicker', 'style': 'width: 50%', 'ng-model': 'date', 'ng-required': 'true'}))
	update_fp = forms.ChoiceField(label="False Positive", choices=LENGTH_CHOICES, required=False, widget=forms.Select(attrs={'class': 'form-control', 'style': 'width: 30%', 'ng-model': 'falsepositive', 'ng-required': 'true'}))
	update_tag = forms.ChoiceField(label="Tag", choices=TAG_CHOICES, required=False, widget=forms.Select(attrs={'class': 'form-control', 'style': 'width: 30%', 'ng-model': 'tag', 'ng-required': 'true'}))


class SearchForm(forms.Form):
	search = forms.CharField(label="", widget=forms.TextInput(attrs={'class': 'form-control', 'style': 'width: 30%','placeholder': 'ex. .exe OR .dll'}))