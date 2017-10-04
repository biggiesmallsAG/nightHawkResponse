import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [

	{ path: 'audits', loadChildren: "app/modules/nh-core/children/audits/audits.module#AuditsModule" },
	{ path: 'config', loadChildren: "app/modules/nh-core/children/config/config.module#ConfigModule" },
	{ path: 'stack', loadChildren: "app/modules/nh-core/children/stack/stack.module#StackModule"},
	{ path: 'search', loadChildren: "app/modules/nh-core/children/search/search.module#SearchModule"},
	{ path: 'upload', loadChildren: "app/modules/nh-core/children/upload/upload.module#UploadModule" }
];

@NgModule({
  imports: [ RouterModule.forChild(routes)],
  exports: [ RouterModule ]
})

export class NhLazyRouteModule { }
