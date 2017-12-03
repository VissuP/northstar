/**
 * Copyright 2016 Verizon Laboratories. All rights reserved.
 * See provided LICENSE file for use of this source code.
 */
import { Injectable } from "@angular/core";
import {Http, Response} from "@angular/http";
import {NsError} from "../models/error.model";
import {Template} from "../models/template.model";
import {Observable} from "rxjs";

// Defines the templates API url.
const templateUrl: string = "/ns/v1/templates";

@Injectable()
export class TemplateService {
    private http: Http;

    constructor(http: Http) {
        this.http = http;
    }

    // Returns templates.
    public getTemplates(): Observable<Template[]> {
        return this.http.get(templateUrl)
            .map(this.extractTemplates)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Helper method used to extract collection Notebook objects from a response.
    private extractTemplates(response: Response): Template[] {
        let templates = new Array<Template>();
        let json = response.json();

        if (json) {
            for (let template of json) {
                templates.push(new Template(template));
            }
        }

        return templates;
    }
}