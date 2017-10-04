import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhFormControlComponent } from './nh-form-control.component';

describe('NhFormControlComponent', () => {
  let component: NhFormControlComponent;
  let fixture: ComponentFixture<NhFormControlComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhFormControlComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhFormControlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
