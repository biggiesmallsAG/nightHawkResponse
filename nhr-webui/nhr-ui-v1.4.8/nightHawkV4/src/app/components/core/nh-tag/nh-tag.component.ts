import { Component, OnInit, Input } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TagHandler } from 'app/components/core/interfaces/tag.interface';
import { NhCoreService } from 'app/services/nh-core.service';
import { AuditHandler } from 'app/components/core/interfaces/audit.interface';

@Component({
  selector: 'app-nh-tag',
  templateUrl: './nh-tag.component.html',
  styleUrls: ['./nh-tag.component.sass']
})
export class NhTagComponent implements OnInit {
  @Input() tagMeta: AuditHandler;
  @Input() docId: string;
	public tagOptions: Array<Object> = [
    {name: "Benign", opt: "benign"},
    {name: "Follow Up", opt: "follow_up"},
    {name: "Malicious", opt: "malicious"},
    {name: "For Review", opt: "for_review"}
    ];
  public tagForm: FormGroup;
  public tagResponse: string;
  
  constructor(
    private _fb: FormBuilder,
    private _NhCoreSvc: NhCoreService
  ) { }

  ngOnInit() {
    this.tagForm  = this._fb.group({
      comment: [''],
      TagName: [''],
      tag_path: ['/tag'],
      comment_path: ['/comment']
    })
  }

  save(model: TagHandler, isValid: boolean, event: Event) {
    event.preventDefault();
    const user = localStorage.getItem('nhr-user');
    if (model.comment) {
      const commentBody = {
        Comment: model.comment,
        CreatedBy: user
      };
      this._NhCoreSvc.POSTJSON(model.comment_path + `/add/${this.tagMeta.case_name}/${this.tagMeta.endpoint}/${this.tagMeta.audit_type}/${this.docId}`,
      commentBody)
      .toPromise()
      .then(
        (response) => {
          console.log(response)
        }
      )
    };
    if (model.TagName) {
      let tagBody = {
        TagName: model.TagName,
        CreatedBy: user
      }
      this._NhCoreSvc.POSTJSON(model.tag_path + `/add/${this.tagMeta.case_name}/${this.tagMeta.endpoint}/${this.tagMeta.audit_type}/${this.docId}`, tagBody)
      .toPromise()
      .then(
        (response) => {
          this.tagResponse = response;
        }
      )
    }
  }

}
