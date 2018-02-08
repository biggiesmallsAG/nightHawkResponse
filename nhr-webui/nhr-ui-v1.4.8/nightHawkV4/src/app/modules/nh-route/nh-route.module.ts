import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { NhDashComponent } from 'app/components/core/nh-dash/nh-dash.component';
import { NhAuthComponent } from 'app/components/core/nh-auth/nh-auth.component';

const routes: Routes = [
	{ path: '', redirectTo: '/app/(dashoutlet:dashboard)', pathMatch: 'full' },
	{ path: 'app', children: [
		{ path: 'dashboard', component: NhAuthComponent, outlet: 'dashoutlet' }
	] },
	{ path: 'core', loadChildren: 'app/modules/nh-core/nh-core.module#NhCoreModule' }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ],
  declarations: []
})

export class NhRouteModule { }
