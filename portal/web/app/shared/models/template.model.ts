/**
 * Copyright 2016 Verizon Laboratories. All rights reserved.
 * See provided LICENSE file for use of this source code.
 */

export class Template {
    public static UnknownType: string = "unknown";
    public static CellType: string = "cell";
    public static NotebookType: string = "notebook";

    private id: string;
    private name: string;
    private createdOn: string;
    private type: string;
    private data: any;

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as Template);
        }
    }

    // Returns the template name.
    public getName(): string {
        return this.name;
    }

    // Returns the template type.
    public getType(): string {
        return this.type;
    }

    // Returns the template data.
    public getData(): any {
        return this.data;
    }

    // Helper method used to copy templates.
    private unmarshal(obj: Template): void {
        this.id = obj.id;
        this.name = obj.name;
        this.createdOn = obj.createdOn;
        this.type = obj.type;
        this.data = obj.data;
    }
}
