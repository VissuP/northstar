import {Component, ViewChild} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {Overlay, overlayConfigFactory} from "angular2-modal";
import {BSModalContext, Modal} from "angular2-modal/plugins/bootstrap";
import {LoggingService} from "../../shared/services/logging.service";
import {AlertService} from "../../shared/services/alerts.service";
import {TransformationService} from "./transformation.service";
import {TransformationModal} from "./editor/editor.component";
import {EventSchema, Transformation} from "../../shared/models/transformation.model";
import {Category, CategoryOption, CategoryOptionField} from "../../shared/models/schedule.model";
import {Router} from "@angular/router";
import {confirm} from "../../shared/confirm/confirmation";
import {Alert, AlertType} from "../../shared/components/alerts/alerts.model";
import {NsError} from "../../shared/models/error.model";

@Component({
    providers: [Modal, Overlay],
    selector: "transformation",
    styleUrls: ["transformation.component.css", "./../../shared/styles/common.css"],
    templateUrl: "transformation.component.html",
})
export class TransformationComponent {
    public modal: Modal;

    private transformationService: TransformationService;
    private log: LoggingService;
    private alertService: AlertService;
    private router: Router;

    private transformations: Transformation[] = [];
    private updates: Transformation[] = [];
    private schemas: EventSchema[];
    private categories: Category[];
    private categoryOptions: CategoryOption[];
    private categoryOptionFields: CategoryOptionField[];
    private alert: Alert;

    @ViewChild("transformationtable") table;

    constructor(transformationService: TransformationService, log: LoggingService,
                alertService: AlertService, modal: Modal, router: Router) {
        this.transformationService = transformationService;
        this.log = log;
        this.alertService = alertService;
        this.modal = modal;
        this.router = router;

        // Get event schemas.
        transformationService.getEventSchemas()
            .catch((error: NsError) => {
                this.alert = new Alert(AlertType.Error, "Failed to get event schemas, " + error.Get());
                return Observable.of(null);
            })
            .subscribe((schemas: EventSchema[]) => {
                if (schemas) {
                    this.schemas = schemas;
                    this.categories = this.getCategories(this.schemas);
                    this.categoryOptions = this.getCategoryOptions(this.schemas);
                    this.categoryOptionFields = this.getCategoryOptionFields(this.schemas);
                }
            });

        // Get transformations.
        transformationService.getTransformations()
            .catch((error: NsError) => {
                this.alert = new Alert(AlertType.Error, "Failed to get transformation, " + error.Get());
                return Observable.of(null);
            })
            .subscribe((transformations: Transformation[]) => {
                if (transformations) {
                    this.updates = transformations;
                    this.updateTransformations();
                }
            });
    }

    private delete(transformation: Transformation): void {
        let title: string = "Delete Transformation";
        let message: string = "Are you sure you want to delete this transformation?";

        confirm(title, message, this.modal)
            .catch((err) => {
                // Catch errors opening dialog
                this.log.error("Failed to open dialog with error: ", err);
                this.alert = new Alert(AlertType.Error, "Failed to open dialog.");
            })
            .then((dialog) => {
                return (dialog as any).result;
            })
            .then((result) => {
                // Returns the result of the execution.
                if (result) {
                    return this.transformationService.deleteTransformation(transformation.id)
                        .catch((error: NsError) => {
                            this.alert = new Alert(AlertType.Error, "Failed to delete transformation, " + error.Get());
                            return Observable.of(null);
                        }).subscribe((response: boolean) => {
                            if (response) {
                                let index = this.updates.indexOf(transformation);
                                if (index > -1) {
                                    this.updates.splice(index, 1);
                                }
                                this.updateTransformations();
                            }
                        });
                } else {
                    return null;
                }
            })
            .catch(() => {
                // Catch dialog cancel.
                this.log.debug("User canceled transformation deletion.");
            });
    }

    // loadResults uses a resolver to load the transformation history, catches any errors that occur,
    // and then routes on to the transformation history page.
    private loadResults(transformationId: string): void {
        this.router.navigate(["/portal/transformations/" + transformationId + "/results"]).catch((error: NsError) => {
            this.alert = new Alert(AlertType.Error, "Failed to retrieve transformation results, " + error.Get());
        });
    }

    // Opens the transformation dialog for user to create transformation.
    // Note that this is used for both adding and editing.
    private launchModal(transformation: Transformation): void {
        let newTransformation: Transformation = new Transformation();
        let update: boolean = false;

        if (transformation) {
            newTransformation = new Transformation(transformation.encode());
            update = true;
        }

        let modal = this.modal.open(TransformationModal, overlayConfigFactory({
            categories: this.categories,
            categoryOptionFields: this.categoryOptionFields,
            categoryOptions: this.categoryOptions,
            selection: newTransformation,
        }));
        modal
            .catch((error) => {
                this.log.error("Open transformation modal returned error: ", error);
            })
            .then((dialog) => {
                return ( dialog as any ).result;
            })
            .then((result: Transformation) => {
                // Make sure result is not nil. E.g., the user did not cancel selection.
                if (result) {
                    result.setScheduled();
                    result.setState();

                    // If we have an id, we're updating
                    if (update) {
                        let index = this.updates.findIndex((element) => {return element.id === result.id});
                        if (index > -1) {
                            this.updates[index] = result;
                        }
                    } else {
                        this.updates.push(result);
                    }

                    // Regardless, update the list of transformations
                    this.updateTransformations();
                }
            });
    }

    // Helper method used to Expand/Collapse row details.
    private toggleExpandRow(row: any): void {
       this.table.rowDetail.toggleExpandRow(row);
    }


    // Helper method used to update transformations.
    private updateTransformations(): void {
        let updates = [...this.updates];
        this.transformations = updates;
    }

    private getCategories(schemas: EventSchema[]): Category[] {
        let index = 0;
        let categoryNames = new Array<string>();
        let categories = new Array<Category>();

        for (let schema of schemas) {
            if (categoryNames.indexOf(schema.category) === -1) {
                index++;
                categoryNames.push(schema.category);
                categories.push(new Category(index, schema.category, schema.description));
            }
        }
        return categories;
    }

    private getCategoryOptions(schemas: EventSchema[]): CategoryOption[] {
        let index = 0;
        let options = new Array<CategoryOption>();

        let desc = "";
        for (let schema of schemas) {
            let categoryIndex = this.categories.findIndex((element) => {
                return element.name === schema.category;
            });
            categoryIndex++;

            if (schema.category === "Timer") {
                for (let field of schema.fields) {
                    index++;
                    desc = "Schedule \"" + field.name +
                        "\" - Executes the transformation every hour, starting from the time when schedule. ";
                    options.push(new CategoryOption(index, categoryIndex, field.name, desc, field.value));
                }
            } else {
                index++;
                options.push(new CategoryOption(index, categoryIndex, schema.deviceKind, schema.deviceKind, ""));
            }
        }

        return options;
    }

    private getCategoryOptionFields(schemas: EventSchema[]): CategoryOptionField[] {
        let index = 0;
        let optionFields = new Array<CategoryOptionField>();
        let devSchemas = schemas.filter((item) => item.category === "Device");

        for (let schm of devSchemas) {
            index++;
            let i = 0;
            if (schm.fields) {
                for (let field of schm.fields) {
                    i++;
                    let desc = field.name +
                        "- Executes the transformation when my device generates a " + field.name + " event. "
                    let category = this.categoryOptions.filter((item) => item.name === schm.deviceKind);
                    optionFields.push(new CategoryOptionField(i, category[0].id, field.name, desc, field.value));
                }
            }
        }

        return optionFields;
    }
}
