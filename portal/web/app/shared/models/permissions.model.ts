// Defines the type that represents notebook permissions.
export const NotebookPermissions = {
    Read: "Read",
    ReadExecute: "ReadExecute",
    ReadWrite: "ReadWrite",
    ReadWriteExecute: "ReadWriteExecute",
    Owner: "Owner",
};

// Helper method used to create array of notebook permissions.
export function PermissionsArray(): Array<string> {
    let permissions = Object.keys(NotebookPermissions);
    return permissions;
}
