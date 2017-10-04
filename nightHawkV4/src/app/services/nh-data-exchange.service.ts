import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { AuditHandler } from '../components/core/interfaces/audit.interface'

@Injectable()
export class NhDataExchangeService {
  private _dataObject = new BehaviorSubject<Object>(null);
  
  dO$ = this._dataObject.asObservable();

  moveAuditData(data: AuditHandler) {
    this._dataObject.next(data);
  };
}
