import {Injectable} from "@angular/core";
import {Http, Response} from "@angular/http";
import {Observable} from "rxjs/Observable";
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";
import {EventSchema, Transformation} from "../../shared/models/transformation.model";
import {LoggingService} from "../../shared/services/logging.service";
import {ExecutionOutput} from "ngx-vz-cell";
import {NsError} from "../../shared/models/error.model";

const transformationUrl: string = "/ns/v1/transformations/";
const eventSchemasUrl: string = "/ns/v1/events/schemas/";

@Injectable()
export class TransformationService {
    private http: Http;
    private log: LoggingService;

    constructor(http: Http, log: LoggingService) {
        this.http = http;
        this.log = log;
    }

    // Returns collection of transformations.
    public getTransformations(): Observable<Transformation[]> {
        return this.http.get(transformationUrl)
            .map(this.extractTransformationList)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Returns transformation execution results for a given id.
    public getTransformationResults(id: string): Observable<ExecutionOutput[]> {
        return this.http.get(transformationUrl + id + "/results")
            .map(this.extractTransformationResults)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Creates a new transformation.
    public createTransformation(transformation: Transformation): Observable<Transformation> {
        return this.http.post(transformationUrl, transformation.encode())
            .map(this.extractTransformation)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    public updateTransformation(transformation: Transformation): Observable<Transformation> {
        return this.http
            .put(transformationUrl, transformation.encode())
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Deletes transformation with id.
    public deleteTransformation(id: string): Observable<boolean> {
        return this.http.delete(transformationUrl + id)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Returns collection of transformations.
    public getEventSchemas(): Observable<EventSchema[]> {
        return this.http.get(eventSchemasUrl)
            .map(this.extractEventSchemaList)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Extract transformation gets the transformation from the response
    private extractTransformation(res: Response) {
        if (res.json()) {
            return new Transformation(res.json());
        }
        return null;
    }

    // Extract transformation list gets the list of transformations from the response
    private extractTransformationList(res: Response): Transformation[] {
        let body: Transformation[];
        if (res.json()) {
            body = new Array<Transformation>();
            for (let transformationEntry of res.json()) {
                // Code is stored base64 encoded. Decode it here.
                let transformation = new Transformation(transformationEntry);

                body.push(transformation);
            }
        }

        return body || [];
    }

    // Extract transformation list gets the list of transformations from the response
    private extractEventSchemaList(res: Response): EventSchema[] {
        let body: EventSchema[];
        if (res.json()) {
            body = new Array<EventSchema>();

            for (let eventSchema of res.json()) {
                body.push(new EventSchema(eventSchema));
            }
        }

        return body || [];
    }

    private extractTransformationResults(res: Response): ExecutionOutput[] {
        let body: ExecutionOutput[];
        if (res.json()) {
            body = new Array<ExecutionOutput>();
            for (let execution of res.json()) {
                body.push(new ExecutionOutput(execution));
            }
        }

        return body;
    }
}
