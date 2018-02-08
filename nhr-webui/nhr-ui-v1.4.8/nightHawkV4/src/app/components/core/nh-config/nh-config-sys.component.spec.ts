/* tslint:disable:no-unused-variable */
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

import { NhConfigSysComponent } from './nh-config-sys.component';

describe('NhConfigSysComponent', () => {
  let component: NhConfigSysComponent;
  let fixture: ComponentFixture<NhConfigSysComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhConfigSysComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhConfigSysComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
