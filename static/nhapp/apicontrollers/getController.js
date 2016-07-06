(function(){

angular
	.module('getController', [])

	.controller('loadController', ['$scope', 'ngDialog', function($scope, ngDialog) {

		    $scope.$on('LOAD', function() {
		    	$scope.loading = true;
		    });
		    $scope.$on('UNLOAD', function() {
		    	$scope.loading = false;
		    })

	}])

	.controller('homeView', ['$scope', function($scope) {

			$scope.data = "Welcome to the nightHawk Response framework. Use the node tree on the left to select an endpoint."

	}])

	.controller('uploadFile', ['$scope', function($scope) {

			$scope.data = ""

	}]);


})();