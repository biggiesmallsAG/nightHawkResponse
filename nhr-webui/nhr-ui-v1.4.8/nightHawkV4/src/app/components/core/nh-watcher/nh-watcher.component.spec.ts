import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NhWatcherComponent } from './nh-watcher.component';

describe('NhWatcherComponent', () => {
  let component: NhWatcherComponent;
  let fixture: ComponentFixture<NhWatcherComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhWatcherComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhWatcherComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
