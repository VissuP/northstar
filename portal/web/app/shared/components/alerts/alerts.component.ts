import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import {Alert} from "./alerts.model";

@Component({
    selector: "alerts",
    templateUrl: "alerts.component.html",
    styleUrls: ["alerts.component.css"],
})
export class AlertsComponent implements OnChanges {
    @Input() message: Alert;
    private timeoutId: number;

    constructor() {

    }

    public closeAlert() {
        this.message = null;
    }

    public ngOnChanges(changes: SimpleChanges) {
        // When an input object changes, we get a notification. Make sure the notification parameter is present. Finally, make sure that the value we get is not null.
        if (changes.message && changes.message.currentValue) {
            // clear the old timeout if it's still running
            clearTimeout(this.timeoutId);

            // Set a timeout for the alert's timeout.
            this.timeoutId = window.setTimeout(() => {
               this.message = null;
            }, this.message.timeout);
        }
    }
}
