import {File} from "./file.model";

export class Portfolio { 

    public name: string;
    public creationDate: string;
    public files: File[];

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as Portfolio);
        } else {
            this.name="";
            this.creationDate="";
            this.files = new Array<File>();
        }
    };

    private unmarshal(obj: Portfolio)  {
        this.name = obj.name;
        this.creationDate = obj.creationDate;
        this.files = obj.files;
    }

    // Adds the specified file to the portfolio.
    public addFile(obj?: Object): void {
        let file = new File(obj);
        this.files.push(file);
    }
}