import {Injectable} from "@angular/core";
import {LoggingService} from "./logging.service";
import {Alert} from "../components/alerts/alerts.model";

@Injectable()
export class AlertService {
    public alerts: Alert[];
    private logs: LoggingService;

    constructor(private logger: LoggingService) {
        this.logs = logger;
        this.alerts = new Array<Alert>();
    }

    public addAlert(alert: Alert) {
        this.alerts.push(alert);
    }

    public closeAlert(index: number) {
        this.alerts.splice(index, 1);
    }
}