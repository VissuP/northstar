import {Component, ViewChild} from "@angular/core";
import {ActivatedRoute} from "@angular/router";
import {TransformationService} from "../transformation.service";
import {LoggingService} from "../../../shared/services/logging.service";
import {ExecutionOutput} from "ngx-vz-cell";

@Component({
    selector: "transformation-history",
    styleUrls: ["results.component.css"],
    templateUrl: "results.component.html",
})
export class TransformationResults {

    private transformationService: TransformationService;
    private log: LoggingService;
    private history: ExecutionOutput[];

    @ViewChild("historyTable") table;

    constructor(route: ActivatedRoute, transformationService: TransformationService, log: LoggingService) {
        this.transformationService = transformationService;
        this.log = log;
        this.history = route.snapshot.data.history;
    }

    // Helper method used to Expand/Collapse row details.
    private toggleExpandRow(row: any): void {
        this.table.rowDetail.toggleExpandRow(row);
    }
}
