(function(){

angular
	.module('postApi', [])
	
	.factory('postApiServices', ['$resource', function($resource) {

		return $resource(':api_path/', 
				{api_path: '@api_path'});

	}])

	.factory('postGetArray', ['$resource', function($resource) {

		return $resource(':api_path/', 
				{api_path: '@api_path'}, { query: {
					method: 'POST',
					isArray: true
				}
				});

	}]);
	
})();

