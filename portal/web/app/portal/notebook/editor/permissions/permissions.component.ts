import {Component, Input, OnDestroy} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {DialogRef, ModalComponent} from "angular2-modal";
import {BSModalContext} from "angular2-modal/plugins/bootstrap";
import {Router} from "@angular/router";
import {LoggingService} from "../../../../shared/services/logging.service";
import {RemotePostData, getRemotePostData} from "./users-data.service";
import {Http} from "@angular/http";
import {NotebookService} from "../../../../shared/services/notebook.service";
import {NotebookPermissions, PermissionsArray} from "../../../../shared/models/permissions.model";
import {User} from "../../../../shared/models/user.model";
import {Alert, AlertType} from "../../../../shared/components/alerts/alerts.model";
import {AlertService} from "../../../../shared/services/alerts.service";
import {NsError} from "../../../../shared/models/error.model";
import {Subscription} from "rxjs";

export class PermissionsModalContext extends BSModalContext {
    public id: string;
}

@Component({
    styleUrls: ["permissions.component.css"],
    templateUrl: "permissions.component.html",
})

export class PermissionsModal implements ModalComponent<PermissionsModalContext>, OnDestroy {
    @Input() context: PermissionsModalContext;
    public dialog: DialogRef<PermissionsModalContext>;

    private dataService: RemotePostData;
    private log: LoggingService;
    private users: User[];
    private Permissions: string[];
    private newUser: User;
    private notebookService: NotebookService;
    private usersChanged: boolean;
    private alerts: AlertService;
    private alert: Alert;
    private navigateSubscription: Subscription;
    private router: Router;

    constructor(dialog: DialogRef<PermissionsModalContext>, logger: LoggingService, http: Http,
                notebookService: NotebookService, alerts: AlertService, router: Router) {
        this.log = logger;
        this.alerts = alerts;
        this.router = router;
        this.newUser = new User();
        this.usersChanged = false;

        // Set up the dialog and grab our arguments passed in via the context
        this.dialog = dialog;
        this.context = dialog.context;
        dialog.setCloseGuard(this);

        // this should not happen, but if it does, close
        if (!this.context.id) {
            this.log.error("Permissions dialog opened without notebook ID.");
            this.dialog.close();
        }

        // This prevents a modal dialog from being closed while it is open.
        this.navigateSubscription = this.router.events.subscribe((event) => {
            if (event.constructor.name === "NavigationStart") {
                this.router.navigate([this.router.url]);
                this.alert = new Alert(AlertType.Warning, "Please close dialog before navigating.");
            }
        });

        // Initialize the permissions service and get our current permissions.
        this.notebookService = notebookService;
        this.notebookService.getNotebookUsers(this.context.id)
            .catch((error: NsError) => {
                this.alert = new Alert(AlertType.Error, "Failed to get notebook users, " + error.Get());
                return Observable.of(null);
            })
            .subscribe((response) => {
                this.users = response;
            });

        this.dataService = getRemotePostData("/ns/v1/user/actions/query", "displayName,email", "displayName", http);
        this.dataService.queryFormatter((arg) => {
            return {
                displayName: arg,
                email: arg,
            };
        });

        // Get our permissions. Note: We remove owner from the list since it can't be configured.
        this.Permissions = PermissionsArray();
        let ownerIndex: number = this.Permissions.indexOf(NotebookPermissions.Owner);
        if (ownerIndex > -1) {
            this.Permissions.splice(ownerIndex, 1);
        }
    }

    public ngOnDestroy() {
        // Whenever we destroy our component, make sure to turn off the dialog watcher.
        this.navigateSubscription.unsubscribe();
    }

    private add(user: User): void {
        // check and see if this user already exists
        for (let userEntry of this.users) {
            if (user.displayName === userEntry.displayName) {
                this.alert = new Alert(AlertType.Error, "User permissions already configured.");
                return;
            }
        }

        let actualUser = this.dataService.lookupUser(user.displayName);
        if (!actualUser) {
            this.alert = new Alert(AlertType.Error, "User not found. Could not add.");
            return;
        }

        user.accountId = actualUser.accountId;
        user.id = actualUser.id;
        user.email = actualUser.email;
        this.usersChanged = true;
        this.users.push(user);
        this.newUser = new User();
    }

    private addDisabled(): boolean {
        return (!this.newUser || !this.newUser.displayName);
    }

    private remove(user: User): void {
        let index = this.users.indexOf(user);
        if (index > -1) {
            this.users.splice(index, 1);
            this.usersChanged = true;
        }
    }

    private submit(): void {
        this.notebookService
            .setNotebookUsers(this.context.id, this.users)
            .catch((error: NsError) => {
                this.alert = new Alert(AlertType.Error, "Failed to set notebook users, " + error.Get());
                return Observable.of(null);
            }).subscribe((response) => {
            if (response) {
                this.dialog.close();
            }
        });
    }

    private exit(): void {
        this.dialog.close();
    }
}
