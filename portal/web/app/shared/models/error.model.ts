import {Response} from "@angular/http";

export class NsError {
    private error: string;
    private description: string;

    constructor(obj?: Response) {
       if (obj) {
           try {
               let body = obj.json();
               this.error = (body as any).error;
               this.description = (body as any).error_description;
           } catch (error) {
               this.description = obj +"";
           }
       }
    }

    public Get(): string {
        return this.description;
    }
}
