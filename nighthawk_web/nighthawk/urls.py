from django.conf.urls import url, include, patterns
from django.contrib import admin
from django.contrib.auth.decorators import login_required

from nighthawk.views.home import Home, HomeSearch
from nighthawk.views.update_doc import UpdateDoc
from nighthawk.views.upload import Upload
from nighthawk.views.comment import Comment
from nighthawk.views.stack_framework import StackView, StackResponse
from nighthawk.views.timeline import TimeLine, TimeLineResponse
from nighthawk.views.platform_stats import PlatformStats
from nighthawk.views.tasks import Tasks

from nighthawk.views.datatypes.w32registryraw import W32Registry
from nighthawk.views.datatypes.w32services import W32Services
from nighthawk.views.datatypes.filedownloadhistory import FiledownloadHistory
from nighthawk.views.datatypes.w32tasks import W32Tasks
from nighthawk.views.datatypes.w32network_route import W32Network_Route
from nighthawk.views.datatypes.urlhistory import UrlHistory
from nighthawk.views.datatypes.w32network_arp import W32Network_Arp
from nighthawk.views.datatypes.w32ports import W32Ports
from nighthawk.views.datatypes.w32prefetch import W32Prefetch
from nighthawk.views.datatypes.w32useraccounts import W32UserAccounts
from nighthawk.views.datatypes.w32network_dns import W32Network_DNS
from nighthawk.views.datatypes.w32scripting_filepersistence import W32ScriptingFilePersistence
from nighthawk.views.datatypes.stateagentinspector import StateagentInspector
from nighthawk.views.datatypes.w32processestree import W32ProcessesTree
from nighthawk.views.datatypes.w32volumes import W32Volumes
from nighthawk.views.datatypes.w32apifiles import W32APIFiles
from nighthawk.views.datatypes.w32rawfiles import W32RAWFiles
from nighthawk.views.datatypes.w32system import W32System
from nighthawk.views.datatypes.w32evtlog import W32EvtLog

data_types = patterns('',
    url(r'^filedownloadhistory_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', FiledownloadHistory.as_view(), name="filedownloadhistory"),
    url(r'^urlhistory_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', UrlHistory.as_view(), name="urlhistory"),
    url(r'^w32tasks_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Tasks.as_view(), name="w32tasks"),
    url(r'^w32services_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Services.as_view(), name="w32services"),
    url(r'^w32system_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32System.as_view(), name="w32system"),
    url(r'^w32registryraw_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Registry.as_view(), name="w32registryraw"),
    url(r'^w32apifiles_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32APIFiles.as_view(), name="w32apifiles"),
    url(r'^w32rawfiles_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32RAWFiles.as_view(), name="w32rawfiles"),
    url(r'^w32network-route_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Network_Route.as_view(), name="w32network_route"),
    url(r'^w32network-arp_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Network_Arp.as_view(), name="w32network_arp"),
    url(r'^w32ports_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Ports.as_view(), name="w32ports"),
    url(r'^w32prefetch_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Prefetch.as_view(), name="w32prefetch"),
    url(r'^w32useraccounts_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32UserAccounts.as_view(), name="w32useraccounts"),
    url(r'^w32network-dns_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Network_DNS.as_view(), name="w32network_dns"),
    url(r'^w32volumes_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32Volumes.as_view(), name="w32volumes"),
    url(r'^w32scripting-persistence_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32ScriptingFilePersistence.as_view(), name="w32network_dns"),
    url(r'^stateagentinspector_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', StateagentInspector.as_view(), name="stateagentinspector"),
    url(r'^w32processes-tree_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32ProcessesTree.as_view(), name="w32ptree"),
    url(r'^w32evtlogs_anchor/(?P<case>[^/]+)/(?P<hostname>\w+)$', W32EvtLog.as_view(), name="w32evtlogs"),
	)

context_urls = patterns('',
    url(r'^admin/', include(admin.site.urls)),
    url(r'^errorpage/$', login_required(Home().Home404), name='home_404'),
    url(r'^update_doc/$', login_required(UpdateDoc.as_view()), name="update_doc"),
    url(r'^comments/get_comment_doc/dialog$', login_required(Comment().DocDiaglog), name="comments_doc_dialog"),
    url(r'^comments/get_comment_doc/$', login_required(Comment().CommentDoc), name="comments_doc"),
    url(r'^comments/$', login_required(Comment.as_view()), name="comments"),
    url(r'^upload/$', login_required(Upload.as_view()), name="upload"),
    url(r'^home/load_cases$', login_required(Home().LoadCaseTree), name="load_cases"),
    url(r'^home/load_cases_audit$', login_required(Home().LoadCaseTreeAudit), name="load_cases_audit"),
    url(r'^home/load_stack$', login_required(StackView().LoadStackTree), name="load_stack"),
    url(r'^home/load_timeline$', login_required(TimeLine().LoadTLTree), name="load_timeline"),
    url(r'^home/main_search/$', login_required(HomeSearch.as_view()), name="home_search"),
    url(r'^stack_response/$', login_required(StackResponse.as_view()), name="stack_response"),
    url(r'^stack/$', login_required(StackView.as_view()), name="stack_framework"),
    url(r'^timeline_response/$', login_required(TimeLineResponse.as_view()), name="timeline_response"),
    url(r'^timeline/$', login_required(TimeLine.as_view()), name="timeline"),
    url(r'^platform_stats/$', login_required(PlatformStats.as_view()), name="platform_stats"),    
    url(r'^tasks/analyst$', login_required(Tasks().GetAnalyst), name="task_analyst"),
    url(r'^tasks/active$', login_required(Tasks().GetActiveTasks), name="task_active"),
    url(r'^tasks/$', login_required(Tasks.as_view()), name="tasks"),    
    url(r'^delete_case/$', login_required(Upload().DeleteCaseList), name="del_case_list"),   
    url(r'^delete_endpoint/$', login_required(Upload().DeleteEndpointList), name="del_endpoint_list"),
    url(r'^delete/$', login_required(Upload().Delete), name="del_endpoint_list"),    
    url(r'$', login_required(Home.as_view()), name="home"),
    )

urlpatterns = [
    url(r'^', include(data_types)),
    url(r'^', include(context_urls))
]
