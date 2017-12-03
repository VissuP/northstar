import {Component, Input, OnInit, Output, EventEmitter} from "@angular/core";
import {BaseRenderer, getRenderer} from "./render.interface";
import {DomSanitizer, SafeResourceUrl} from "@angular/platform-browser";
import {LanguageConfig} from "../../../../../shared/services/language.service";
import {Cell} from "ngx-vz-cell";

@Component({
    selector: "static-cell",
    styleUrls: ["static.component.css"],
    templateUrl: "static.component.html",
})

export class StaticCellComponent implements OnInit {
    @Input() cell: Cell;
    @Input() id: string;
    @Output() valueChange: EventEmitter<string> = new EventEmitter<string>();
    @Input() activeCell: boolean;
    @Input() readOnly: boolean;
    @Input() config: LanguageConfig;

    private editing: boolean;
    private renderer: BaseRenderer;
    private htmlSrc: SafeResourceUrl;
    private sanitizer: DomSanitizer;

    constructor(sanitizer: DomSanitizer) {
        this.sanitizer = sanitizer;
    }

    public ngOnInit() {
        this.renderer = getRenderer(this.cell.language);

        // Editing selects the output or editor. If there is no content, editing has to
        // be enabled (no content can't be clicked).
        this.editing = false;
        if (this.cell.code.length === 0) {
            this.editing = true;
        } else {
            let html = this.renderer.Render(this.cell.code);
            this.htmlSrc = this.sanitizer.bypassSecurityTrustHtml(html);
        }
    }

    public startEdit(): void {
        if (!this.readOnly) {
            this.editing = true;
        }
    }

    public stopEdit(): void {
        if (this.cell.code.trim().length > 0) {
            this.editing = false;
        }

        // TODO: Consider the ramifications of this. See
        // https://github.com/showdownjs/showdown/wiki/Markdown's-XSS-Vulnerability-(and-how-to-mitigate-it)
        let html = this.renderer.Render(this.cell.code);
        this.htmlSrc = this.sanitizer.bypassSecurityTrustHtml(html);
    }
}
