/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { NhPageTitleService } from './nh-page-title.service';

describe('NhPageTitleService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NhPageTitleService]
    });
  });

  it('should ...', inject([NhPageTitleService], (service: NhPageTitleService) => {
    expect(service).toBeTruthy();
  }));
});
