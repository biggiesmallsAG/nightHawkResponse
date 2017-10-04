import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UploadRouteModule } from 'app/modules/nh-lazy-route/upload-route/upload-route.module';
import { NhSharedModule } from 'app/modules/nh-shared/nh-shared.module';
import { NhUploadComponent } from 'app/components/core/nh-upload/nh-upload.component';

@NgModule({
	imports: [
	CommonModule,
	UploadRouteModule,
	NhSharedModule
	],
	declarations: [
	NhUploadComponent
	],
	exports: [
	NhUploadComponent
	]
})
export class UploadModule { }
