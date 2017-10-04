import {Injectable}      from '@angular/core'
import {BehaviorSubject} from 'rxjs/BehaviorSubject';

@Injectable()
export class NhPageTitleService {
  private _pageTitleSource = new BehaviorSubject<string>(null);
  
  navItem$ = this._pageTitleSource.asObservable();
  
  updateTitle(title: string) {
    this._pageTitleSource.next(title);
  }
}