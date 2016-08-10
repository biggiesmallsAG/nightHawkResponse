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

	.controller('deleteController', ['$scope', 'getApiServices', 'postApiServices', 'ngDialog', function($scope, getApiServices, postApiServices, ngDialog) {

		$scope.cases = getApiServices.query({api_path: 'delete_case/'});
		$scope.endpoints = getApiServices.query({api_path: 'delete_endpoint/'});

		$scope.update = function(){
			if (angular.isUndefined($scope.case)) {
				$scope.selected = {
					"case": "",
					"endpoint": $scope.endpoint.endpoint,
					"singular": false
				};				
			}
			else if (angular.isUndefined($scope.endpoint)) {
				$scope.selected = {
					"case": $scope.case.case,
					"endpoint": "",
					"singular": false
				};					
			}
			else {
				$scope.selected = {
					"case": $scope.case.case,
					"endpoint": $scope.endpoint.endpoint,
					"singular": true
				};			
			}

		$scope.postDelete = function (selected) {
				$('.input_hide').hide(function(){
					$('.loader').show();
				});
				postApiServices.save({api_path: 'delete'}, selected, 
					function (data) {
						ngDialog.open({
							template: "<span style='color:#4AAAC4'>Successfull delete took: </span><span style='color:#F7973D'>" + data.took + "ms</span>",
							plain: true
						});
						$('.loader').hide(function(){
							$('.input_hide').show()
						})
				},
					function (err) {
						ngDialog.open({
							template: "<span style='color:#4AAAC4'>Error in delete, reason: </span><span style='color:#F7973D'>" + err.data.reason + "</span>",
							plain: true
						});
						$('.loader').hide(function(){
							$('.input_hide').show()
						})
					})
			}
		};

	}])

})();