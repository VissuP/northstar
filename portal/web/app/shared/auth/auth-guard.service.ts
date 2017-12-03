import {Injectable} from '@angular/core';
import {CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot, CanActivateChild} from '@angular/router';
import {Http, Response, Headers} from '@angular/http';
import {Observable} from 'rxjs/Observable';


@Injectable()
export class AuthGuard implements CanActivate, CanActivateChild {

    // TODO - We should not have this url in here. Maybe we should create
    // a service that enables us to store user selections in local storage.
    public redirectUrl: string = '';
    private http: Http;
    private router: Router;

    constructor(router: Router, http: Http) {
        this.http = http;
        this.router  = router;
    }

    // CanActivate is used for auth guard on parent routes (i.e. /portal).
    public canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
        let url: string = state.url;
        return this.http.post('/users/actions/verify', null)
                        .map(
                            response => {
                              return response.status == 204; 
                                    }
                                        
                            ).catch((error: any) => {
                                this.redirectUrl = url;
                                this.router.navigate(['./login']);
                                return Observable.of(false)
                                                    }
                                    )
                                                                                                        }

    // CanActivateChild is used for auth guard on child routes (i.e. /portal/dashboard)
    public canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
        return this.canActivate(route, state);
    }

    public login(formValues: Object) {
        // Send request to login the user.
        let headers = new Headers();
        headers.append("Content-Type", "application/json");

        return this.http.post(
            '/users/actions/login',
            JSON.stringify({email: (formValues as any).username, password: ( formValues as any).password}),
            {headers: headers})
            .map(
            response => {
                if (response.status == 204) {
                    return response.statusText;
                }
                return null
            }
        );
    }

    public logout() {
        // TODO: Add return type.
        return this.http.post("/users/actions/logout",null);
    }
}