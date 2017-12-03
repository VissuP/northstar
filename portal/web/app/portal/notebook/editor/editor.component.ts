import {Component, OnDestroy} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {Subscription} from "rxjs/Subscription";
import {ActivatedRoute, Router} from "@angular/router";
import {overlayConfigFactory} from "angular2-modal";
import {BSModalContext, Modal} from "angular2-modal/plugins/bootstrap";
import {AlertService} from "../../../shared/services/alerts.service";
import {LoggingService} from "../../../shared/services/logging.service";
import {EventService} from "ngx-vz-websocket";
import {PermissionsModal} from "./permissions/permissions.component";
import {Cell, ExecutionOutput, StatusIndicator} from "ngx-vz-cell";
import {LanguageConfig, LanguageMode, LanguageService} from "../../../shared/services/language.service";
import {Event, EventType} from "ngx-vz-websocket";
import {NotebookService} from "../../../shared/services/notebook.service";
import {Template} from "../../../shared/models/template.model";
import {TemplateService} from "../../../shared/services/template.service";
import {Alert, AlertType} from "../../../shared/components/alerts/alerts.model";
import {NsError} from "../../../shared/models/error.model";
import {Notebook, NotebookCell} from "../../../shared/models/notebook.model";


@Component({
    selector: "editor",
    styleUrls: ["editor.component.css", "./../../../shared/styles/common.css"],
    templateUrl: "editor.component.html",
})
export class EditorComponent implements OnDestroy {
    private notebook: Notebook;
    private activeCell: Cell;
    private log: LoggingService;
    private alerts: AlertService;
    private notebookService: NotebookService;
    private templateService: TemplateService;
    private connections: EventService;
    private modal: Modal;
    private languageService: LanguageService;
    private languages: LanguageConfig[];
    private router: Router;
    private alert: Alert;
    private running: boolean;
    private lastRun: string;
    private authRefreshObservable: Subscription;

    private templates: Template[];

    constructor(route: ActivatedRoute,
                notebookService: NotebookService,
                templateService: TemplateService,
                logger: LoggingService,
                alertService: AlertService,
                modal: Modal,
                languageService: LanguageService,
                router: Router,
                connections: EventService) {
        this.notebook = new Notebook();
        this.activeCell = null;
        this.log = logger;
        this.alerts = alertService;
        this.notebookService = notebookService;
        this.templateService = templateService;
        this.connections = connections;
        this.modal = modal;
        this.languageService = languageService;
        this.languages = languageService.languageArray(LanguageMode.Code, LanguageMode.Static, LanguageMode.Query);
        this.router = router;
        this.running = false;

        // If our notebook has been populated by the resolver, use that. Otherwise create a new notebook
        if (route.snapshot.data.notebook) {
            this.notebook = route.snapshot.data.notebook;
        } else if (!(route.params as any).id) {
            this.notebook = new Notebook();
        }

        // Get templates.
        templateService.getTemplates()
            .catch((error: NsError) => {
                this.alert = new Alert(AlertType.Error, error.Get());
                return Observable.of(null);
            })
            .subscribe((templates: Template[]) => {
                if (templates) {
                    this.templates = templates;
                }
            });

        this.connections.observe(EventType.ExecuteResult)
            .subscribe((msg: Event) => {
                this.processExecutionResult(msg);
            });

        // Register to receive error events.
        this.connections.observe(EventType.Error)
            .subscribe((errorEvent: Event) => {
                this.log.error("Error received on socket:", errorEvent);
                let output = JSON.parse(atob(errorEvent.payload));
                this.notebook.setRunningCellsTo(StatusIndicator.Error, output.Value);
            });

    };

    // Implements OnDestroy.
    public ngOnDestroy(): void {
        if (this.authRefreshObservable) {
            this.log.debug("Removing auth-refresh observer.");
            this.authRefreshObservable.unsubscribe();
        }
    }

