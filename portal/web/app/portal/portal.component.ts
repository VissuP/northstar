import {RouterOutlet, Router} from "@angular/router";
import {Component} from "@angular/core";
import {HttpModule} from "@angular/http";
import {FormsModule} from "@angular/forms";
import {Observable, Subject} from "rxjs/Rx";
import {AuthGuard} from "../shared/auth/auth-guard.service";
import {UserService} from "../shared/services/user.service";
import {AlertService} from "../shared/services/alerts.service";
import {DisconnectReason, EventService} from "ngx-vz-websocket";
import {LoggingService} from "../shared/services/logging.service";
import {User} from "../shared/models/user.model";
import {Event, EventType} from "ngx-vz-websocket";
import {NsError} from "../shared/models/error.model";
import {Utilities} from "../shared/resources/utilities";

// Defines the connections API url.
const connectionUrl: string = "/ns/v1/connections/";

@Component({
    selector: "app",
    templateUrl: "portal.component.html",
    styleUrls: ["portal.component.css"],
})
export class PortalComponent {
    private router: Router;
    private authGuard: AuthGuard;
    private userService: UserService;
    private connections: EventService;
    private log: LoggingService;
    private alertService: AlertService;
    private mobileView: number = 992;
    private toggle: boolean = false;
    private title: String = "NorthStar";
    private user: User = new User();

    constructor(router: Router, authGuard: AuthGuard, userService: UserService, logger: LoggingService,
                alerts: AlertService, connections: EventService) {
        this.router = router;
        this.authGuard = authGuard;
        this.userService = userService;

        this.connections = connections;
        this.connections.setURL(this.getWebsocketUrl());
        this.log = logger;
        this.alertService = alerts;

        this.attachEvents();

        // Get the authenticated user.
        userService.getUser()
            .catch((error: NsError) => {
                this.log.error("Failed to get user with error " + error.Get());
                return Observable.of(new User());
            })
            .subscribe((user: User) => {
                this.user = user;
            });

        // Set toggle value.
        this.setToggle();
    }

    public getTitle(): String {
        return this.title;
    }

    public closeAlert(index: number): void {
        this.alertService.closeAlert(index);
    }

    public getWidth(): number {
        return window.innerWidth;
    }

    public toggleSidebar(): void {
        this.toggle = !this.toggle;
        localStorage.setItem("toggle", this.toggle.toString());
    }

    public setToggle(): void {
        if (this.getWidth() >= this.mobileView) {
            if (localStorage.getItem("toggle")) {
                this.toggle = localStorage.getItem("toggle") === "true" ? true : false;
            } else {
                this.toggle = true;
            }
        } else {
            this.toggle = false;
        }
    }

    // Logs out the authenticated user.
    public logout(): void {

        // Unregister from receive async events.
        let event = new Event();
        event.type = EventType.UnregisterObserver;
        event.id = Utilities.NewGuid();

        this.connections.write(event);

        this.authGuard.logout()
            .subscribe((response: any) => {
                // User logged out. Close their connection.
                this.connections.disconnect(DisconnectReason.NORMAL_CLOSURE, "User logged out");
                this.router.navigate(["/login"]);
            });
    }

    private attachEvents(): void {
        window.onresize = () => {
            this.setToggle();
        };
    }

    // Using (click) doesn't allow changing of variables. Need to abstract to a function to allow us to navigate.
    private navigate(href: string): void {
        window.location.href = href;
    }

    // We need to build the websocket address so that it points at the host of the service.
    // If we used https to access the service, use wss (secure). Otherwise ues ws.
    private getWebsocketUrl(): string {
        let protocol = "ws";
        if (window.location.protocol === "https:") {
            protocol = "wss";
        }
        return protocol + "://" + window.location.host + connectionUrl;
    }
}
