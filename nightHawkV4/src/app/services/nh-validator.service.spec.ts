import { TestBed, inject } from '@angular/core/testing';

import { NhValidatorService } from './nh-validator.service';

describe('NhValidatorService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhValidatorService]
    });
  });

  it('should be created', inject([NhValidatorService], (service: NhValidatorService) => {
    expect(service).toBeTruthy();
  }));
});
