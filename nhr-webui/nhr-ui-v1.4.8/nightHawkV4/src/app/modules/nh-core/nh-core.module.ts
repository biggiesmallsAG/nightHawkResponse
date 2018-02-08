import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NhSharedModule } from 'app/modules/nh-shared/nh-shared.module';
import { NhLazyRouteModule } from '../nh-lazy-route/nh-lazy-route.module';

@NgModule({
	imports:      [  	NhLazyRouteModule,
						CommonModule,
						NhSharedModule
	],
	declarations: [],
	exports: []
})

export class NhCoreModule { }
