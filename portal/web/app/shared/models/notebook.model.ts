/**
 * Copyright 2016 Verizon Laboratories. All rights reserved.
 * See provided LICENSE file for use of this source code.
 */
import {Cell, ExecutionOutput, StatusIndicator} from "ngx-vz-cell";
import {NotebookPermissions} from "./permissions.model";
import {Utilities} from "../resources/utilities";

export class Notebook {
    public version: string;
    public id: string;
    public name: string;
    public cells: NotebookCell[];
    private permissions: string;

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as Notebook);
        } else {
            this.name = "Untitled Notebook";
            this.cells = new Array<NotebookCell>();
            this.permissions = NotebookPermissions.ReadWriteExecute;
        }
    };

    // Creates a new cell for the specified language.
    public newCell(language: string): void {
        let cell = new NotebookCell();
        cell.setLanguage(language);
        this.cells.push(cell);
    }

    // Adds the specified cell to the notebook.
    public addCell(obj?: Object): void {
        let cell = NotebookCell.clone(obj);
        this.cells.push(cell);
    }

    // Removes the specified cell.
    public removeCell(cell: NotebookCell) {
        let index = this.cells.indexOf(cell, 0);
        if (index > -1) {
            this.cells.splice(index, 1);
        }
    }

    // Returns true if user is notebook owner.
    public isOwner(): boolean {
        return this.permissions === NotebookPermissions.Owner;
    }

    // Returns true if notebook has no cells.
    public isEmpty(): boolean {
        return this.cells.length === 0;
    }

    // Returns true if notebook can be executed by user.
    public canExecute(): boolean {
        return (
            this.permissions === NotebookPermissions.Owner ||
            this.permissions === NotebookPermissions.ReadExecute ||
            this.permissions === NotebookPermissions.ReadWriteExecute
        );
    }

    // Returns true if notebook can be edited by user.
    public canWrite(): boolean {
        return (
            this.permissions === NotebookPermissions.Owner ||
            this.permissions === NotebookPermissions.ReadWrite ||
            this.permissions === NotebookPermissions.ReadWriteExecute
        );
    }

    // Update cells status.
    public setRunningCellsTo(status: string, message: string) {
        for (let cell of this.cells) {
            if (cell.options.status === StatusIndicator.Running) {
                cell.options.status = status;
                cell.output = new ExecutionOutput();
                // Note that this is error condition.
                cell.output.setExecutionError(message);
            }
        }
    }

    private unmarshal(obj: Notebook)  {
        this.version = obj.version;
        this.id = obj.id;
        this.name = obj.name;
        this.permissions = obj.permissions;
        this.cells = new Array<NotebookCell>();
        if (obj.cells) {
            for (let cell of obj.cells) {
                this.cells.push(new NotebookCell(cell));
            }
        }
    }
}

export class NotebookCell extends Cell {

    // Returns a new copy of the cell.
    public static clone(obj: Object): NotebookCell {
        let cell = new NotebookCell(obj);
        cell.id = Utilities.NewGuid();
        return cell;
    }

    public id: string;

    constructor(obj?: Object) {
        super(obj);
        if(obj) {
            super.unmarshal(obj as Cell);
            this.unmarshal(obj as NotebookCell);
        } else {
            this.id = Utilities.NewGuid();
        }
    }

    // Sets the cell language.
    public setLanguage(language: string): void {
        this.language = language;
    }

    protected unmarshal(cell: NotebookCell) {
        this.id = cell.id;
    }
}