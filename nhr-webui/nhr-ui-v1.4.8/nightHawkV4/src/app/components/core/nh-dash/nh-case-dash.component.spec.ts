/* tslint:disable:no-unused-variable */
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

import { NhCaseDashComponent } from './nh-case-dash.component';

describe('NhCaseDashComponent', () => {
  let component: NhCaseDashComponent;
  let fixture: ComponentFixture<NhCaseDashComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhCaseDashComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhCaseDashComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
