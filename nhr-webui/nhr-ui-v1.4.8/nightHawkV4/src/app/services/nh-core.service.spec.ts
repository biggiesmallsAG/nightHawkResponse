/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { NhCoreService } from './nh-core.service';

describe('NhCoreService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhCoreService]
    });
  });

  it('should ...', inject([NhCoreService], (service: NhCoreService) => {
    expect(service).toBeTruthy();
  }));
});
