import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhSearchComponent } from './nh-search.component';

describe('NhSearchComponent', () => {
  let component: NhSearchComponent;
  let fixture: ComponentFixture<NhSearchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhSearchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhSearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
