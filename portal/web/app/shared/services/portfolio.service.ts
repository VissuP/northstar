import {Injectable} from "@angular/core";
import {Http, Response} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {LoggingService} from "./logging.service";
import {Portfolio} from "../models/portfolio.model";
import {File} from "../models/file.model";
import {NsError} from "../models/error.model";

// Defines bucket url.
const bucketUrl: string = "/ns/v1/portfolios/";

@Injectable()
export class PortfolioService {
    public portfolio: Portfolio;
    private logs: LoggingService;
    private http: Http;

    constructor(http: Http, logs: LoggingService) {
        this.logs = logs;
        this.http = http;
    }

    // getPortfolios returns list of portfolios.
    public getPortfolios(): Observable<Portfolio[]> {
        return this.http.get(bucketUrl)
            .map(this.extractPortfolioList)
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }

    // Return portfolio with files list.
    public getFiles(name: string): Observable<Portfolio> {
        return this.http.get(bucketUrl + name)
            .map(res => this.extractFilesList(res, name))
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }
 
    // Helper method used to extract collection portfolio objects from a response.
    private extractPortfolioList(res: Response): Portfolio[] {
        let body: Portfolio[];
        if (res.json()) {
            body = new Array<Portfolio>();
            for (let portfolio of res.json()) {
                body.push(new Portfolio(portfolio));
            }
        }

        return body;
    }

    // Helper method used to extract files list from a response.
    private extractFilesList(res: Response, name: string): Portfolio {
        let body: Portfolio;
        if (res.json()) {
            body = new Portfolio();
            body.name = name;
            for (let file of res.json()) {
                body.addFile(file);
            }
        }
        return body;
    }
}