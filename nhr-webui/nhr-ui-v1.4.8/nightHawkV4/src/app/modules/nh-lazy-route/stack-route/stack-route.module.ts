import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { NhStackComponent } from 'app/components/core/nh-stack/nh-stack.component';

const routes: Routes = [
	{ path: 'view', component: NhStackComponent },
];

@NgModule({
  imports: [ RouterModule.forChild(routes)],
  exports: [ RouterModule ]
})

export class StackRouteModule { }
