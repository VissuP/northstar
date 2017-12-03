import { Injectable } from "@angular/core";
import { Http, Response } from "@angular/http";
import { Observable } from "rxjs/Observable";
import { LoggingService } from "./logging.service";
import { Notebook, NotebookCell } from "../models/notebook.model";
import { User } from "../models/user.model";
import { NsError } from "../models/error.model";
import { Cell } from "ngx-vz-cell";

// Defines notebook url.
const notebooksUrl: string = "/ns/v1/notebooks/";

@Injectable()
export class NotebookService {
    public notebook: Notebook;
    private logs: LoggingService;
    private http: Http;

    constructor(http: Http, logs: LoggingService) {
        this.logs = logs;
        this.http = http;
    }

    // Return notebook with specified id.
    public getNotebook(id: string): Observable<Notebook> {
        return this.http.get(notebooksUrl + id)
            .map(this.extractNotebook)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // getNotebooks returns list of notebooks.
    public getNotebooks(): Observable<Notebook[]> {
        return this.http.get(notebooksUrl)
            .map(this.extractNotebookList)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Creates a new notebook.
    public createNotebook(notebook: Notebook): Observable<Notebook> {
        return this.http.post(notebooksUrl, JSON.stringify(notebook))
            .map(this.extractNotebook)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Updates an existing notebook. Returns true if notebook updated successfully, false otherwise.
    public updateNotebook(notebook: Notebook): Observable<boolean> {
        return this.http.put(notebooksUrl, JSON.stringify(notebook))
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Deletes notebook with id.
    public deleteNotebook(id: string): Observable<boolean> {
        return this.http.delete(notebooksUrl + id)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Returns the users associated with notebook id.
    public getNotebookUsers(notebookId: string): Observable<User[]> {
        return this.http.get(notebooksUrl + notebookId + "/users")
            .map(this.extractUserList)
            .catch((error: Response) => {
                throw new NsError(error);
            });
        ;
    }

    // Updates the users associated with notebook id.
    public setNotebookUsers(notebookId: string, users: User[]): Observable<boolean> {
        return this.http.put(notebooksUrl + notebookId + "/users", JSON.stringify(users))
            .catch((error: Response) => {
                throw new NsError(error);
            });
        ;
    }

    // Helper method used to extract notebook gets the notebook from the response
    private extractNotebook(res: Response) {
        if (res.json()) {
            return new Notebook(res.json());
        }
        return null;
    }

    // Helper method used to extract collection Notebook objects from a response.
    private extractNotebookList(res: Response): Notebook[] {
        let body: Notebook[];
        if (res.json()) {
            body = new Array<Notebook>();
            for (let notebook of res.json()) {
                body.push(new Notebook(notebook));
            }
        }

        return body;
    }

    // Helper method used to extract collection of notebook users from a response.
    private extractUserList(res: Response): User[] {
        let body = new Array<User>();
        if (res.text()) {
            for (let user of res.json()) {
                body.push(new User(user));
            }
        }

        return body;
    }

    // Helper method used to return a new notebook with base64 decoded cell content, if necessary
    public decodeNotebook(notebook: Notebook) {
        let decodedNotebook: Notebook = new Notebook(notebook);
        decodedNotebook.cells = new Array<NotebookCell>();
        notebook.cells.forEach(notebookCell => {
            let decodedCell = new Cell(notebookCell);
            if (this.isBase64(notebookCell.code)) { // decode only if code is base64 encoded
                decodedCell.code = atob(notebookCell.code);
            }
            decodedNotebook.addCell(decodedCell);
        });
        return decodedNotebook;
    }

    // Helper method used to return a new notebook with base64 encoded cell content
    public encodeNotebook(notebook: Notebook) {
        let encodedNotebook: Notebook = new Notebook(notebook);
        encodedNotebook.cells = new Array<NotebookCell>();
        notebook.cells.forEach(notebookCell => {
            let encodedCell = new Cell(notebookCell);
            encodedCell.code = btoa(notebookCell.code);
            encodedNotebook.addCell(encodedCell);
        });
        return encodedNotebook;
    }

    // Helper method to dtermine is input string is a valid base64 string
    // https://stackoverflow.com/questions/8571501/how-to-check-whether-the-string-is-base64-encoded-or-not
    public isBase64(input: string): boolean {
        let base64RegEx = /^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$/i;
        return base64RegEx.test(input);
    }


}