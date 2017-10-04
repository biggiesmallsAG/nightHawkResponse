import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuditsRouteModule } from 'app/modules/nh-lazy-route/audits-route/audits-route.module';
import { NhSharedModule } from 'app/modules/nh-shared/nh-shared.module';

import { NhAuditsComponent } from 'app/components/core/nh-audits/nh-audits.component';
import { NhAuditsContainerComponent } from 'app/components/core/nh-audits/nh-audits-container.component';
import { NhAuditOutletComponent } from 'app/components/core/nh-audits/nh-audit-outlet.component';
import { NhCaseDashComponent } from 'app/components/core/nh-dash/nh-case-dash.component';
import { NhProcesstreeComponent } from 'app/components/core/nh-processtree/nh-processtree.component';

@NgModule({
	imports: [
	CommonModule,
	AuditsRouteModule,
	NhSharedModule
	],
	declarations: [
	NhAuditsComponent,
	NhAuditsContainerComponent,
	NhAuditOutletComponent,
	NhCaseDashComponent,
	NhProcesstreeComponent
	],
	exports: [
	NhAuditsComponent,
	NhAuditsContainerComponent,
	NhAuditOutletComponent,
	NhCaseDashComponent,
	NhProcesstreeComponent
	]
})
export class AuditsModule { }
