import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchRouteModule } from 'app/modules/nh-lazy-route/search-route/search-route.module';
import { NhSharedModule } from 'app/modules/nh-shared/nh-shared.module';

import { NhSearchComponent } from 'app/components/core/nh-search/nh-search.component';
import { NhTimelineComponent } from 'app/components/core/nh-timeline/nh-timeline.component';
import { NhDiffComponent } from 'app/components/core/nh-diff/nh-diff.component';

@NgModule({
	imports: [
	CommonModule,
	SearchRouteModule,
	NhSharedModule
	],
	declarations: [
	NhSearchComponent,
	NhTimelineComponent,
	NhDiffComponent
	],
	exports: [
	NhSearchComponent,
	NhTimelineComponent,
	NhDiffComponent
	]
})
export class SearchModule { }
