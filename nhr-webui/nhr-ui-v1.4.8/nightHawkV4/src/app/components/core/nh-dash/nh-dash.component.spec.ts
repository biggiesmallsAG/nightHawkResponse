/* tslint:disable:no-unused-variable */
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

import { NhDashComponent } from './nh-dash.component';

describe('NhDashComponent', () => {
  let component: NhDashComponent;
  let fixture: ComponentFixture<NhDashComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhDashComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhDashComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
