/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { NhDataExchangeService } from './nh-data-exchange.service';

describe('NhDataExchangeService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhDataExchangeService]
    });
  });

  it('should ...', inject([NhDataExchangeService], (service: NhDataExchangeService) => {
    expect(service).toBeTruthy();
  }));
});
