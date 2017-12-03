import {Component, Input, OnDestroy} from "@angular/core";
import {DialogRef, ModalComponent} from "angular2-modal";
import {BSModalContext} from "angular2-modal/plugins/bootstrap";
import {Router} from "@angular/router";
import "rxjs/add/operator/pairwise";
import "codemirror/mode/javascript/javascript";
import "codemirror/mode/lua/lua";
import "codemirror/mode/r/r";
import "codemirror/mode/xml/xml";
import "codemirror/mode/htmlmixed/htmlmixed";
import "codemirror/addon/fold/foldcode";
import "codemirror/addon/fold/foldgutter";
import "codemirror/addon/fold/brace-fold";
import "codemirror/addon/fold/xml-fold";
import "codemirror/addon/fold/markdown-fold";
import "codemirror/addon/fold/comment-fold";
import "codemirror/addon/selection/active-line";
import {Observable} from "rxjs/Observable";
import {LoggingService} from "../../../shared/services/logging.service";
import {EventSchema, Transformation, TransformationEvent} from "../../../shared/models/transformation.model";
import {
    EditorConfig,
    ExecutionOutput,
    LanguageConfig,
    LanguageMode,
} from "ngx-vz-cell";
import {LanguageService} from "../../../shared/services/language.service";
import {TransformationService} from "../transformation.service";
import {Category, CategoryOption, CategoryOptionField} from "../../../shared/models/schedule.model";
import {EventService} from "ngx-vz-websocket";
import {Utilities} from "../../../shared/resources/utilities";
import {Event, EventType} from "ngx-vz-websocket";
import {Alert, AlertType} from "../../../shared/components/alerts/alerts.model";
import {NsError} from "../../../shared/models/error.model";
import {Subscription} from "rxjs";

export class TransformationModalContext extends BSModalContext {
    public transformation: Transformation;

    // Scheduler
    public schemas: EventSchema[];
    public selection: Transformation;
    public categories: Category[];
    public categoryOptions: CategoryOption[];
    public categoryOptionFields: CategoryOptionField[];
}

@Component({
    selector: "modal-content",
    styleUrls: ["editor.component.css"],
    templateUrl: "editor.component.html",
})
export class TransformationModal implements ModalComponent<TransformationModalContext>, OnDestroy {
    public dialog: DialogRef<TransformationModalContext>;

    @Input() context: TransformationModalContext;
    private codeConfig: LanguageConfig;
    private jsonConfig: LanguageConfig;
    private transformation: Transformation;
    private languages: LanguageConfig[];
    private languageService: LanguageService;
    private log: LoggingService;
    private jsonCode: string;
    private executionID: string;

    private categories: Category[];
    private categoryOptions: CategoryOption[];
    private categoryOptionFields: CategoryOptionField[];

    private selectedCategory: Category;
    private selectedCategoryOption: CategoryOption;
    private selectedCategoryOptionField: CategoryOptionField;

    private transformationService: TransformationService;
    private scheduleChanged: boolean;
    private connections: EventService;

    private alert: Alert;
    private router: Router;
    private navigateSubscription: Subscription;

    private testResults: ExecutionOutput;

    constructor(dialog: DialogRef<TransformationModalContext>,
                log: LoggingService,
                languageService: LanguageService,
                transformationService: TransformationService,
                router: Router,
                connections: EventService) {

        this.router = router;

        // This prevents a modal dialog from being closed while it is open.
        this.navigateSubscription = this.router.events.subscribe((event) => {
            if (event.constructor.name === "NavigationStart") {
                this.router.navigate([this.router.url]);
                this.alert = new Alert(AlertType.Warning, "Please close dialog before navigating.");
            }
        });

        this.connections = connections;
        this.log = log;
        this.transformationService = transformationService;
        this.languageService = languageService;

        // Initialize the modal context.
        this.dialog = dialog;
        this.context = dialog.context;
        this.dialog.setCloseGuard(this);

        // TODO: Wonder if we should make this a separate event type from a regular execution.
        // If a regular execution completes while on this tab, we'll show the output.
        this.connections.observe(EventType.ExecuteResult)
            .subscribe((msg: Event) => {
                if (msg.type === EventType.InternalError) {
                    this.testResults = new ExecutionOutput();
                    this.testResults.setExecutionError("Error, network connection interrupted. Please retry.");
                }
                if (msg.id === this.executionID) {
                    let payload = JSON.parse(atob(msg.payload));
                    this.testResults = new ExecutionOutput(payload);
                }
            });

        this.categories = this.context.categories;
        this.categoryOptions = this.context.categoryOptions;
        this.categoryOptionFields = this.context.categoryOptionFields;

        this.transformation = this.context.selection;
        this.initializeScheduleSettings(this.transformation.schedule.event);

        // Initialize support for the codemirror objects and the select of language.
        this.languages = this.languageService.languageArray(LanguageMode.Code);
        this.codeConfig = this.languageService.configs.get(this.transformation.language);
        this.jsonConfig = this.languageService.configs.get("json");
    }

    public ngOnDestroy(): void {
        // Whenever we destroy our component, make sure to turn off the dialog watcher.
        this.navigateSubscription.unsubscribe();
    }

