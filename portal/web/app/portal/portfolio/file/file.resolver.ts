import {Injectable} from "@angular/core";
import {ActivatedRouteSnapshot, Resolve} from "@angular/router";
import {Observable} from "rxjs/Observable";

import {Portfolio} from "../../../shared/models/portfolio.model";
import {PortfolioService} from "../../../shared/services/portfolio.service";

@Injectable()
export class PortfolioResolver implements Resolve<Portfolio> {
    private portfolioService: PortfolioService;
    constructor(portfolioService: PortfolioService) {
        this.portfolioService = portfolioService;
    }

    // Resolve allows us to load major resources before attempting the page
    // itself and catch any errors before navigating.
    public resolve (route: ActivatedRouteSnapshot): Observable<Portfolio> {
        return this.portfolioService.getFiles(route.params.id);
    }
}
