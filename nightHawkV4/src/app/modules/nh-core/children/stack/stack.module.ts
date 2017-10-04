import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StackRouteModule } from 'app/modules/nh-lazy-route/stack-route/stack-route.module';
import { NhSharedModule } from 'app/modules/nh-shared/nh-shared.module';
import { NhStackComponent } from 'app/components/core/nh-stack/nh-stack.component';


@NgModule({
	imports: [
	CommonModule,
	StackRouteModule,
	NhSharedModule
	],
	declarations: [
	NhStackComponent
	],
	exports: [
	NhStackComponent
	]
})
export class StackModule { }
