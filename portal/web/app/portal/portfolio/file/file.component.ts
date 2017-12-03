import {Component, OnDestroy} from "@angular/core";
import {Subscription} from "rxjs/Subscription";
import {ActivatedRoute, Router} from "@angular/router";
import {overlayConfigFactory} from "angular2-modal";
import {BSModalContext, Modal} from "angular2-modal/plugins/bootstrap";
import {AlertService} from "../../../shared/services/alerts.service";
import {LoggingService} from "../../../shared/services/logging.service";
import {PortfolioService} from "../../../shared/services/portfolio.service";
import {Alert, AlertType} from "../../../shared/components/alerts/alerts.model";
import {NsError} from "../../../shared/models/error.model";
import {Portfolio} from "../../../shared/models/portfolio.model";
import {File} from "../../../shared/models/file.model";

@Component({
    selector: "file",
    styleUrls: ["file.component.css", "./../../../shared/styles/common.css"],
    templateUrl: "file.component.html",
})
export class FileComponent implements OnDestroy {
    private portfolio: Portfolio;
    private log: LoggingService;
    private alerts: AlertService;
    private portfolioService: PortfolioService;
    private modal: Modal;
    private router: Router;
    private alert: Alert;
    private running: boolean;
    private lastRun: Date;
    private authRefreshObservable: Subscription;

    constructor(route: ActivatedRoute,
                portfolioService: PortfolioService,
                logger: LoggingService,
                alertService: AlertService,
                modal: Modal,
                router: Router) {
        this.portfolio = new Portfolio();
        this.log = logger;
        this.alerts = alertService;
        this.portfolioService = portfolioService;
        this.modal = modal;
        this.router = router;
        this.running = false;

        // If our portfolio has been populated by the resolver, use that.
        if (route.snapshot.data.portfolio) {
            this.portfolio = route.snapshot.data.portfolio;
        } 
    };

    private downloadFile(portfolioName: string, fileName: string): void { 
        let fileUrl = "/ns/v1/portfolios/"+portfolioName+"/"+fileName;
        window.open(fileUrl);
    }

    // Implements OnDestroy.
    public ngOnDestroy(): void {
        if (this.authRefreshObservable) {
            this.log.debug("Removing auth-refresh observer.")
            this.authRefreshObservable.unsubscribe();
        }
    }
}
