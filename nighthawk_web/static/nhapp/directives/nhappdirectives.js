(function(){

var stack_data = {};

angular
	.module('customDirectives', [])

	.directive('jstree', function() {
		return {
			restrict: 'E',
            scope: {
            	dataurl: '=',
                dataurlaudit: '=',
                selectnode: '=',
            	url: '=',
            	themes: '=',
            	types: '=',
            	plugins: '=',
            	check: '=',
            },
            template: '<div id="jstree"></div>',
            replace: true,
            link: function(scope, elem){
                    var jstree = $('#jstree').jstree({
                        'core' : {
                            'data' : {
                                'url' : function (node) {
                                    if (node.parent == null){
                                        return scope.dataurl
                                    } else if (node.parent == 'current_inv' || node.parent == 'timeline' || node.parent == 'stackable'){
                                        return scope.dataurl
                                    } else {
                                        return scope.dataurlaudit
                                    }

                                },
                                'contentType': "application/json; charset=utf-8",
                                'data': function (data) {
                                    return { "id": data.id, "parent": data.parent, "text": data.text }
                                  }
                            },
                            'themes': {
                                'name': 'default-dark'
                            }
                        },
                      'types': scope.types,
                        "checkbox" : {
                            "keep_selected_style" : false
                        },
                        "plugins" : scope.plugins
                    })
                .on('open_node.jstree', function (event, data) {
                    if (scope.selectnode === 'mainview') {
                        data.instance.set_icon(data.node, "glyphicon glyphicon-minus")
                    }
                })
                .on('select_node.jstree', function (event, data) {

                    switch (scope.selectnode) {
                        
                        case "mainview": 
                            event.preventDefault()
                            var parent = data.node.parent
                            var href = data.node.a_attr.id.split("_")[0] + '/' + data.node.parents[1] + '/' + data.node.parent
                            document.location.href = "#" + href
                            break;
                        
                        case "stackview":

                            event.preventDefault()
                            var parent = data.node.id
                            var audittype = data.node.a_attr.id.split("_")[0]

                            if (!isInDict(parent, stack_data)) {
                                stack_data[parent] = data.node.children
                                localStorage.setItem('stack_data', JSON.stringify(stack_data))
                            }
                            else {
                                // do nothing
                            }
                            break;
                        
                        case "timelineview":
                            var timeline_data = {};
                            event.preventDefault()
                            var endpoint = data.node.id

                            timeline_data = {
                                "case": data.node.parent,
                                "endpoint": endpoint
                            };
                            
                            localStorage.setItem('timeline_data', JSON.stringify(timeline_data))
                            break;

                     }
                })
                .on('deselect_node.jstree', function (event, data) {
                    event.preventDefault()
                    var audittype = data.node.a_attr.id.split("_")[0]

                    if (audittype in stack_data) {
                        delete stack_data[audittype]
                        localStorage.setItem('stack_data', JSON.stringify(stack_data))
                    }
                    else {
                        // do nothing
                    }
                });

                return jstree;
             }
         }
	})

})();

