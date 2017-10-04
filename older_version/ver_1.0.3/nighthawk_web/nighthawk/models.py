from django.db import models

class TaskModel(models.Model):
	
	task_date = models.DateTimeField(max_length=30)
	task_name = models.CharField(max_length=120)
	task_analyst = models.CharField(max_length=50)
	task_assignee = models.CharField(max_length=50, null=True)
	task_endpoints = models.TextField(max_length=200)
	task_urgency = models.CharField(max_length=20, null=True)
	task_inactive = models.BooleanField()