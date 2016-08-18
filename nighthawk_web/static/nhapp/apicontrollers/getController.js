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

	.controller('mainTreeController', ['$scope', 'getApiServices', function($scope, getApiServices) {

		$scope.dataurl = "home/load_cases";

		$scope.dataurlaudit = "home/load_cases_audit";

		$scope.types = {
                        'root': {
                          'icon': 'glyphicon glyphicon-expand'
                        },
                        'endpoint': {
                          'icon': 'glyphicon glyphicon-blackboard'
                        },
                        'case': {
                          'icon': 'glyphicon glyphicon-leaf'
                        },
                        'audit': {
                          'icon': 'glyphicon glyphicon-log-in'
                        }
                      }

        $scope.plugins = ["sort", "types"];

        $scope.selectnode = "mainview"

	}])

	.controller('timelineTreeController', ['$scope', 'getApiServices', function($scope, getApiServices) {

		$scope.dataurl = "/home/load_timeline";

		$scope.types = {
	                'root': {
	                  'icon': 'glyphicon glyphicon-expand'
	                },
	                'endpoint': {
	                  'icon': 'glyphicon glyphicon-blackboard'
	                },
	                'case': {
	                  'icon': 'glyphicon glyphicon-leaf'
	                },
	                'audit': {
	                  'icon': 'glyphicon glyphicon-log-in'
	                }
	              }

        $scope.plugins = ["sort", "types"];

        $scope.selectnode = "timelineview"

	}])

	.controller('stackTreeController', ['$scope', 'getApiServices', function($scope, getApiServices) {

		$scope.dataurl = "/home/load_stack";

		$scope.types = {
	            'root': {
	              'icon': 'glyphicon glyphicon-expand'
	            },
	            'endpoint': {
	              'icon': 'glyphicon glyphicon-blackboard'
	            },
	            'stack': {
	              'icon': 'glyphicon glyphicon-fire'
	            },
	            'audit': {
	              'icon': 'glyphicon glyphicon-log-in'
	            }
	          }

        $scope.plugins = ["sort", "types", "checkbox"];

        $scope.selectnode = "stackview"

	}])

	.controller('taskController', ['$scope', 'getApiServicesObject', 'getApiServices', 'postGetArray', function($scope, getApiServicesObject, getApiServices, postGetArray) {

		$scope.tasks = getApiServices.query({api_path: 'tasks/active'});
		$scope.analyst = getApiServicesObject.query({api_path: 'tasks/analyst'});
		$scope.assignee = $scope.analyst;

		$scope.addTask = function() {
			$scope.tasks.push({
				"datetime": "Pending save..",
				"id": "-",
				"endpoints": $scope.endpoints,
				"task_name": $scope.task_name,
				"urgency": $scope.urgency,
				"task_analyst": $scope.analyst.analyst,
				"task_assignee": $scope.assigned,
				"task_inactive": false,
				"task_done": false
			});

			$scope.task_name = '';
			$scope.endpoints = '';
			$scope.urgency = '';

		};

		$scope.update = function (){
			
		};

		$scope.saveTasks = function(){
			for (var i = 0; i < $scope.tasks.length; i++) {
				if ($scope.tasks[i].id === '-') {
					var newTask = [];
					newTask.push($scope.tasks[i])

					postGetArray.save({api_path: 'tasks'}, newTask, function (success){
						$scope.tasks = getApiServices.query({api_path: 'tasks/active'});
					})
				}
			}
		};

		$scope.removeTask = function(){
			for (var i = 0; i < $scope.tasks.length; i++) {
				if ($scope.tasks[i].task_inactive) {
					var removeTasks = [];
					removeTasks.push($scope.tasks[i])

					postGetArray.save({api_path: 'tasks'}, removeTasks, function (success){
						$scope.tasks = getApiServices.query({api_path: 'tasks/active'});
					})
				}
			}
		};

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