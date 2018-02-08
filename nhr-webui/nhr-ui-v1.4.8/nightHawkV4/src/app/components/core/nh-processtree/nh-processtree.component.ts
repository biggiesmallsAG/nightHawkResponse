import { Component, OnInit, AfterViewInit, Input, NgZone, ElementRef, ViewChild, ViewEncapsulation, OnDestroy } from '@angular/core';
import {
	D3Service,
	D3,
	Axis,
	Selection,
	TreeLayout,
	ForceLink,
	HierarchyNode,
	D3ZoomEvent
} from 'd3-ng2-service';

import { AuditHandler } from '../interfaces/audit.interface';

export interface Node {
	id: string,
	group: number
}

export interface Link {
	source: string,
	target: Node,
}

@Component({
	selector: 'app-nh-processtree',
	templateUrl: './nh-processtree.component.html',
	styleUrls: ['./nh-processtree.component.sass'],
	encapsulation: ViewEncapsulation.None
})
export class NhProcesstreeComponent implements AfterViewInit {

	@ViewChild('svgWrapper') elementView: ElementRef;
	@Input() visualizeItem: AuditHandler;

	private d3: D3;
	private parentNativeElement: any;
	private d3Svg: Selection<SVGSVGElement, any, null, undefined>;
	private graphSelected;

	constructor(
		private element: ElementRef,
		private ngZone: NgZone,
		d3service: D3Service) {
		this.d3 = d3service.getD3();
		this.parentNativeElement = element.nativeElement;
	}

