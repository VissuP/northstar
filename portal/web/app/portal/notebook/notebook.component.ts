import {Component} from "@angular/core";
import {Overlay} from "angular2-modal";
import {Modal} from "angular2-modal/plugins/bootstrap";
import {Observable} from "rxjs/Observable";
import {NotebookService} from "../../shared/services/notebook.service";
import {LoggingService} from "../../shared/services/logging.service";
import {Notebook} from "../../shared/models/notebook.model";
import {AlertService} from "../../shared/services/alerts.service";
import {confirm} from "../../shared/confirm/confirmation";
import {Alert, AlertType} from "../../shared/components/alerts/alerts.model";
import {Router} from "@angular/router";
import {NsError} from "../../shared/models/error.model";

@Component({
    providers: [Modal, Overlay],
    selector: "notebook",
    styleUrls: ["notebook.component.css", "./../../shared/styles/common.css"],
    templateUrl: "notebook.component.html",
})
export class NotebookComponent {
    private notebooks: Notebook[];
    private updates: Notebook[];
    private logs: LoggingService;
    private notebookService: NotebookService;
    private modal: Modal;
    private alerts: AlertService;
    private alert: Alert;
    private router: Router;

    constructor(notebookService: NotebookService, logs: LoggingService, modal: Modal,
                alertService: AlertService, router: Router) {
        this.logs = logs;
        this.modal = modal;
        this.notebookService = notebookService;
        this.alerts = alertService;
        this.router = router;
        this.notebooks = [];

        // returns a list of notebooks accessible via this account.
        this.notebookService.getNotebooks()
            .catch((error: NsError) => {
                this.alert = new Alert(AlertType.Error, "Failed to retrieve notebooks list, " + error.Get());
                return Observable.of(null);
            })
            .subscribe((notebooks: Notebook[]) => {
                if (!notebooks) {
                    return;
                }
                this.updates = notebooks;
                this.updateNotebooks();
            });
    }

    private delete(notebook: Notebook) {
        let title: string = "Delete Notebook";
        let message: string = "Are you sure you want to delete this notebook?";

        confirm(title, message, this.modal)
            .catch((err) => {
                // Catch errors opening dialog
                this.logs.error("Failed to open dialog with error: ", err);
                this.alert = new Alert(AlertType.Error, "Failed to open dialog.");
            })
            .then((dialog) => {
                return (dialog as any).result;
            })
            .then((result) => {
                // Returns the result of the execution.
                if (result) {
                    this.notebookService.deleteNotebook(notebook.id)
                        .catch((error: NsError) => {
                            this.alert = new Alert(AlertType.Error, "Failed to delete notebook, " + error.Get());
                            return Observable.of(null);
                        })
                        .subscribe((deleted: boolean) => {
                            if (deleted) {
                                let index = this.updates.indexOf(notebook, 0);
                                if (index > -1) {
                                    this.updates.splice(index, 1);
                                }
                                this.updateNotebooks();
                            }
                        });
                }
            })
            .catch(() => {
                // Catch dialog cancel.
                this.logs.debug("User canceled notebook deletion.");
            });
    }

    // Helper method used to update notebooks. Note that the two collections are
    // needed because the whole collection associated with the table needs to be
    // change. See https://github.com/swimlane/angular2-data-table/blob/master/demo/basic/filter.ts
    private updateNotebooks() {
        let updates = [...this.updates];
        this.notebooks = updates;
    }

    // loadNotebook uses a resolver to load the notebook, catches any errors that occur,
    // and then routes on to the editor.
    private loadNotebook(notebookId: string): void {
        this.router.navigate(["/portal/editor", notebookId]).catch((error: NsError) => {
            this.logs.error("Editor router navigate returned error: ", error);
            this.alert = new Alert(AlertType.Error, error.Get());
        });
    }

    private import(event: Event) {
        let files: FileList = (event.target as any).files;
        this.updates = this.notebooks;

        // TODO: Add better notebook validation
        // Loop through the uploaded files
        for (let i = 0; i < files.length; i++) {
            let file = files[i];

            // For each file, get a reader to read the contents.
            let reader: FileReader = new FileReader();
            reader.onload = (input: ProgressEvent) => {
                let json: Object;

                // Attempt to parse the uploaded notebook as JSON. Confirms that the file is at least JSON.
                try {
                    json = JSON.parse((input.target as any).result);
                } catch (error) {
                    // If we failed to parse, notify the user.
                    this.alert = new Alert(AlertType.Error, "Failed to parse notebook: " + error);
                    return;
                }

                let notebook = new Notebook(json);

                // Notify the user if there are any backend errors.
                this.notebookService.createNotebook(notebook)
                    .catch((error: NsError) => {
                        this.alert = new Alert(AlertType.Error, "Failed to import notebook, " + error.Get());
                        return Observable.of(null);
                    })
                    .subscribe((createdNotebook: Notebook) => {
                            this.updates.push(createdNotebook);
                            this.updateNotebooks();
                        },
                    );
            };
            reader.readAsText(file);
        }
    }
}
