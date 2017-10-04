(function(){  

app = angular

	.module('brsapp', 
		[
		'ngRoute', 
		'ngResource',
		'ngDialog',
		'ngCookies',
		'getController', 
		'postController', 
		'getApi', 
		'postApi',
		'datatables',
		'datatables.colreorder',
		'datatables.buttons',
		'datatables.bootstrap',
		'datatables.columnfilter',
		'jsonFormatter',
		'customDirectives'
		],
	function($interpolateProvider, $httpProvider, $resourceProvider) 
	{
		$httpProvider.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
		$httpProvider.defaults.xsrfCookieName = 'csrftoken';
		$httpProvider.defaults.xsrfHeaderName = 'X-CSRFToken';
		$httpProvider.defaults.headers.common['X-CSRFToken'] = '{{ csrf_token }}';
	    $interpolateProvider.startSymbol('[[');
	    $interpolateProvider.endSymbol(']]');
	    $resourceProvider.defaults.stripTrailingSlashes = false;
	})

	.config(['$routeProvider', function($routeProvider){

		$routeProvider.

		when('/admin', {
			templateUrl: function() {
				return '/admin/'
			}
		}).

		when('/Tasks', {
			templateUrl: function()  {
				return '/tasks/'
			}
		}).

		when('/UploadFile', {
			templateUrl: '/upload/'
		}).

		when('/PlatformStats', {
			templateUrl: '/platform_stats/'
		}).

		when('/w32system/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32system_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			}
		}).

		when('/w32eventlogs/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32evtlogs_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			}
		}).

		when('/w32registryraw/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32registryraw_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32registry'
		}).

		when('/w32apifiles/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32apifiles_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			}
		}).

		when('/w32rawfiles/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32rawfiles_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			}
		}).

		when('/w32services/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32services_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32services'
		}).

		when('/filedownloadhistory/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/filedownloadhistory_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'filedownloadhistory'
		}).

		when('/w32tasks/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32tasks_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32tasks'
		}).

		when('/w32network-route/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32network-route_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32network_route'
		}).

		when('/urlhistory/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/urlhistory_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'urlhistory'
		}).

		when('/w32network-arp/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32network-arp_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32network_arp'
		}).

		when('/w32ports/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32ports_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32ports'
		}).

		when('/w32prefetch/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32prefetch_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32prefetch'
		}).

		when('/w32useraccounts/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32useraccounts_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32useraccounts'
		}).

		when('/w32network-dns/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32network-dns_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32network_dns'
		}).

		when('/w32scripting-persistence/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32scripting-persistence_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32network_dns'
		}).

		when('/stateagentinspector/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/stateagentinspector_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32network_dns'
		}).

		when('/w32processes-tree/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32processes-tree_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'w32network_dns'
		}).

		when('/w32volumes/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/w32volumes_anchor/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: ''
		}).

		when('/timeline/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/timeline/' + urlattr.casename + '/' + urlattr.hostname;
			},
			controller: 'timeline'
		}).

		when('/stack_response/', {
			templateUrl: function(urlattr) {
				return '/stack_response/';
			},
			controller: 'stackView'
		}).

		when('/Comment', {
			templateUrl: function(urlattr) {
				return '/comments/';
			},
			controller: 'commentController'
		}).

		when('/w32disks/:casename/:hostname', {
			templateUrl: function(urlattr) {
				return '/errorpage/';
			}
		}).

		otherwise('/', {
			templateUrl: function(urlattr) {
				return '/'
			}
		})
	}])

})();
