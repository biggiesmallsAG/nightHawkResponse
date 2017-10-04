import { Injectable } from '@angular/core';

@Injectable()
export class NhGridHelperService {

	public columnDefs: Array<Object> = [];

	constructor() { }

	public iterate(obj: Object, type?: string) {
		this.columnDefs = [];
		var walked = [];
		var stack = [{obj: obj, stack: ''}];
		while(stack.length > 0)
		{
			var item = stack.pop();
			var obj = item.obj;
			for (var property in obj) {
				if (obj.hasOwnProperty(property)) {
					if (typeof obj[property] == "object") {
						var alreadyFound = false;
						for(var i = 0; i < walked.length; i++)
						{
							if (walked[i] === obj[property])
							{
								alreadyFound = true;
								break;
							}
						}
						if (!alreadyFound)
						{
							walked.push(obj[property]);
							stack.push({obj: obj[property], stack: item.stack + '.' + property});
						}
					}
					else
					{
						var exp = `${item.stack}.${property}`;
						switch (exp) {
							case "._index":
							case "._type":
							case "._uid":
							case "._routing":
							case "._parent":
							case ".sort.0":
							break;
							case "._id":
							this.columnDefs.push({
								headerName: "_id",
								field: "_id",
								hide: true
							});
							break;						
							default:
							if (type) {
								var regex = new RegExp("\.(.*)");
								var header = regex.exec(exp);
								this.columnDefs.push({
									headerName: header[1],
									field: header[1]
								});
								break;
							} else {
								var regex = new RegExp("\._source\.(.*)");
								var header = regex.exec(exp);
								var prefetch_strip = new RegExp("^Record.AccessedFileList\.\\d{1,2}")
								if (prefetch_strip.exec(header[1])) {
									// pass through and dont build a col
								} else {
									switch (header[1]) {
										case "ComputerName":
										case "Record.JobCreated":
										case "CaseInfo.case_name":
										case "CaseInfo.case_date":
										case "CaseInfo.computer_name":
										this.columnDefs.push({
											headerName: header[1],
											field: `_source.${header[1]}`,
											hide: true
										})			
										break;
										case "Record.ReportedSizeInBytes":
										case "Record.TimesExecuted":
										case "Record.SizeInBytes":
										case "Record.BytesDownloaded":																				case "Record.CacheHitCount":
										case "Record.CacheHitCount":
										case "Record.MaxBytes":
										case "Record.VisitCount":
										case "Record.ReportedLengthInBytes":
										case "Record.Registry.ReportedLengthInBytes":
										case "Record.File.SizeInBytes":
										case "Record.DataLength":
										case "Record.Pid":
										case "Record.Index":
										case "Record.EID":
										case "Record.LocalPort":
										case "Record.RemotePort":
										this.columnDefs.push({
											headerName: header[1],
											field: `_source.${header[1]}`,
											filter: 'number'
										});
										break;
										case "Record.TlnTime":
										this.columnDefs.push({
											headerName: header[1],
											field: `_source.${header[1]}`,
											sort: "desc"
										});
										break;
										default:
										this.columnDefs.push({
											headerName: header[1],
											field: `_source.${header[1]}`,
											filter: 'text'
										})
										break;
									}
								}
								break;
							}
						}
					}
				}
			};
		};
		return this.columnDefs
	};

}
