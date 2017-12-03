export class File { 

    public name: string;
    public lastModified: string;
    public size: number;
    public etag: string;

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as File);
        } 
    };

    protected unmarshal(obj: File)  {
        this.name = obj.name;
        this.lastModified = obj.lastModified;
        this.size = obj.size;
        this.etag = obj.etag;
    }
}