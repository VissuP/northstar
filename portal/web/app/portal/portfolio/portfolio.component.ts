import {Component} from "@angular/core";
import {Overlay} from "angular2-modal";
import {Modal} from "angular2-modal/plugins/bootstrap";
import {Observable} from "rxjs/Observable";
import {PortfolioService} from "../../shared/services/portfolio.service";
import {LoggingService} from "../../shared/services/logging.service";
import {Portfolio} from "../../shared/models/portfolio.model";
import {AlertService} from "../../shared/services/alerts.service";
import {confirm} from "../../shared/confirm/confirmation";
import {Alert, AlertType} from "../../shared/components/alerts/alerts.model";
import {Router} from "@angular/router";
import {NsError} from "../../shared/models/error.model";

@Component({
    providers: [Modal, Overlay],
    selector: "portfolio",
    styleUrls: ["portfolio.component.css", "./../../shared/styles/common.css"],
    templateUrl: "portfolio.component.html",
})
export class PortfolioComponent {
    private portfolios: Portfolio[];
    private logs: LoggingService;
    private portfolioService: PortfolioService;
    private modal: Modal;
    private alerts: AlertService;
    private alert: Alert;
    private router: Router;

    constructor(portfolioService: PortfolioService, logs: LoggingService, modal: Modal,
                alertService: AlertService, router: Router) {
        this.logs = logs;
        this.modal = modal;
        this.portfolioService = portfolioService;
        this.alerts = alertService;
        this.router = router;
        this.portfolios = [];

        // returns a list of portfolios accessible via this account.
        this.portfolioService.getPortfolios()
            .catch((error: NsError) => {
                this.alert = new Alert(AlertType.Error, "Failed to retrieve portfolios list, " + error.Get());
                return Observable.of(null);
            })
            .subscribe((portfolios: Portfolio[]) => {
                if (!portfolios) {
                    return;
                }
                this.portfolios = portfolios;
            });
    }

    // loadPortfolio uses a resolver to load the notebook, catches any errors that occur,
    // and then routes on to the editor.
    private loadPortfolio(portfolioName: string): void {
        this.router.navigate(["/portal/portfolios", portfolioName]).catch((error: NsError) => {
            this.logs.error("Editor router navigate returned error: ", error);
            this.alert = new Alert(AlertType.Error, error.Get());
        });
    }
}
