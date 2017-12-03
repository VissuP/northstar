import {Injectable} from "@angular/core";
import {ActivatedRouteSnapshot, Resolve} from "@angular/router";
import {Observable} from "rxjs/Observable";
import {ExecutionOutput} from "ngx-vz-cell";
import {TransformationService} from "../transformation.service";

@Injectable()
export class TransformationResultsResolver implements Resolve<ExecutionOutput[]> {
    private transformationService: TransformationService;

    constructor(transformationService: TransformationService) {
        this.transformationService = transformationService;
    }

    // Resolve allows us to load major resources before attempting the page
    // itself and catch any errors before navigating.
    public resolve (route: ActivatedRouteSnapshot): Observable<ExecutionOutput[]> {
        return this.transformationService.getTransformationResults(route.params.id);
    }
}
