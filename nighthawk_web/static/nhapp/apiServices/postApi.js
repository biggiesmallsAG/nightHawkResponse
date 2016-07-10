(function(){

var BrsTriage = "http://localhost:9090/";

angular
	.module('postApi', [])
	
	.factory('postApiServices', ['$resource', function($resource) {

		return $resource(BrsTriage+':api_path/', 
				{api_path: '@api_path'});

	}])

	.factory('postGetArray', ['$resource', function($resource) {

		return $resource(BrsTriage+':api_path/', 
				{api_path: '@api_path'}, { query: {
					method: 'POST',
					isArray: true
				}
				});

	}]);
	
})();

