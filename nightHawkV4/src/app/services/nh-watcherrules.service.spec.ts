import { TestBed, inject } from '@angular/core/testing';

import { NhWatcherrulesService } from './nh-watcherrules.service';

describe('NhWatcherrulesService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhWatcherrulesService]
    });
  });

  it('should be created', inject([NhWatcherrulesService], (service: NhWatcherrulesService) => {
    expect(service).toBeTruthy();
  }));
});
