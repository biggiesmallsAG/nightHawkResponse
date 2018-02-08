import { Directive, Input, OnInit, ElementRef } from '@angular/core';

declare var renderjson: any;

@Directive({
	selector: '[RenderJson]'
})
export class RenderjsonDirective implements OnInit {
	
	@Input() jsonObject: Object;

	private rjson;
	private rawObject: Object;
	private opts: Object;
	
	constructor(private el: ElementRef) {}

	ngOnInit() {
		this.rawObject = this.jsonObject;
		this.rjson = new renderjson(this.rawObject)
		this.el.nativeElement.append(this.rjson)
	}
}
