import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhLoaderComponent } from './nh-loader.component';

describe('NhLoaderComponent', () => {
  let component: NhLoaderComponent;
  let fixture: ComponentFixture<NhLoaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhLoaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhLoaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
