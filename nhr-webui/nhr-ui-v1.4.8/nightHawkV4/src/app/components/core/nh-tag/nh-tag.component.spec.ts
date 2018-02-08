import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhTagComponent } from './nh-tag.component';

describe('NhTagComponent', () => {
  let component: NhTagComponent;
  let fixture: ComponentFixture<NhTagComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhTagComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhTagComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