    // processExecutionResult is called when an event of type executionResult is received on the websocket.
    // It transforms the content so the notebook can use it.
    private processExecutionResult(result: Event) {
        // If we receive an internal error on the notebook, stop all executing cells.
        if (result.type === EventType.InternalError) {
            this.log.error("Error received on socket:", result);
            this.notebook.setRunningCellsTo(StatusIndicator.Error, "Network connection lost. Please try again.");
            return;
        }

        let cell: Cell = this.notebook.cells.find((notebookCell) => notebookCell.id === result.id);

        if (cell) {
            // This shouldn't be here.... shouldn't be encoded.
            cell.output = new ExecutionOutput(JSON.parse(atob(result.payload)));

            // TODO - Business logic does not belong in the component. It belongs in
            // the model.

            if (cell.output) {
                if (cell.output.failedExecution()) {
                    cell.options.status = StatusIndicator.Error;
                    return;
                }
                cell.options.status = StatusIndicator.Success;
            }
        }
    }

    // moveCellUp takes the selected cell and swaps it with the cell above it.
    private moveCellUp(cell: Cell): void {
        let cellIndex: number = this.notebook.cells.findIndex((notebookCell) => notebookCell.id === (cell as NotebookCell).id);
        if (cellIndex !== 0) {
            // This is a one line implementation of a swap. See
            // https://basarat.gitbooks.io/typescript/content/docs/destructuring.html
            [this.notebook.cells[cellIndex], this.notebook.cells[cellIndex - 1]] =
                [this.notebook.cells[cellIndex - 1], this.notebook.cells[cellIndex]];
        }
    }

    // moveCellDown takes the selected cell and swaps it with the cell below it.
    private moveCellDown(cell: Cell): void {
        let cellIndex: number = this.notebook.cells.findIndex((notebookCell) => notebookCell.id === (cell as NotebookCell).id);
        if (cellIndex !== this.notebook.cells.length - 1) {
            // This is a one line implementation of a swap. See
            // https://basarat.gitbooks.io/typescript/content/docs/destructuring.html
            [this.notebook.cells[cellIndex], this.notebook.cells[cellIndex + 1]] =
                [this.notebook.cells[cellIndex + 1], this.notebook.cells[cellIndex]];
        }
    }

    // executeCell forms and sends out an execution request for a cell.
    private executeCell(cell: Cell): void {
        this.log.debug("executeCell:", cell);

        if (!this.languageService.configs.get(cell.language).isExecutable()) {
            this.log.error("Cannot execute a non-executable cell.");
            return;
        }

        cell.options.status = StatusIndicator.Running;

        // Base64 the cell code before. Note that this
        // is needed to prevent errors while sending
        // in JSON.
        let payload = new NotebookCell(cell);
        payload.code = btoa(cell.code);

        // Create the event. Sending cell as JSON payload.
        let event = new Event();
        event.type = EventType.ExecuteCell;
        event.id = payload.id;
        event.payload = payload;

        this.log.debug("Sending event", event);
        this.connections.write(event);
    }

    // executeAll triggers an execute on each cell successively, giving up when the response
    // is an error or we've reached the end of the list.
    private executeAll() {
        this.log.debug("executeAll");

        // Start by getting and executing the first cell in the array
        let cell = this.getExecutableCell(0);
        if (cell) {
            this.executeCell(cell);
        } else {
            // If we don't have any executable cells don't continue.
            this.log.debug("Failed to find any executable cells in:", this.notebook.cells);
            return;
        }

        let execution = this.connections.observe(EventType.ExecuteResult)
            .subscribe(
                (executionResult: Event) => {
                    // Check our payload response for a value in stderr. If we have a value,
                    // this cell didn't execute properly so cancel the executeAll.
                    let payload: ExecutionOutput = new ExecutionOutput(JSON.parse(atob(executionResult.payload)));
                    if (payload.failedExecution()) {
                        this.log.debug("ExecuteAll received error response:", payload.getExecutionError());
                        execution.unsubscribe();
                        return;
                    }

                    // If we're not able to execute our next cell, our execute all is over. Cancel the subscription.
                    if (!this.executeNext(executionResult)) {
                        this.log.debug("No more cells to execute. Canceling subscription...");
                        execution.unsubscribe();
                        return;
                    }
                },
                (error) => {
                    // We will probably never hit this, but if we do, we're not going to
                    // process any more so unsubscribe.
                    this.log.error("Execution failed with error: ", error);
                    execution.unsubscribe();
                    return;
                },
            );
    }

