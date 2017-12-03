import { Injectable } from "@angular/core";
import { ActivatedRouteSnapshot, Resolve } from "@angular/router";
import { Observable } from "rxjs/Observable";

import { Notebook } from "../../../shared/models/notebook.model";
import { NotebookService } from "../../../shared/services/notebook.service";

@Injectable()
export class NotebookResolver implements Resolve<Notebook> {
    private notebookService: NotebookService;
    constructor(notebookService: NotebookService) {
        this.notebookService = notebookService;
    }

    // Resolve allows us to load major resources before attempting the page
    // itself and catch any errors before navigating.
    public resolve(route: ActivatedRouteSnapshot): Observable<Notebook> {
        return this.notebookService.getNotebook(route.params.id)
            .map(notebook => this.notebookService.decodeNotebook(notebook));
    }
}
