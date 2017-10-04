import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { MaterializeDirective } from "angular2-materialize";
import { AgGridModule } from 'ag-grid-angular/main';
import { RenderjsonDirective } from 'app/directives/renderjson.directive';

import { NhFormControlComponent } from 'app/components/forms/nh-form-control/nh-form-control.component';

@NgModule({
	imports: [
	CommonModule,
	AgGridModule.withComponents([]),
	FormsModule, 
	ReactiveFormsModule
	],
	declarations: [
	MaterializeDirective,
	RenderjsonDirective,
	NhFormControlComponent
	],
	exports: [
	NhFormControlComponent,
	MaterializeDirective,
	RenderjsonDirective,
	AgGridModule,
	FormsModule, 
	ReactiveFormsModule,
	]
})
export class NhSharedModule { }
