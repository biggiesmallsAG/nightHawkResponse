/* tslint:disable:no-unused-variable */
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

import { NhAuditOutletComponent } from './nh-audit-outlet.component';

describe('NhAuditOutletComponent', () => {
  let component: NhAuditOutletComponent;
  let fixture: ComponentFixture<NhAuditOutletComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NhAuditOutletComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NhAuditOutletComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
