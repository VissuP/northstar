import * as showdown from "showdown";
import * as katex from "katex";

export const RenderType = {
    Latex: "latex",
    Markdown: "markdown",
    Html: "html",
};

export function getRenderer(type: string) {
    switch (type) {
        case RenderType.Latex:
            return new LatexRenderer();
        case RenderType.Markdown:
            return new MarkdownRenderer();
        case RenderType.Html:
            return new HtmlRenderer();
        default:
            console.log("Rendeerer not found.");
            return new BaseRenderer();
    }
}

export class BaseRenderer {
    constructor() {
    }

    public Render(code: string): string {
        return "Not found";
    };
}

export class LatexRenderer extends BaseRenderer {
    constructor() {
        super();
    }

    public Render(code: string): string {
        let html: string;
        try {
            html = katex.renderToString(code,
                {
                    displayMode: false,
                    throwOnError: false,
                });
        } catch(err) {
            html = err.message;
        }
        return html;
    };
}

export class MarkdownRenderer extends BaseRenderer {
    private showdown;

    constructor() {
        super();
        this.showdown = new showdown.Converter();
        this.showdown.setOption("parseImgDimensions", true);
        this.showdown.setOption("simplifiedAutoLink", true);
        this.showdown.setOption("strikethrough", true);
        this.showdown.setOption("tables", true);

    }

    public Render(code: string): string {
        return this.showdown.makeHtml(code);
    };
}

export class HtmlRenderer extends BaseRenderer {
    constructor() {
        super();
    }

    public Render(code: string): string {
        return code;
    };
}

