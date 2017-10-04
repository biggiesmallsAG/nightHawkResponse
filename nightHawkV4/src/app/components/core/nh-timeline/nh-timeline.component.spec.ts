import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhTimelineComponent } from './nh-timeline.component';

describe('NhTimelineComponent', () => {
  let component: NhTimelineComponent;
  let fixture: ComponentFixture<NhTimelineComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhTimelineComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhTimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
