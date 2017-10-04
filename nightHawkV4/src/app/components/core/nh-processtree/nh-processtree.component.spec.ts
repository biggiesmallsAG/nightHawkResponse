import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhProcesstreeComponent } from './nh-processtree.component';

describe('NhProcesstreeComponent', () => {
  let component: NhProcesstreeComponent;
  let fixture: ComponentFixture<NhProcesstreeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhProcesstreeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhProcesstreeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