	ngAfterViewInit() {
		let self = this;
		let d3 = this.d3;
		let d3ParentElement: Selection<HTMLElement, any, null, undefined>;
		let d3Svg: Selection<SVGSVGElement, any, null, undefined>;
		var duration = 750;
		var i = 0, root;

		setTimeout(() => { // this is a hack because there is some viewchild bug 
			// build SVG tree layout
			if (this.parentNativeElement !== null) {
				d3ParentElement = d3.select(this.parentNativeElement);
				d3Svg = this.d3Svg = d3ParentElement.select<SVGSVGElement>('svg');

				var div = d3.select("body")
				.append("div")
				.attr("class", "tooltip")
				.style("opacity", 0);

				var tree = d3.tree()
				.size([this.elementView.nativeElement.offsetHeight, this.elementView.nativeElement.offsetWidth]);

				var svg = d3.select("body").append("svg")
				.attr("width", this.elementView.nativeElement.offsetWidth + 240)
				.attr("height", this.elementView.nativeElement.offsetHeight + 40)
				.attr("transform", "translate(" + 120 + "," + 20 + ")")
				.append("g");

				// d3Svg.call(d3.zoom<SVGSVGElement, any>()
				// 	.scaleExtent([1 / 2, 4])
				// 	.on('zoom', zoomed));

				root = d3.hierarchy(this.visualizeItem.data[0]._source.Record, function(d) {
					return d._children; 
				});
				root.x0 = (this.elementView.nativeElement.offsetHeight - 40) / 2;
				root.y0 = 0;

				root.children.forEach(collapse)
				update(root);
			};
				// function zoomed(this: SVGSVGElement) {
				// 	let e: D3ZoomEvent<SVGSVGElement, any> = d3.event;
				// 	d3Svg.attr('transform', e.transform.toString());
				// };
				function collapse(d) {
					if (d.children) {
						d._children = d.children;
						d._children.forEach(collapse);
						d.children = null;
					}
				};
				function update(source) {
					// Assigns the x and y position for the nodes
					var treeData = tree(root);

					// Compute the new tree layout.
					var nodes = treeData.descendants(),
					links = treeData.descendants().slice(1);

					// Normalize for fixed-depth.
					nodes.forEach(function(d){ d.y = d.depth * 180});

					// ****************** Nodes section ***************************

					// Update the nodes...
					var node = d3Svg.selectAll('g.node')
					.data(nodes, (d: any) => {return d.id || (d.id = ++i); });

					// Enter any new modes at the parent's previous position.
					var nodeEnter = node.enter().append('g')
					.attr('class', 'node')
					.attr("transform", function(d) {
						return "translate(" + source.y0 + "," + source.x0 + ")";
					})
					.on('click', click) .on("mouseover", (d: any) => {
						div.transition()
						.duration(200)
						.style("opacity", .95);
						div.html(
							"PID: " + '<span class="tooltiphighlight">' + d.data.pid + "</span><br/>" + 
							"Args: " + '<span class="tooltiphighlight">' + d.data.arguments + "</span><br/>" +
							"Start Time: " + '<span class="tooltiphighlight">' + d.data.startTime + "</span><br/>" +
							"Path: " + '<span class="tooltiphighlight">' + d.data.path + "</span><br/>" + 
							"Parent: " + '<span class="tooltiphighlight">' + d.data.parent + "</span><br/>"
							)
						.style("left", (d3.event.pageX) + "px")
						.style("top", (d3.event.pageY - 28) + "px");
					})
					.on("mouseout", function(d) {
						div.transition()
						.duration(500)
						.style("opacity", 0);
					});;

					// Add Circle for the nodes
					nodeEnter.append('circle')
					.attr('class', 'node')
					.attr('r', 1e-6)
					.style("fill", function(d: any) {
						return d._children ? "lightsteelblue" : "#fff";
					});

					// Add labels for the nodes
					nodeEnter.append('text')
					.attr("dy", ".35em")
					.attr("x", function(d: any) {
						return d.children || d._children ? -13 : 13;
					})
					.attr("text-anchor", function(d: any) {
						return d.children || d._children ? "end" : "start";
					})
					.text(function(d: any) { return d.data.name; });

					// UPDATE
					var nodeUpdate = nodeEnter.merge(node);

					// Transition to the proper position for the node
					nodeUpdate.transition()
					.duration(duration)
					.attr("transform", function(d) { 
						return "translate(" + d.y + "," + d.x + ")";
					});

					// Update the node attributes and style
					nodeUpdate.select('circle.node')
					.attr('r', 10)
					.style("fill", function(d: any) {
						return d._children ? "#FF773D" : "#4aaabf";
					})
					.attr('cursor', 'pointer');


					// Remove any exiting nodes
					var nodeExit = node.exit().transition()
					.duration(duration)
					.attr("transform", function(d) {
						return "translate(" + source.y + "," + source.x + ")";
					})
					.remove();

					// On exit reduce the node circles size to 0
					nodeExit.select('circle')
					.attr('r', 1e-6);

					// On exit reduce the opacity of text labels
					nodeExit.select('text')
					.style('fill-opacity', 1e-6);

					// ****************** links section ***************************

					// Update the links...
					var link = d3Svg.selectAll('path.link')
					.data(links, function(d: any) { return d.id; });

					// Enter any new links at the parent's previous position.
					var linkEnter = link.enter().insert('path', "g")
					.attr("class", "link")
					.attr('d', function(d){
						var o = {x: source.x0, y: source.y0}
						return diagonal(o, o)
					});

					// UPDATE
					var linkUpdate = linkEnter.merge(link);

					// Transition back to the parent element position
					linkUpdate.transition()
					.duration(duration)
					.attr('d', function(d){ return diagonal(d, d.parent) });

					// Remove any exiting links
					var linkExit = link.exit().transition()
					.duration(duration)
					.attr('d', function(d) {
						var o = {x: source.x, y: source.y}
						return diagonal(o, o)
					})
					.remove();

					// Store the old positions for transition.
					nodes.forEach(function(d: any){
						d.x0 = d.x;
						d.y0 = d.y;
					});

					// Creates a curved (diagonal) path from parent to the child nodes
					function diagonal(s, d) {

						var path = `M ${s.y} ${s.x}
						C ${(s.y + d.y) / 2} ${s.x},
						${(s.y + d.y) / 2} ${d.x},
						${d.y} ${d.x}`

						return path
					}

					// Toggle children on click.
					function click(d) {
						if (d.children) {
							d._children = d.children;
							d.children = null;
						} else {
							d.children = d._children;
							d._children = null;
						}
						update(d);
					}
				}			
			}, 100)		
		}

		ngOnInit() {
		}

		private selectItem = (param) => {
			this.graphSelected = param.id
		}

		ngOnDestroy() {
			if (this.d3Svg.empty && !this.d3Svg.empty()) {
				this.d3Svg.selectAll('*').remove();
			}		
		}
	}
