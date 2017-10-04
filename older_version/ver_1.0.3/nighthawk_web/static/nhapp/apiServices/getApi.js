(function(){

var API_PATH = '/';

angular
	.module('getApi', [])

	.factory('getApiServices', ['$resource', function($resource) {

			return $resource(API_PATH+':api_path', {api_path: '@api_path'});
				
	}])

	.factory('getApiServicesObject', ['$resource', function($resource) {

			return $resource(API_PATH+':api_path', {api_path: '@api_path'}, { query: {
					isArray: false
					}
				});
				
	}])

})();

