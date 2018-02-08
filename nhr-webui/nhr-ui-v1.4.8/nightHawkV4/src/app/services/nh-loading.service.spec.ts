import { TestBed, inject } from '@angular/core/testing';

import { NhLoadingService } from './nh-loading.service';

describe('NhLoadingService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhLoadingService]
    });
  });

  it('should be created', inject([NhLoadingService], (service: NhLoadingService) => {
    expect(service).toBeTruthy();
  }));
});
