import {Injectable} from "@angular/core";
import {Http} from "@angular/http";
import "rxjs/Rx";
import {Summary} from "./summary.model";

const simapiURL = '/sim/v1/';

@Injectable()
export class NSSimService {
  private http: Http;

  constructor(http: Http) {
    this.http = http;
    console.log(this.http);
  }

  public getInfo() {
    return this.http.get(simapiURL + 'info')
      .map((response) => {
        return response;
      });
  }

  public getTestSummary() {
    return this.http.get(simapiURL + 'tests')
      .map((response) => {
        return new Summary(response.json());
      });
  }

  public executeTest(name: string) {
    return this.http.post(simapiURL + 'tests/name/execute', null)
      .map((response) => {
        return response;
      });
  }
}
