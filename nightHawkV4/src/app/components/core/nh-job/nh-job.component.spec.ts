/* tslint:disable:no-unused-variable */
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

import { NhJobComponent } from './nh-job.component';

describe('NhJobComponent', () => {
  let component: NhJobComponent;
  let fixture: ComponentFixture<NhJobComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhJobComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhJobComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
