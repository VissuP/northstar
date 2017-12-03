import {Component, ViewEncapsulation} from "@angular/core";
import {Router} from "@angular/router";
import {FormControl, FormGroup, Validators} from "@angular/forms";

import {AuthGuard} from "../shared/auth/auth-guard.service";
import {LoggingService} from "../shared/services/logging.service";

@Component({
    encapsulation: ViewEncapsulation.None,
    selector: "login",
    templateUrl: "login.component.html"
})
export class LoginComponent {

    private router: Router;
    private authGuard: AuthGuard;
    private log: LoggingService;

    private loginForm: FormGroup;
    private username: FormControl;
    private password: FormControl;

    constructor(router: Router, authGuard: AuthGuard, log: LoggingService) {
        this.router = router;
        this.authGuard = authGuard;
        this.log = log;

        this.username = new FormControl("", Validators.compose([Validators.required]));
        this.password = new FormControl("", Validators.compose([Validators.required]));

        this.loginForm = new FormGroup({
            password: this.password,
            username: this.username,
        });
    }

    // Used to authenticate the user.
    public login(event: Event) {
        event.preventDefault();

        // Call the auth service to authenticate the user.
        this.authGuard.login(this.loginForm.value)
            .subscribe((response) => {
                this.router.navigate([this.authGuard.redirectUrl]);
            }, (err) => {
                this.log.error("Logging error : ", JSON.stringify(err));
            });
    }
}
