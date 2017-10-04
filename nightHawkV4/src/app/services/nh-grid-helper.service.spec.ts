import { TestBed, inject } from '@angular/core/testing';

import { NhGridHelperService } from './nh-grid-helper.service';

describe('NhGridHelperService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhGridHelperService]
    });
  });

  it('should be created', inject([NhGridHelperService], (service: NhGridHelperService) => {
    expect(service).toBeTruthy();
  }));
});
