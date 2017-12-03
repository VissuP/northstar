import { Http, Response, Headers, RequestOptions } from "@angular/http";
import {Subscription} from "rxjs/Subscription";
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";

import {CompleterBaseData} from "ng2-completer";
import {User} from "../../../../shared/models/user.model";

export function getRemotePostData(url:string, searchFields:string, titleField:string, http:Http):RemotePostData {
    return new RemotePostData(http).remoteUrl(url).searchFieldss(searchFields).titleField(titleField);
};

export class RemotePostData extends CompleterBaseData {
    private _remoteUrl: string;
    private remoteSearch: Subscription;
    private _dataField: string = null;
    private _headers: Headers;
    private _queryFormatter: (args:any) => any = null;
    private users: User[];
    private http: Http;

    constructor(http:Http) {
        super();
        this.http = http;
    }

    public remoteUrl(remoteUrl: string) {
        this._remoteUrl = remoteUrl;
        return this;
    }

    // queryFormatter expects a function to use to generate a query object.
    public queryFormatter(formatter: (arg: any) => any): void {
        this._queryFormatter = formatter;
    }

    public dataField(dataField: string): void {
        this._dataField = dataField;
    }

    public headers(headers: Headers): void {
        this._headers = headers;
    }

    public search(term: string): void {
        this.cancel();

        let url = this._remoteUrl;

        let params = {};
        if (this._queryFormatter) {
            params = this._queryFormatter(term);
        }

        let options = new RequestOptions({
            headers: this._headers || new Headers(),
        });

        this.remoteSearch = this.http.post(url, params, options)
            .map((res: Response) => {
                // store a copy of the user collection for later lookup of the specific object
                this.users = new Array<User>();

                for (let user of res.json()) {
                    let u = new User(user);

                    // Note that this is to make sure the UI does not show empty
                    // user information. E.g., empty display name.
                    if (u.displayName === "") {
                        u.displayName = u.email;
                    }

                    this.users.push(u);
                }

                return this.users;
            })
            .map((data: any) => {
                let matches = this.extractValue(data, this._dataField);
                return this.extractMatches(matches, term);
            })
            .map((matches: any[]) => {
                let results = this.processResults(matches, term);
                this.next(results);
                return results;
            })
            .catch((err) => {
                this.error(err);
                return null;
            })
            .subscribe();
    }

    public lookupUser(displayName: string): User {
        for (let user of this.users) {
            if (user.displayName === displayName) {
                return user;
            }
        }

        return null;
    }

    public cancel(): void {
        if (this.remoteSearch) {
            this.remoteSearch.unsubscribe();
        }
    }
}
