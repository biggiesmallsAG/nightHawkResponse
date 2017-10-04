(function(){
	var nHResponse = "/";

	angular
	.module('postController', [])

	.controller('w32system', GetDataSystem)

	.controller('w32registry', GetDataRegistry)

	.controller('w32apifiles', GetDataApiFiles)

	.controller('w32rawfiles', GetDataRawFiles)

	.controller('w32services', GetDataServices)

	.controller('filedownloadhistory', GetDataFileDownload)

	.controller('w32tasks', GetDataTasks)

	.controller('w32network_route', GetDataNetworkRoute)

	.controller('urlhistory', GetDataUrlHistory)

	.controller('w32network_arp', GetDataNetworkArp)

	.controller('w32ports', GetDataPorts)

	.controller('w32prefetch', GetDataPrefetch)

	.controller('w32useraccounts', GetDataUserAccounts)

	.controller('w32network_dns', GetDataDNS)

	.controller('w32scripting-persistence', GetDataPersistence)

	.controller('w32volumes', GetDataVolumes)

	.controller('w32disks', GetDataDisks)

	.controller('w32evtlogs', GetDataEvtLogs)

	.controller('stackView', GetStackResponse)

	.controller('searchAllDocs', GetSearchAllDocs)

	.controller('commentController', CommentView)

	.controller('confirmController', ['$scope', function($scope){

		$scope.comment = '';
		$scope.falsepositive = false;
		$scope.tag = '';

		$scope.submitUpdate = function() {
			return {
				comment: $scope.comment,
				falsepositive: $scope.falsepositive,
				tag: $scope.tag, 
				rowId: $scope.ngDialogData._id,
				parent: $scope.ngDialogData._parent
			}
		}

	}]);

	function GetSearchAllDocs($scope) {
		$scope.search_main = 'To search all documents:'
	}

	function GetDataSystem($scope, ngDialog, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32system_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withBootstrap()
	    .withOption('rowCallback', rowCallback)
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.SystemInfo.BiosVersion').withTitle('BIOS Version'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.Drives').withTitle('Drives'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.InstallDate').withTitle('Install Date'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.Domain').withTitle('Domain'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.LoggedOnUser').withTitle('Logged On User'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.Mac').withTitle('Mac Addr'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.OS').withTitle('OS'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.OsBitness').withTitle('OS (32/64bit)'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.PrimaryIpAddress').withTitle('Primary IP'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.TimezoneStandard').withTitle('TimezoneStandard'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.GmtOffset').withTitle('GMTOffset'),
		    DTColumnBuilder.newColumn('_source.SystemInfo.TotalPhysical').withTitle('TotalPhysical'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataRegistry($scope, ngDialog, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, $cookies) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;

	    vm.dtOptions = DTOptionsBuilder.newOptions()
	        .withOption('ajax', {
	         url: nHResponse + 'w32registryraw_anchor/' + $routeParams.casename + '/' + $routeParams.hostname,
	         type: 'POST',
	         data: {
	         	data: 'data',
	         	csrfmiddlewaretoken: $cookies.get('csrftoken')
	         }
	     })
        .withOption('processing', true)
        .withOption('serverSide', true)
        .withLanguage({
        	"sProcessing": '<img src="static/images/loader.gif">'
        })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    .withOption('order', [0, 'desc'])
	    .withColReorder()
	    .withPaginationType('numbers')
	    .withDisplayLength(200)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);

	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Modified').withTitle('Modified'),
		    DTColumnBuilder.newColumn('_source.Record.Path').withTitle('Path'),
		    DTColumnBuilder.newColumn('_source.Record.Type').withTitle('Type'),
		    DTColumnBuilder.newColumn('_source.Record.Text').withTitle('Text'),
		    DTColumnBuilder.newColumn('_source.Record.ReportedLengthInBytes').withTitle('Bytes'),
		    DTColumnBuilder.newColumn('_source.Record.Hive').withTitle('Hive'),
		    DTColumnBuilder.newColumn('_source.Record.Username').withTitle('Username'),
		    DTColumnBuilder.newColumn('_source.Record.JobCreated').withTitle('JobCreated'),
		    DTColumnBuilder.newColumn('_source.Record.ValueName').withTitle('ValueName'),
		    DTColumnBuilder.newColumn('_source.Record.IsKnownKey').withTitle('IsKnownKey'),
		    DTColumnBuilder.newColumn('_source.Record.KeyPath').withTitle('KeyPath'),
		    DTColumnBuilder.newColumn('_source.Record.SecurityID').withTitle('SecurityID'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];

	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataEvtLogs($scope, ngDialog, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, $cookies) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;

	    vm.dtOptions = DTOptionsBuilder.newOptions()
	        .withOption('ajax', {
	         url: nHResponse + 'w32evtlogs_anchor/' + $routeParams.casename + '/' + $routeParams.hostname,
	         type: 'POST',
	         data: {
	         	data: 'data',
	         	csrfmiddlewaretoken: $cookies.get('csrftoken')
	         }
	     })
        .withOption('processing', true)
        .withOption('serverSide', true)
        .withLanguage({
        	"sProcessing": '<img src="static/images/loader.gif">'
        })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    .withOption('order', [0, 'desc'])
	    .withColReorder()
	    .withPaginationType('numbers')
	    .withDisplayLength(200)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);

	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.GenTime').withTitle('Generate Time'),
		    DTColumnBuilder.newColumn('_source.Record.Source').withTitle('Source'),
		    DTColumnBuilder.newColumn('_source.Record.EID').withTitle('EID'),
		    DTColumnBuilder.newColumn('_source.Record.Index').withTitle('Index'),
		    DTColumnBuilder.newColumn('_source.Record.ExecutionProcessId').withTitle('ExecutionProcessId'),
		    DTColumnBuilder.newColumn('_source.Record.ExecutionThreadId').withTitle('ExecutionThreadId'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.AccountDomain').withTitle('AccountDomain'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.AccountName').withTitle('AccountName'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LinkedLogonId').withTitle('LinkedLogonId'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LogonId').withTitle('LogonId'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.NetworkAccountDomain').withTitle('NetworkAccountDomain'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.NetworkAccountName').withTitle('NetworkAccountName'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.SecurityId').withTitle('SecurityID'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.AuthenticationPackage').withTitle('AuthenticationPackage'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.KeyLength').withTitle('KeyLength'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LogonInformation.ElevatedToken').withTitle('ElevatedToken'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LogonInformation.LogonType').withTitle('LogonType'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LogonInformation.LogonTypeDetail').withTitle('LogonTypeDetail'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LogonInformation.RestrictedAdminMode').withTitle('RestrictedAdminMode'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LogonInformation.VirtualAccount').withTitle('VirtualAccount'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.LogonProcess').withTitle('LogonProcess'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.NetworkInformation.SourceNetworkAddress').withTitle('SourceNetworkAddress'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.NetworkInformation.SourcePort').withTitle('SourcePort'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Account.NetworkInformation.WorkstationName').withTitle('WorkstationName'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.PackageName').withTitle('PackageName'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.ProcessInformation.ProcessCommandline').withTitle('ProcessCommandline'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.ProcessInformation.ProcessId').withTitle('ProcessId'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.ProcessInformation.ProcessIdHex').withTitle('ProcessIdHex'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.ProcessInformation.ProcessName').withTitle('ProcessName'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.AccountDomain').withTitle('AccountDomain'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.AccountName').withTitle('AccountName'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.LinkedLogonId').withTitle('LinkedLogonId'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.LogonGuid').withTitle('LogonGuid'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.LogonId').withTitle('LogonId'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.NetworkAccountDomain').withTitle('NetworkAccountDomain'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.NetworkAccountName').withTitle('NetworkAccountName'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.Subject.SecurityId').withTitle('SecurityId'),
			DTColumnBuilder.newColumn('_source.Record.MessageDetail.TransitedServices').withTitle('TransitedServices'),
		    DTColumnBuilder.newColumn('_source.Record.Message').withTitle('Message'),
		    DTColumnBuilder.newColumn('_source.Record.Category').withTitle('Category'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
	    ];
	
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataApiFiles($scope, ngDialog, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, $cookies) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;

	    vm.dtOptions = DTOptionsBuilder.newOptions()
	        .withOption('ajax', {
	         url: nHResponse + 'w32apifiles_anchor/' + $routeParams.casename + '/' + $routeParams.hostname,
	         type: 'POST',
	         data: {
	         	data: 'data',
	         	csrfmiddlewaretoken: $cookies.get('csrftoken')
	         }
	     })
        .withOption('processing', true)
        .withOption('serverSide', true)
        .withLanguage({
        	"sProcessing": '<img src="static/images/loader.gif">'
        })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withOption('order', [0, 'desc'])
	    .withColReorder()
	    .withPaginationType('numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);

	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Modified').withTitle('Modified'),
		    DTColumnBuilder.newColumn('_source.Record.Accessed').withTitle('Accessed'),
		    DTColumnBuilder.newColumn('_source.Record.Changed').withTitle('Changed'),
		    DTColumnBuilder.newColumn('_source.Record.Username').withTitle('User'),
		    DTColumnBuilder.newColumn('_source.Record.Path').withTitle('Path'),
		    DTColumnBuilder.newColumn('_source.Record.SizeInBytes').withTitle('SizeInBytes'),
		    DTColumnBuilder.newColumn('_source.Record.Md5sum').withTitle('Md5sum'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodHash').withTitle('IsGoodHash'),
		    DTColumnBuilder.newColumn('_source.Record.FileAttributes').withTitle('FileAttributes'),
			DTColumnBuilder.newColumn('_source.Record.PeInfo.Type').withTitle('File.PeInfo.Type'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.Subsystem').withTitle('File.PeInfo.Subsystem'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.BaseAddress').withTitle('File.PeInfo.BaseAddress'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.PETimeStamp').withTitle('File.PeInfo.PETimeStamp'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.SignatureExists').withTitle('File.DigitalSignature.SignatureExists'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.SignatureVerified').withTitle('File.DigitalSignature.SignatureVerified'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.Description').withTitle('File.DigitalSignature.Description'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.CertificateSubject').withTitle('File.DigitalSignature.CertificateSubject'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.CertificateIssuer').withTitle('File.DigitalSignature.CertificateIssuer'),
		    DTColumnBuilder.newColumn('_source.Record.JobCreated').withTitle('JobCreated'),
		    DTColumnBuilder.newColumn('_source.Record.SecurityID').withTitle('SecurityID'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataRawFiles($scope, ngDialog, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, $cookies) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;

	    vm.dtOptions = DTOptionsBuilder.newOptions()
	        .withOption('ajax', {
	         url: nHResponse + 'w32rawfiles_anchor/' + $routeParams.casename + '/' + $routeParams.hostname,
	         type: 'POST',
	         data: {
	         	data: 'data',
	         	csrfmiddlewaretoken: $cookies.get('csrftoken')
	         }
	     })
        .withOption('processing', true)
        .withOption('serverSide', true)
        .withLanguage({
        	"sProcessing": '<img src="static/images/loader.gif">'
        })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withOption('order', [0, 'desc'])
	    .withColReorder()
	    .withPaginationType('numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);

	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.FilenameCreated').withTitle('FNCreated'),
		    DTColumnBuilder.newColumn('_source.Record.FilenameModified').withTitle('FNModified'),
		    DTColumnBuilder.newColumn('_source.Record.FilenameAccessed').withTitle('FNAccessed'),
		    DTColumnBuilder.newColumn('_source.Record.FilenameChanged').withTitle('FNChanged'),
		    DTColumnBuilder.newColumn('_source.Record.Created').withTitle('SICreated'),
		    DTColumnBuilder.newColumn('_source.Record.Modified').withTitle('SIModified'),
		    DTColumnBuilder.newColumn('_source.Record.Accessed').withTitle('SIAccessed'),
		    DTColumnBuilder.newColumn('_source.Record.Changed').withTitle('SIChanged'),
		    DTColumnBuilder.newColumn('_source.Record.Username').withTitle('User'),
		    DTColumnBuilder.newColumn('_source.Record.Path').withTitle('Path'),
		    DTColumnBuilder.newColumn('_source.Record.FileName').withTitle('FileName'),
		    DTColumnBuilder.newColumn('_source.Record.FileExtension').withTitle('Extension'),
		    DTColumnBuilder.newColumn('_source.Record.SizeInBytes').withTitle('SizeInBytes'),
		    DTColumnBuilder.newColumn('_source.Record.Md5sum').withTitle('Md5sum'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodHash').withTitle('IsGoodHash'),
		    DTColumnBuilder.newColumn('_source.Record.FileAttributes').withTitle('FileAttributes'),
			DTColumnBuilder.newColumn('_source.Record.PeInfo.Type').withTitle('File.PeInfo.Type'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.Subsystem').withTitle('File.PeInfo.Subsystem'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.BaseAddress').withTitle('File.PeInfo.BaseAddress'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.PETimeStamp').withTitle('File.PeInfo.PETimeStamp'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.SignatureExists').withTitle('File.DigitalSignature.SignatureExists'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.SignatureVerified').withTitle('File.DigitalSignature.SignatureVerified'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.Description').withTitle('File.DigitalSignature.Description'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.CertificateSubject').withTitle('File.DigitalSignature.CertificateSubject'),
		    DTColumnBuilder.newColumn('_source.Record.PeInfo.DigitalSignature.CertificateIssuer').withTitle('File.DigitalSignature.CertificateIssuer'),
		    DTColumnBuilder.newColumn('_source.Record.JobCreated').withTitle('JobCreated'),
		    DTColumnBuilder.newColumn('_source.Record.SecurityID').withTitle('SecurityID'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataServices($scope, ngDialog, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32services_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withBootstrap()
	    .withOption('rowCallback', rowCallback)
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.DescriptiveName').withTitle('Name'),
		    DTColumnBuilder.newColumn('_source.Record.Status').withTitle('Status'),
		    DTColumnBuilder.newColumn('_source.Record.Mode').withTitle('Mode'),
		    DTColumnBuilder.newColumn('_source.Record.Path').withTitle('Path'),
		    DTColumnBuilder.newColumn('_source.Record.StartedAs').withTitle('Started As'),
		    DTColumnBuilder.newColumn('_source.Record.PathMd5sum').withTitle('PathMD5'),
		    DTColumnBuilder.newColumn('_source.Record.PathSignatureExists').withTitle('PathSignatureExists'),
		    DTColumnBuilder.newColumn('_source.Record.PathSignatureVerified').withTitle('PathSignatureVerified'),
		    DTColumnBuilder.newColumn('_source.Record.PathSignatureDescription').withTitle('PathSignatureDescription'),
		    DTColumnBuilder.newColumn('_source.Record.Pid').withTitle('PID'),
		    DTColumnBuilder.newColumn('_source.Record.Type').withTitle('Type'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodService').withTitle('IsGoodService'),
		    DTColumnBuilder.newColumn('_source.Record.VTResults').withTitle('VTResults'),
		    DTColumnBuilder.newColumn('_source.Record.Arguments').withTitle('Arguments'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    // DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataFileDownload($scope, ngDialog, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, $cookies) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;

	    vm.dtOptions = DTOptionsBuilder.newOptions()
	        .withOption('ajax', {
	         url: nHResponse + 'filedownloadhistory_anchor/' + $routeParams.casename + '/' + $routeParams.hostname,
	         type: 'POST',
	         data: {
	         	data: 'data',
	         	csrfmiddlewaretoken: $cookies.get('csrftoken')
	         }
	     })
        .withOption('processing', true)
        .withOption('serverSide', true)
        .withLanguage({
        	"sProcessing": '<img src="static/images/loader.gif">'
        })
        .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withPaginationType('numbers')
	    .withDisplayLength(200)
	    .withOption('order', [0, 'desc'])
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.LastModifiedDate').withTitle('LastModifiedDate'),
		    DTColumnBuilder.newColumn('_source.Record.SourceUrl').withTitle('SourceUrl'),
		    DTColumnBuilder.newColumn('_source.Record.Filename').withTitle('FileName'),
		    DTColumnBuilder.newColumn('_source.Record.Username').withTitle('Username'),
		    DTColumnBuilder.newColumn('_source.Record.DownloadType').withTitle('Download Type'),
		    DTColumnBuilder.newColumn('_source.Record.Profile').withTitle('Profile'),
		    DTColumnBuilder.newColumn('_source.Record.BrowserName').withTitle('BrowserName'),
		    DTColumnBuilder.newColumn('_source.Record.BrowserVersion').withTitle('BrowserVersion'),
		    DTColumnBuilder.newColumn('_source.Record.TargetDirectory').withTitle('TargetDirectory'),
		    DTColumnBuilder.newColumn('_source.Record.FullHttpHeader').withTitle('FullHttpHeader'),
		    DTColumnBuilder.newColumn('_source.Record.LastCheckedDate').withTitle('LastCheckedDate'),
		    DTColumnBuilder.newColumn('_source.Record.BytesDownloaded').withTitle('BytesDownloaded'),
		    DTColumnBuilder.newColumn('_source.Record.CacheFlags').withTitle('CacheFlags'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodFile').withTitle('IsGoodFile'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    // DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataTasks($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32tasks_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withPaginationType('full_numbers')
	    .withOption('order', [0, 'desc'])
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.NextRunTime').withTitle('NextRunTime'),
		    DTColumnBuilder.newColumn('_source.Record.Name').withTitle('Name'),
		    DTColumnBuilder.newColumn('_source.Record.AccountLogonType').withTitle('AccountLogonType'),
		    DTColumnBuilder.newColumn('_source.Record.AccountName').withTitle('AccountName'),
		    DTColumnBuilder.newColumn('_source.Record.AccountRunLevel').withTitle('AccountRunLevel'),
		    DTColumnBuilder.newColumn('_source.Record.Creator').withTitle('Creator'),
		    DTColumnBuilder.newColumn('_source.Record.CreationDate').withTitle('CreationDate'),
		    DTColumnBuilder.newColumn('_source.Record.ExitCode').withTitle('ExitCode'),
		    DTColumnBuilder.newColumn('_source.Record.MostRecentRunTime').withTitle('MostRecentRunTime'),
		    DTColumnBuilder.newColumn('_source.Record.Status').withTitle('Status'),
		    DTColumnBuilder.newColumn('_source.Record.ActionList').withTitle('ActionList'),
		    DTColumnBuilder.newColumn('_source.Record.TriggerList').withTitle('TriggerList'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),	 
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataNetworkRoute($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32network-route_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Interface').withTitle('Interface'),
		    DTColumnBuilder.newColumn('_source.Record.Destination').withTitle('Destination'),
		    DTColumnBuilder.newColumn('_source.Record.Netmask').withTitle('Netmask'),
		    DTColumnBuilder.newColumn('_source.Record.RouteType').withTitle('RouteType'),
		    DTColumnBuilder.newColumn('_source.Record.Gateway').withTitle('Gateway'),
		    DTColumnBuilder.newColumn('_source.Record.Protocol').withTitle('Protocol'),
		    DTColumnBuilder.newColumn('_source.Record.RouteAge').withTitle('RouteAge'),
		    DTColumnBuilder.newColumn('_source.Record.Metric').withTitle('Metric'),
		    DTColumnBuilder.newColumn('_source.Record.ValidLifeTime').withTitle('ValidLifeTime'),
		    DTColumnBuilder.newColumn('_source.Record.PreferredLifeTime').withTitle('PreferredLifeTime'),
		    DTColumnBuilder.newColumn('_source.Record.IsLoopback').withTitle('IsLoopback'),
		    DTColumnBuilder.newColumn('_source.Record.IsAutoconfiguredAddress').withTitle('IsAutoconfiguredAddress'),
		    DTColumnBuilder.newColumn('_source.Record.IsPublished').withTitle('IsPublished'),
		    DTColumnBuilder.newColumn('_source.Record.IsImortal').withTitle('IsImortal'),
		    DTColumnBuilder.newColumn('_source.Record.Origin').withTitle('Origin'),
		    // DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataUrlHistory($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog, $cookies) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;

	    vm.dtOptions = DTOptionsBuilder.newOptions()
	        .withOption('ajax', {
	         url: nHResponse + 'urlhistory_anchor/' + $routeParams.casename + '/' + $routeParams.hostname,
	         type: 'POST',
	         data: {
	         	data: 'data',
	         	csrfmiddlewaretoken: $cookies.get('csrftoken')
	         }
	     })
        .withOption('processing', true)
        .withOption('serverSide', true)        
        .withLanguage({
        	"sProcessing": '<img src="static/images/loader.gif">'
        })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'desc'])
	    .withPaginationType('numbers')
	    .withDisplayLength(200)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.LastVisitDate').withTitle('LastVisitDate'),
		    DTColumnBuilder.newColumn('_source.Record.Hidden').withTitle('Hidden'),
		    DTColumnBuilder.newColumn('_source.Record.Username').withTitle('Username'),
		    DTColumnBuilder.newColumn('_source.Record.Url').withTitle('Url'),
		    DTColumnBuilder.newColumn('_source.Record.Profile').withTitle('Profile'),
		    DTColumnBuilder.newColumn('_source.Record.BrowserName').withTitle('BrowserName'),
		    DTColumnBuilder.newColumn('_source.Record.BrowserVersion').withTitle('BrowserVersion'),
		    DTColumnBuilder.newColumn('_source.Record.PageTitle').withTitle('PageTitle'),
		    DTColumnBuilder.newColumn('_source.Record.VisitFrom').withTitle('VisitFrom'),
		    DTColumnBuilder.newColumn('_source.Record.VisitType').withTitle('VisitType'),
		    DTColumnBuilder.newColumn('_source.Record.VisitType').withTitle('VisitType'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodDomain').withTitle('IsGoodDomain'),
		    DTColumnBuilder.newColumn('_source.Record.IsNewDomain').withTitle('IsNewDomain'),
		    // DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataNetworkArp($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32network-arp_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'asc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Interface').withTitle('Interface'),
		    DTColumnBuilder.newColumn('_source.Record.PhysicalAddress').withTitle('PhysicalAddress'),
		    DTColumnBuilder.newColumn('_source.Record.Ipv4Address').withTitle('Ipv4Address'),
		    DTColumnBuilder.newColumn('_source.Record.CacheType').withTitle('CacheType'),
		    DTColumnBuilder.newColumn('_source.Record.Ipv6Address').withTitle('Ipv6Address'),
		    DTColumnBuilder.newColumn('_source.Record.InterfaceType').withTitle('InterfaceType'),
		    DTColumnBuilder.newColumn('_source.Record.State').withTitle('State'),
		    DTColumnBuilder.newColumn('_source.Record.IsRouter').withTitle('IsRouter'),
		    DTColumnBuilder.newColumn('_source.Record.LastReachable').withTitle('LastReachable'),
		    DTColumnBuilder.newColumn('_source.Record.LastUnreachable').withTitle('LastUnreachable'),
		    // DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),	 
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataPorts($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32ports_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'asc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Pid').withTitle('Pid'),
		    DTColumnBuilder.newColumn('_source.Record.Process').withTitle('Process'),
		    DTColumnBuilder.newColumn('_source.Record.Path').withTitle('Path'),
		    DTColumnBuilder.newColumn('_source.Record.State').withTitle('State'),
		    DTColumnBuilder.newColumn('_source.Record.LocalIp').withTitle('LocalIp'),
		    DTColumnBuilder.newColumn('_source.Record.RemoteIp').withTitle('RemoteIp'),
		    DTColumnBuilder.newColumn('_source.Record.LocalPort').withTitle('LocalPort'),
		    DTColumnBuilder.newColumn('_source.Record.RemotePort').withTitle('RemotePort'),
		    DTColumnBuilder.newColumn('_source.Record.Protocol').withTitle('Protocol'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodPort').withTitle('IsGoodPort'),
		    // DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataPrefetch($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32prefetch_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'desc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Created').withTitle('Created'),
		    DTColumnBuilder.newColumn('_source.Record.Path').withTitle('FullPath'),
		    DTColumnBuilder.newColumn('_source.Record.SizeInBytes').withTitle('SizeInBytes'),
		    DTColumnBuilder.newColumn('_source.Record.ReportedSizeInBytes').withTitle('ReportedSizeInBytes'),
		    DTColumnBuilder.newColumn('_source.Record.ApplicationFileName').withTitle('ApplicationFileName'),
		    DTColumnBuilder.newColumn('_source.Record.LastRun').withTitle('LastRun'),
		    DTColumnBuilder.newColumn('_source.Record.AccessedFileList').withTitle('AccessedFileList'),
		    DTColumnBuilder.newColumn('_source.Record.ApplicationFullPath').withTitle('ApplicationFullPath'),
		    DTColumnBuilder.newColumn('_source.Record.VolumeDevicePath').withTitle('VolumeDevicePath'),
			DTColumnBuilder.newColumn('_source.Record.VolumeCreationTime').withTitle('VolumeCreationTime'),
			DTColumnBuilder.newColumn('_source.Record.VolumeDevicePath').withTitle('VolumeDevicePath'),
			DTColumnBuilder.newColumn('_source.Record.VolumeSerialNumber').withTitle('VolumeSerialNumber'),
			DTColumnBuilder.newColumn('_source.Record.IsGoodPrefetch').withTitle('IsGoodPrefetch'),		    		    		    		    		    
		    // DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataUserAccounts($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32useraccounts_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'desc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Username').withTitle('Username'),
		    DTColumnBuilder.newColumn('_source.Record.Fullname').withTitle('Fullname'),
		    DTColumnBuilder.newColumn('_source.Record.SecurityType').withTitle('SecurityType'),
		    DTColumnBuilder.newColumn('_source.Record.SecurityID').withTitle('SecurityID'),
		    DTColumnBuilder.newColumn('_source.Record.Description').withTitle('Description'),
		    DTColumnBuilder.newColumn('_source.Record.LastLogin').withTitle('LastLogin'),
		    DTColumnBuilder.newColumn('_source.Record.Disabled').withTitle('Disabled'),
		    DTColumnBuilder.newColumn('_source.Record.HomeDirectory').withTitle('HomeDirectory'),
		    DTColumnBuilder.newColumn('_source.Record.ScriptPath').withTitle('ScriptPath'),
			DTColumnBuilder.newColumn('_source.Record.LockedOut').withTitle('LockedOut'),
			DTColumnBuilder.newColumn('_source.Record.PasswordRequired').withTitle('PasswordRequired'),
			DTColumnBuilder.newColumn('_source.Record.UserPasswordAge').withTitle('UserPasswordAge'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataDNS($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32network-dns_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'desc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Host').withTitle('Host'),
		    DTColumnBuilder.newColumn('_source.Record.RecordName').withTitle('RecordName'),
		    DTColumnBuilder.newColumn('_source.Record.RecordType').withTitle('RecordType'),
		    DTColumnBuilder.newColumn('_source.Record.TimeToLive').withTitle('TimeToLive'),
		    DTColumnBuilder.newColumn('_source.Record.Flags').withTitle('Flags'),
		    DTColumnBuilder.newColumn('_source.Record.DataLength').withTitle('DataLength'),
		    DTColumnBuilder.newColumn('_source.Record.RecordDataList.0.Ipv4Address').withTitle('Ipv4Address'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodEntry').withTitle('IsGoodEntry'),
			// DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataVolumes($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32volumes_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'desc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Name').withTitle('Name'),
		    DTColumnBuilder.newColumn('_source.Record.Type').withTitle('Type'),
		    DTColumnBuilder.newColumn('_source.Record.DevicePath').withTitle('DevicePath'),
		    DTColumnBuilder.newColumn('_source.Record.VolumeName').withTitle('VolumeName'),
		    DTColumnBuilder.newColumn('_source.Record.FileSystemName').withTitle('FileSystemName'),
		    DTColumnBuilder.newColumn('_source.Record.ActualAvailableAllocationUnits').withTitle('ActualAvailableAllocationUnits'),
		    DTColumnBuilder.newColumn('_source.Record.TotalAllocationUnits').withTitle('TotalAllocationUnits'),
		    DTColumnBuilder.newColumn('_source.Record.CreationTime').withTitle('CreationTime'),
			// DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.IsMounted').withTitle('Mounted'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataDisks($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32disks_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'desc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.Host').withTitle('Host'),
		    DTColumnBuilder.newColumn('_source.Record.RecordName').withTitle('RecordName'),
		    DTColumnBuilder.newColumn('_source.Record.RecordType').withTitle('RecordType'),
		    DTColumnBuilder.newColumn('_source.Record.TimeToLive').withTitle('TimeToLive'),
		    DTColumnBuilder.newColumn('_source.Record.Flags').withTitle('Flags'),
		    DTColumnBuilder.newColumn('_source.Record.DataLength').withTitle('DataLength'),
		    DTColumnBuilder.newColumn('_source.Record.RecordDataList.0.Ipv4Address').withTitle('Ipv4Address'),
		    DTColumnBuilder.newColumn('_source.Record.IsGoodEntry').withTitle('IsGoodEntry'),
			// DTColumnBuilder.newColumn('_source.Record.NHScore').withTitle('nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetDataPersistence($scope, DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "w32scripting-persistence_anchor/" + $routeParams.casename + '/' + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withOption('rowCallback', rowCallback)
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withOption('order', [0, 'desc'])
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.Record.PersistenceType').withTitle('Type'),
		    DTColumnBuilder.newColumn('_source.Record.Path').withTitle('FilePath'),
		    DTColumnBuilder.newColumn('_source.Record.FileOwner').withTitle('FileOwner'),
		    DTColumnBuilder.newColumn('_source.Record.FileCreated').withTitle('FileCreated'),
		    DTColumnBuilder.newColumn('_source.Record.FileModified').withTitle('FileModified'),
		    DTColumnBuilder.newColumn('_source.Record.FileAccessed').withTitle('FileAccessed'),
		    DTColumnBuilder.newColumn('_source.Record.FileChanged').withTitle('FileChanged'),
		    DTColumnBuilder.newColumn('_source.Record.RegOwner').withTitle('RegOwner'),
		    DTColumnBuilder.newColumn('_source.Record.RegModified').withTitle('RegModified'),
		    DTColumnBuilder.newColumn('_source.Record.RegPath').withTitle('RegPath'),
		    DTColumnBuilder.newColumn('_source.Record.Md5sum').withTitle('Md5sum'),
		    DTColumnBuilder.newColumn('_source.Record.File.DevicePath').withTitle('File.DevicePath'),
		    DTColumnBuilder.newColumn('_source.Record.File.Path').withTitle('File.FullPath'),
		    DTColumnBuilder.newColumn('_source.Record.File.FilePath').withTitle('File.FilePath'),
		    DTColumnBuilder.newColumn('_source.Record.File.FileName').withTitle('File.FileName'),
		    DTColumnBuilder.newColumn('_source.Record.File.FileExtension').withTitle('File.FileExtension'),
		    DTColumnBuilder.newColumn('_source.Record.File.SizeInBytes').withTitle('File.SizeInBytes'),
		    DTColumnBuilder.newColumn('_source.Record.File.Created').withTitle('File.Created'),
		    DTColumnBuilder.newColumn('_source.Record.File.Modified').withTitle('File.Modified'),
		    DTColumnBuilder.newColumn('_source.Record.File.Accessed').withTitle('File.Accessed'),
		    DTColumnBuilder.newColumn('_source.Record.File.Changed').withTitle('File.Changed'),
		    DTColumnBuilder.newColumn('_source.Record.File.FileAttributes').withTitle('File.FileAttributes'),
		    DTColumnBuilder.newColumn('_source.Record.File.Username').withTitle('File.Username'),
		    DTColumnBuilder.newColumn('_source.Record.File.Md5sum').withTitle('File.Md5sum'),
			DTColumnBuilder.newColumn('_source.Record.File.PeInfo.Type').withTitle('File.PeInfo.Type'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.Subsystem').withTitle('File.PeInfo.Subsystem'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.BaseAddress').withTitle('File.PeInfo.BaseAddress'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.PETimeStamp').withTitle('File.PeInfo.PETimeStamp'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.DigitalSignature.SignatureExists').withTitle('File.DigitalSignature.SignatureExists'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.DigitalSignature.SignatureVerified').withTitle('File.DigitalSignature.SignatureVerified'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.DigitalSignature.Description').withTitle('File.DigitalSignature.Description'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.DigitalSignature.CertificateSubject').withTitle('File.DigitalSignature.CertificateSubject'),
		    DTColumnBuilder.newColumn('_source.Record.File.PeInfo.DigitalSignature.CertificateIssuer').withTitle('File.DigitalSignature.CertificateIssuer'),
			DTColumnBuilder.newColumn('_source.Record.File.IsGoodHash').withTitle('File.IsGoodHash'),
			// DTColumnBuilder.newColumn('_source.Record.File.NHScore').withTitle('File.nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.File.Tag').withTitle('File.Tag'),
		    DTColumnBuilder.newColumn('_source.Record.File.VTResults').withTitle('File.VTResults'),
		    DTColumnBuilder.newColumn('_source.Record.File.Comment.Comment').withTitle('File.Comment'),
			DTColumnBuilder.newColumn('_source.Record.Registry.KeyPath').withTitle('Registry.KeyPath'),
			DTColumnBuilder.newColumn('_source.Record.Registry.Type').withTitle('Registry.Type'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.Modified').withTitle('Registry.Modified'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.ValueName').withTitle('Registry.ValueName'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.Username').withTitle('Username'),
			DTColumnBuilder.newColumn('_source.Record.Registry.Path').withTitle('Registry.Path'),
			DTColumnBuilder.newColumn('_source.Record.Registry.Text').withTitle('Registry.Text'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.ReportedLengthInBytes').withTitle('Registry.ReportedLengthInBytes'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.ValueName').withTitle('Registry.ValueName'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.Username').withTitle('Username'),
			DTColumnBuilder.newColumn('_source.Record.Registry.IsKnownKey').withTitle('Registry.IsKnownKey'),
			// DTColumnBuilder.newColumn('_source.Record.Registry.NHScore').withTitle('Registry.nHScore'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.Tag').withTitle('Registry.Tag'),
		    DTColumnBuilder.newColumn('_source.Record.Registry.Comment.Comment').withTitle('Registry.Comment'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Comment Date'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Comment Analyst'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_id').withTitle('doc_id').notVisible(),
		    DTColumnBuilder.newColumn('_parent').withTitle('parent').notVisible(),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = info;

			ngDialog.openConfirm({
		            template: 'update_doc/',
		            controller: 'confirmController',
		            className: 'ngdialog-theme-default custom-width-800',
		            data: vm.message
				})
			.then(function (success) {

					$http.post(nHResponse + 'update_doc/', JSON.stringify(success)).then(function (data){
						ngDialog.open({
					            template: '<p><center>Document Updated Successfully.</center></p>',
					            plain: true
							})			
					})
				})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}

	function GetStackResponse(DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $routeParams) {
	    var vm = this;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "stack_response/" + $routeParams.hostname).then(function (data) {
	        	return data.data
	        });
	    })
	    .withBootstrap()
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('Attribute').withTitle('Attribute'),
		    DTColumnBuilder.newColumn('Count').withTitle('Count'),
		    DTColumnBuilder.newColumn('Endpoints').withTitle('Endpoints'),
		    DTColumnBuilder.newColumn('AuditType').withTitle('AuditType'),
		    
	    ];		
	}

	function CommentView(DTOptionsBuilder, DTColumnBuilder, DTDefaultOptions, $http, $scope, ngDialog) {
	    var vm = this;
	    vm.message = '';
	    vm.ClickHandler = ClickHandler;
	    DTDefaultOptions.setLoadingTemplate('<img src="static/images/loader.gif">');
	    vm.dtOptions = DTOptionsBuilder.fromFnPromise(function() {
	        return $http.post(nHResponse + "comments/").then(function (data) {
	        	return data.data
	        });
	    })
	    .withBootstrap()
	    .withOption('rowCallback', rowCallback)
	    .withDOM('<"top"flp>rt<"bottom"><"clear">')
	    // .withButtons([
	    // 	'columnsToggle'])
	    .withColReorder()
	    .withPaginationType('full_numbers')
	    .withDisplayLength(100)
	    .withOption('order', [3, 'desc'])
	    .withOption('scrollX', '100%')
	    .withOption('scrollY', '70vh')
	    .withOption('scrollCollapse', true);


	    vm.dtColumns = [
		    DTColumnBuilder.newColumn('_source.CaseInfo.case_name').withTitle('Case #'),
		    DTColumnBuilder.newColumn('_id').withTitle('Doc ID'),
		    DTColumnBuilder.newColumn('_parent').withTitle('Endpoint'),
		    DTColumnBuilder.newColumn('_source.Record.Comment.Date').withTitle('Date'),
			DTColumnBuilder.newColumn('_source.Record.Comment.Analyst').withTitle('Analyst'),
			DTColumnBuilder.newColumn('_source.Record.Comment.Comment').withTitle('Comment'),
		    DTColumnBuilder.newColumn('_source.AuditType.Generator').withTitle('AuditType'),
		    DTColumnBuilder.newColumn('_source.Record.Tag').withTitle('Tag'),
		    
	    ];		
	    function ClickHandler(info) {
	        vm.message = {
	        	"rowId": info._id,
	        	"parent": info._parent
	        };

			$http.post(nHResponse + 'comments/get_comment_doc/', JSON.stringify(vm.message)).then(function (doc_data){
				ngDialog.open({
			            template: 'comments/get_comment_doc/dialog',
			            data: doc_data.data
					})			
			})
	    }

	    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
	        $('td', nRow).unbind('click');
	        $('td', nRow).bind('click', function() {
	            $scope.$apply(function() {
	                vm.ClickHandler(aData);
	            });
	        });
	        return nRow;
	    }
	}
})();

	