    // setAutoRefresh handles changes in auto refresh ui.
    private setAutoRefresh(value: string): void {
        let interval: number = parseInt(value, 10);

        this.running = false;

        // Stop auto-refresh observable.
        if (this.authRefreshObservable) {
            this.authRefreshObservable.unsubscribe();
        }

        // If the interval is greater than 0, create timer.
        if (interval > 0) {
            this.running = true;
            Observable.timer(0, interval).subscribe(() => {
                this.lastRun = new Date().toLocaleString();
                this.log.debug("Refreshing cells ", this.lastRun);
                this.executeAll();
            });
        }
    }

    // executeNext finds the cell after the provided cell and executes it.
    private executeNext(executionResponse: Event): boolean {
        this.log.debug("executeNext: event: ", executionResponse);

        // Get the index of the cell we just executed
        let cellIndex = this.notebook.cells.findIndex((cell) => cell.id === executionResponse.id);

        // Find our next executable cell (if it exists) and execute it.
        let cell = this.getExecutableCell(cellIndex + 1);
        if (cell != null) {
            this.executeCell(cell);
            return true;
        }

        return false;
    }

    // get next executable cell starting at index
    private getExecutableCell(cellIndex: number): Cell {
        this.log.debug("getExecutableCell: index: ", cellIndex);
        for (let i = cellIndex; i < this.notebook.cells.length; i++) {
            let cell = this.notebook.cells[i];
            if (this.languageService.configs.get(cell.language).isExecutable()) {
                this.log.debug("Found executable cell: ", cell);
                return this.notebook.cells[i];
            }
        }

        return null;
    }

    // saveNotebook saves the notebook currently loaded in the editor.
    private saveNotebook(): void {
        let notebookCopy = new Notebook(this.notebookService.encodeNotebook(this.notebook));

        // We don't want to save a notebook with a status of running since it won't be
        // running when it loads. Set to error instead so the user knows to rerun.
        notebookCopy.setRunningCellsTo(StatusIndicator.Error, "Execution interrupted. Please rerun.");

        // TODO - We need to fix this. It should be always doing PUT
        // on the back end. PUT should create or update.
        if (!notebookCopy.id) {
            this.notebookService.createNotebook(notebookCopy)
                .catch((error: NsError) => {
                    this.alert = new Alert(AlertType.Error, "Failed to create notebook, " + error.Get());
                    return Observable.of(null);
                })
                .subscribe((notebook: Notebook) => {
                    if (notebook) {
                        this.notebook = this.notebookService.decodeNotebook(notebook); // decode the notebook returned after saving
                        this.alert = new Alert(AlertType.Success, "Saved notebook successfully");

                        // Update the URL to represent the new notebook ID
                        history.pushState(null, null, window.location + "/" + this.notebook.id);
                    }
                });

        } else {
            this.notebookService.updateNotebook(notebookCopy)
                .catch((error: NsError) => {
                    this.alert = new Alert(AlertType.Error, "Failed to update notebook, " + error.Get());
                    return Observable.of(null);
                })
                .subscribe((result: boolean) => {
                    if (result) {
                        this.alert = new Alert(AlertType.Success, "Saved notebook successfully");
                    }
                });
        }
    }

    // setActiveCell is called whenever a cell is selected in the editor, it allows us keep track
    // of where we think the user is so for things like styling (i.e. sidebar on selected cell)
    private setActiveCell(cell: Cell): void {
        this.activeCell = cell;
    }

    // configurePermissions opens a dialog to allow the user to set access permissions for this notebook.
    private configurePermissions(): void {

        let modal = this.modal.open(PermissionsModal, overlayConfigFactory({id: this.notebook.id}))
            .catch((err: string) => {
                this.log.error("Could not open permissions modal. Error was:", err);
                // TODO: Consider showing this as a banner
            });
    }

    private downloadNotebook(): void {
        let notebookBlob = new Blob([JSON.stringify(this.notebook)], {type: "application/vznotebook"});
        let notebookUrl = window.URL.createObjectURL(notebookBlob);
        window.open(notebookUrl);
        window.URL.revokeObjectURL(notebookUrl);
    }

    // adds new cell via selected template by first decoding the cell code for UI, if necessary
    private addDecodedTemplate(template: Template) {
        let cell = new Cell(template.getData());
        if (this.notebookService.isBase64(cell.code)) { // decode only if code is base64 encoded
            cell.code = atob(cell.code);
        }
        this.notebook.addCell(cell);
    }
}
