import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhStackComponent } from './nh-stack.component';

describe('NhStackComponent', () => {
  let component: NhStackComponent;
  let fixture: ComponentFixture<NhStackComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhStackComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhStackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
