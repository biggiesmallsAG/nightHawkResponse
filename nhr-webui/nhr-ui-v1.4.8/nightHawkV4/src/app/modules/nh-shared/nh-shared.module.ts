import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { MaterializeModule } from "angular2-materialize";
import { AgGridModule } from 'ag-grid-angular/main';
import { RenderjsonDirective } from 'app/directives/renderjson.directive';

import { NhFormControlComponent } from 'app/components/forms/nh-form-control/nh-form-control.component';
import { NhTagComponent } from 'app/components/core/nh-tag/nh-tag.component';
import { HistoryComponent } from 'app/components/core/nh-tag/history/history.component';
import { CustomDateComponent } from 'app/components/core/nh-audits/custom-renderers/customdate.component';

@NgModule({
	imports: [
	CommonModule,
	FormsModule, 
	ReactiveFormsModule,
	MaterializeModule,
	AgGridModule.withComponents([CustomDateComponent]),
	],
	declarations: [
	RenderjsonDirective,
	NhFormControlComponent,
	NhTagComponent,
	HistoryComponent,
	CustomDateComponent
	],
	exports: [
	NhFormControlComponent,
	MaterializeModule,
	RenderjsonDirective,
	AgGridModule,
	FormsModule, 
	ReactiveFormsModule,
	NhTagComponent,
	HistoryComponent
	]
})
export class NhSharedModule { }
