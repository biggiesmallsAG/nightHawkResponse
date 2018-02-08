import { Component, OnInit, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';
import { FormGroup } from '@angular/forms';
import { Validators } from '@angular/forms';
import { UserModel } from 'app/components/core/interfaces/user.interface';
import { NhCoreService } from 'app/services/nh-core.service';
import { MaterializeAction } from 'angular2-materialize';
import { NhLoadingService } from 'app/services/nh-loading.service';

@Component({
  selector: 'app-nh-auth',
  templateUrl: './nh-auth.component.html',
  styleUrls: ['./nh-auth.component.sass']
})
export class NhAuthComponent implements OnInit {
  public loginForm: FormGroup;
  public loginResponse: boolean;
  public loginError: string;
  modalActions = new EventEmitter<string|MaterializeAction>();
  constructor(
    private _fb: FormBuilder,
    private _NhCoreSvc: NhCoreService,
    private _NhLoaderSvc: NhLoadingService,
    private router: Router
  ) { }

  ngOnInit() {
    this._NhLoaderSvc.hide();
    this.loginForm = this._fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required],
      path: ['/auth/login']
    })
  }

  loginUser(model: UserModel, isValid: boolean, event: Event) {
    event.preventDefault();
    this._NhCoreSvc.POSTJSON(model.path, model)
      .toPromise()
      .then(
      ((response) => {
        try {
          this.loginResponse = response.data.token;
          localStorage.setItem('nhr-token', response.data.token);
          localStorage.setItem('nhr-user', response.data.username);
          setTimeout(() => {
            this.modalActions.emit({action:"modal",params:['close']});
            location.reload()
          }, 5000)
        } catch (error) {
          this.loginResponse = null;
          this.loginError = response.reason
        }
        this.modalActions.emit({action:"modal",params:['open']});
      })
    )
  }
}
