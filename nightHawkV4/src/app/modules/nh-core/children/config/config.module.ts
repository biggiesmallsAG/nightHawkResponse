import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConfigRouteModule } from 'app/modules/nh-lazy-route/config-route/config-route.module';
import { NhSharedModule } from 'app/modules/nh-shared/nh-shared.module';

import { NhConfigComponent } from 'app/components/core/nh-config/nh-config.component';
import { NhConfigSysComponent } from 'app/components/core/nh-config/nh-config-sys.component';
import { NhConfigStatsComponent } from 'app/components/core/nh-config/nh-config-stats.component';
import { NhJobComponent } from 'app/components/core/nh-job/nh-job.component';
import { NhWatcherComponent } from 'app/components/core/nh-watcher/nh-watcher.component';

import { NhWatcherrulesService } from 'app/services/nh-watcherrules.service';

@NgModule({
	imports: [
	CommonModule,
	ConfigRouteModule,
	NhSharedModule
	],
	declarations: [
	NhConfigComponent,
	NhConfigSysComponent,
	NhConfigStatsComponent,
	NhJobComponent,
	NhWatcherComponent
	],
	providers: [
	NhWatcherrulesService
	],
	exports: [
	NhConfigComponent,
	NhConfigSysComponent,
	NhConfigStatsComponent,
	NhJobComponent,
	NhWatcherComponent
	]
})
export class ConfigModule { }
