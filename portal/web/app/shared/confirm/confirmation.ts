import {Modal} from 'angular2-modal/plugins/bootstrap';

//confirm method used to confirm before any action taken.
export function confirm(title: string, message: string, modal: Modal) {

    return modal.confirm()
        .title(title)
        .body(message)
        .size('sm')
        .okBtn('Ok')
        .cancelBtn('Cancel')
        .open();
}
