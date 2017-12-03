import {NotebookPermissions} from "./permissions.model";

// Defines the type that represent a thingspace user.
export class User {
    public displayName: string;
    public email: string;
    public accountId: string;
    public id: string;
    public imageId: string;
    public permissions: string;

    constructor(obj?: Object) {
        if (obj) {
            this.displayName = (obj as any).displayName;
            this.email = (obj as any).email;
            this.accountId = (obj as any).accountId;
            this.id = (obj as any).id;
            this.imageId = (obj as any).imageId;
            this.permissions = (obj as any).permissions;
        } else {
            this.permissions = NotebookPermissions.Read;
        }
    }

    // Returns true is user is resource (e.g., notebook) owner,
    // false otherwise.
    public isOwner() {
        return this.permissions === NotebookPermissions.Owner;
    }
}
