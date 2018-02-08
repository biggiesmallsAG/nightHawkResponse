/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { NhWebsocketService } from './nh-websocket.service';

describe('NhWebsocketService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhWebsocketService]
    });
  });

  it('should ...', inject([NhWebsocketService], (service: NhWebsocketService) => {
    expect(service).toBeTruthy();
  }));
});
