import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { NhSearchComponent } from 'app/components/core/nh-search/nh-search.component';
import { NhTimelineComponent } from 'app/components/core/nh-timeline/nh-timeline.component';
import { NhDiffComponent } from 'app/components/core/nh-diff/nh-diff.component';

const routes: Routes = [
	{ path: 'view/global', component: NhSearchComponent },
	{ path: 'view/timeline', component: NhTimelineComponent },
	{ path: 'view/diff', component: NhDiffComponent }
];

@NgModule({
  imports: [ RouterModule.forChild(routes)],
  exports: [ RouterModule ]
})

export class SearchRouteModule { }
