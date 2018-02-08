import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { NhUploadComponent } from 'app/components/core/nh-upload/nh-upload.component';

const routes: Routes = [
	{ path: 'view', component: NhUploadComponent },
];

@NgModule({
  imports: [ RouterModule.forChild(routes)],
  exports: [ RouterModule ]
})

export class UploadRouteModule { }
