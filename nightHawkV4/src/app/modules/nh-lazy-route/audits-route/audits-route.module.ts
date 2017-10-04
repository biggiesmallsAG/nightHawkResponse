import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';

import { NhAuditsComponent } from 'app/components/core/nh-audits/nh-audits.component';
import { NhAuditsContainerComponent } from 'app/components/core/nh-audits/nh-audits-container.component';
import { NhAuditOutletComponent } from 'app/components/core/nh-audits/nh-audit-outlet.component';
import { NhCaseDashComponent } from 'app/components/core/nh-dash/nh-case-dash.component';

const routes: Routes = [

	{path: 'view', component: NhAuditsComponent, children: [
		{path: 'casedash', component: NhCaseDashComponent, outlet: 'caseoutlet'},
		{path: 'auditdash', component: NhAuditOutletComponent, outlet: 'auditoutlet' }
	]}
];

@NgModule({
	imports: [ RouterModule.forChild(routes)],
	exports: [ RouterModule ]
})

export class AuditsRouteModule { }