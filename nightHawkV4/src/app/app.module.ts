import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { HttpModule } from '@angular/http';
import { NhNavModule } from './modules/nh-nav/nh-nav.module';
import { NhRouteModule } from './modules/nh-route/nh-route.module';

import { NhContainerComponent } from './components/core/nh-container/nh-container.component';
import { NhMainViewComponent } from './components/core/nh-container/nh-main-view.component';
import { NhDashComponent } from './components/core/nh-dash/nh-dash.component';
import { NhLoaderComponent } from './components/core/utility/nh-loader/nh-loader.component';

import { NhPageTitleService } from './services/nh-page-title.service'; 
import { NhCoreService } from './services/nh-core.service';
import { NhWebsocketService } from './services/nh-websocket.service';
import { NhDataExchangeService } from './services/nh-data-exchange.service';
import { NhGridHelperService } from './services/nh-grid-helper.service';
import { NhLoadingService } from './services/nh-loading.service';
import { D3Service } from 'd3-ng2-service';
import { NhSharedModule } from './modules/nh-shared/nh-shared.module';

@NgModule({
	imports:      [ BrowserModule,
					NhNavModule,
					NhRouteModule,
					HttpModule,
					BrowserAnimationsModule,
					NhSharedModule
	],
	declarations: [ 
					NhDashComponent,
					NhContainerComponent,
					NhMainViewComponent,
					NhLoaderComponent
	],
	providers:    [ NhCoreService, 
					NhPageTitleService,
					NhWebsocketService,
					NhDataExchangeService,
					NhGridHelperService,
					D3Service,
					NhLoadingService
	],
	bootstrap:    [ NhContainerComponent ]
})

export class AppModule { }

