import { Injectable } from '@angular/core';
import {Http, Response, Headers} from '@angular/http';
import {NsError} from "../models/error.model";

// Defines the users API url.
const userUrl:string = "/ns/v1/user";

@Injectable()
export class UserService {
    private http: Http;
    constructor(http: Http) {
        this.http = http;
    }

    // Returns information about the authenticated user.
    public getUser() {
        // Note that this will return an observable.
        return this.http.get(userUrl)
            .map((user: Response) => {
                return user.json();
            })
            .catch((error: Response) => {
                throw new NsError(error);
            });
    }
}