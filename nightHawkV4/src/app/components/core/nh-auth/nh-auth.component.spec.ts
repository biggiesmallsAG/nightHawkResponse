import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhAuthComponent } from './nh-auth.component';

describe('NhAuthComponent', () => {
  let component: NhAuthComponent;
  let fixture: ComponentFixture<NhAuthComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhAuthComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhAuthComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
