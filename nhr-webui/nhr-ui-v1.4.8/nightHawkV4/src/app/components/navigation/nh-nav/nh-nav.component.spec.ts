/* tslint:disable:no-unused-variable */
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

import { NhNavComponent } from './nh-nav.component';

describe('NhNavComponent', () => {
  let component: NhNavComponent;
  let fixture: ComponentFixture<NhNavComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhNavComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhNavComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
