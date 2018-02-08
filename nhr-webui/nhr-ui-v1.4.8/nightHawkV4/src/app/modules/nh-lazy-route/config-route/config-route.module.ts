import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { NhConfigComponent } from 'app/components/core/nh-config/nh-config.component';
import { NhConfigSysComponent } from 'app/components/core/nh-config/nh-config-sys.component';
import { NhConfigStatsComponent } from 'app/components/core/nh-config/nh-config-stats.component';
import { NhJobComponent } from 'app/components/core/nh-job/nh-job.component';
import { NhWatcherComponent } from 'app/components/core/nh-watcher/nh-watcher.component';

const routes: Routes = [
	{ path: 'view', component: NhConfigComponent },
	{ path: 'view/sysconfig', component: NhConfigSysComponent },
	{ path: 'view/platformstats', component: NhConfigStatsComponent }, 
	{ path: 'view/job', component: NhJobComponent },
	{ path: 'view/watcher', component: NhWatcherComponent }
];

@NgModule({
  imports: [ RouterModule.forChild(routes)],
  exports: [ RouterModule ]
})

export class ConfigRouteModule { }
