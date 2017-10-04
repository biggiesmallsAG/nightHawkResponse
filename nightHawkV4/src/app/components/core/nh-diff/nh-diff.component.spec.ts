import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhDiffComponent } from './nh-diff.component';

describe('NhDiffComponent', () => {
  let component: NhDiffComponent;
  let fixture: ComponentFixture<NhDiffComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhDiffComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhDiffComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