    // TODO: See if we can simplify this.
    private initializeScheduleSettings(event: TransformationEvent): void {
        for (let category of this.categories) {
            if (category.name === event.type) {
                this.selectedCategory = category;
            }
        }

        if (this.selectedCategory && this.selectedCategory.name === "Timer") {
            for (let categoryOption of this.categoryOptions) {
                if (categoryOption.name === event.name && categoryOption.value === event.value) {
                    this.selectedCategoryOption = categoryOption;
                }
            }
        } else if (this.selectedCategory && this.selectedCategory.name === "Device") {
            for (let categoryOptionField of this.categoryOptionFields) {
                if (categoryOptionField.name === event.name && categoryOptionField.value === event.value) {
                    this.selectedCategoryOptionField = categoryOptionField;
                }
            }
        }
    }

    private save(): void {
        this.context.selection = this.transformation;
        if (!this.parseSchedule()) {
            return;
        }

        if (this.transformation.id) {
            this.transformationService.updateTransformation(this.transformation)
                .catch((error: NsError) => {
                    this.alert = new Alert(AlertType.Error, "Failed to update transformation, " + error.Get());
                    return Observable.of(null);
                }).subscribe((data) => {
                    if (data) {
                        this.dialog.close(this.transformation);
                    }
            });
        } else {
            this.transformationService.createTransformation(this.transformation)
                .catch((error: NsError) => {
                    this.alert = new Alert(AlertType.Error, "Failed to create transformation, " + error.Get());
                    return Observable.of(null);
                }).subscribe((data) => {
                if (data) {
                    this.dialog.close(data);
                }
            });
        }
    }

    private exit(): void {
        this.dialog.close();
    }

    private parseSchedule(): boolean {
        if (!this.selectedCategory) {
            this.alert = new Alert(AlertType.Error, "Event category not set. Cannot save");
            return false;
        }

        this.transformation.schedule.event.type = this.selectedCategory.name;

        if (this.selectedCategory.name === "Timer") {
            if (!this.selectedCategoryOption) {
                this.alert = new Alert(AlertType.Error, "Event category options not set. Cannot save.");
                return false;
            }
            this.transformation.schedule.event.name = this.selectedCategoryOption.name;
            this.transformation.schedule.event.value = this.selectedCategoryOption.value;
            return true;
        } else if (this.selectedCategory.name === "Device") {
            if (!this.selectedCategoryOptionField) {
                this.alert = new Alert(AlertType.Error, "Event field not set. Cannot save.");
                return false;
            }
            this.transformation.schedule.event.name = this.selectedCategoryOptionField.name;
            this.transformation.schedule.event.value = this.selectedCategoryOptionField.value;
            return true;
        } else if (this.selectedCategory.name === "None") {
            return true;
        }

        return null;
    }

    private executeTransformation(): void {
        // If input provided, add to transformation.
        if (this.jsonCode) {
            try {
                // TODO - We should update the UI. Free form is not the proper UI element.
                this.transformation.arguments = JSON.parse(this.jsonCode);
            } catch (error) {
                this.alert = new Alert(AlertType.Error, "Failed to parse arguments as JSON. error: ", error);
            }
        }

        // Create the event for the execution.
        let event = new Event();

        event.type = EventType.ExecuteTransformation;
        event.id = Utilities.NewGuid();
        event.payload = this.transformation.encode();

        // Keep track of our execution ID locally so that we know that the result belongs to us.
        this.executionID = event.id;

        // Send event for execution.
        this.connections.write(event);
    }

    // CodeMirror does some calculations based on the window size when it loads.
    // Because we are using tabs, the codemirror element gets improper dimensions.
    // Refresh here to fix it.
    private activateCodeMirrorByID(id: string, config: EditorConfig): void {
        // NOTE: Angular strongly encourages you not to play with the DOM directly.
        // However, this fixes an issue where codemirror doesn't initialize unless visible at load.
        let CodeMirror = document.getElementById(id);

        if (CodeMirror) {
            let innerCodeMirror = CodeMirror.getElementsByClassName("CodeMirror")[0];
            if (innerCodeMirror && (innerCodeMirror as any).CodeMirror) {
                // Wait for the codemirror to load and refresh it. Note that setTimeout only runs once.
                setTimeout(() => {
                    (innerCodeMirror as any).CodeMirror.refresh();
                }, 20);
            }
        }
    }

    // setCodeMirrorField sets individual fields on the codemirror object. Note that the codemirror object
    // doesn't work with data binding, so this must be done manually.
    private setCodeMirrorField(id: string, field: string, value: any): void {
        // NOTE: Angular strongly encourages you not to play with the DOM directly.
        // However, codemirror doesn't recognize binding properly.
        let CodeMirror = document.getElementById(id);

        if (CodeMirror) {
            let innerCodeMirror = CodeMirror.getElementsByClassName("CodeMirror")[0];
            if (innerCodeMirror && (innerCodeMirror as any).CodeMirror) {
                (innerCodeMirror as any).CodeMirror.setOption(field, value);
            }
        }
    }

    private languageChanged(language: string): void {
        this.codeConfig = this.languageService.configs.get(language);
        this.setCodeMirrorField("code", "mode", (this.codeConfig as any).codemirror.mode);
    }

    // TODO: Change column names depending on selected category
    private categorySelected(): boolean {
        return this.selectedCategory.name !== "None";
    }

    private categoryChanged(): void {
        this.selectedCategoryOption = null;
        this.selectedCategoryOptionField = null;
        this.scheduleChanged = true;
    }

    private categoryOptionsChanged(): void {
        this.selectedCategoryOptionField = null;
        this.scheduleChanged = true;
    }

    private categoryOptionsFieldChanged(): void {
        this.scheduleChanged = true;
    }
}
