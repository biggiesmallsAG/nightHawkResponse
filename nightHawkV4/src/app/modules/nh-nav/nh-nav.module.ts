import { NgModule } from '@angular/core';

import { NhNavComponent } from 'app/components/navigation/nh-nav/nh-nav.component';
import { NhTreeComponent } from 'app/components/navigation/nh-tree/nh-tree.component';

import { NhRouteModule } from '../nh-route/nh-route.module';

@NgModule({
	imports:      [  	NhRouteModule
	],
	declarations: [  	NhNavComponent,
						NhTreeComponent
	],
	exports: [  		NhNavComponent,
						NhTreeComponent

	]
})

export class NhNavModule { }
