export class EditorConfig {
    public autoCloseBackets: boolean;
    public foldGutter: boolean;
    public gutters: string[];
    public lineNumbers: boolean;
    public lineWrapping: boolean;
    public matchBrackets: boolean;
    public mode: any;
    public readOnly: boolean;
    public styleActiveLine: boolean;
    public theme: string;
    public viewPortMargin: number;

    constructor() {
    }
}

export const LanguageMode = {
    Code: "Code",
    Query: "Query",
    Static: "Static",
    Notation: "Notation",
};

export class LanguageConfig {
    public icon: string;
    public type: string;
    public language: string;
    public editorConfig: EditorConfig;

    public isCode() {
        return this.type === LanguageMode.Code;
    }

    public isStatic() {
        return this.type === LanguageMode.Static;
    }

    public isExecutable() {
        return this.type == LanguageMode.Code;
    }
}

const LanguageIconBaseUri = "assets/img/languages/";
export class LanguageService {
    public configs: Map<string, LanguageConfig>;

    constructor() {
        this.configs = new Map<string, LanguageConfig>();
        this.configs.set("lua", this.generateConfig("lua"));
        this.configs.set("html", this.generateConfig("html"));
        this.configs.set("latex", this.generateConfig("latex"));
        this.configs.set("markdown", this.generateConfig("markdown"));
        this.configs.set("json", this.generateConfig("json"));
    }

    public languageArray(...typeFilter: string[]): Array<LanguageConfig> {
        let languages = new Array<LanguageConfig>();

        this.configs.forEach(
            (language: LanguageConfig) => {
                if (typeFilter.length === 0 || typeFilter.indexOf(language.type) !== -1) {
                    languages.push(language);
                }
            });

        return languages.sort(
            (language1, language2): number => {
                if (language1.language > language2.language) {
                    return 1;
                }
                if (language1.language < language2.language) {
                    return -1;
                }
                return 0;
            });
    }

    // TODO: Convert the following to pull specific language settings from portal configuration
    private getCodemirrorPreferences(): EditorConfig {
        let codeMirrorConfig = new EditorConfig();
        codeMirrorConfig.autoCloseBackets = true;
        codeMirrorConfig.foldGutter = true;
        codeMirrorConfig.lineNumbers = true;
        codeMirrorConfig.lineWrapping = true;
        codeMirrorConfig.styleActiveLine = true;
        codeMirrorConfig.gutters = ["CodeMirror-linenumbers", "CodeMirror-foldgutter"];
        codeMirrorConfig.theme = "solarized light";
        codeMirrorConfig.styleActiveLine = true;
        codeMirrorConfig.matchBrackets = true;
        codeMirrorConfig.viewPortMargin = Infinity;
        return codeMirrorConfig;
    }

    private generateConfig(language: string): LanguageConfig {
        let config = new LanguageConfig();
        config.language = language;

        switch (language) {
            case "lua":
                config.type = LanguageMode.Code;
                config.editorConfig = this.getCodemirrorPreferences();
                config.editorConfig.mode = {
                    name: "lua",
                };
                config.icon = LanguageIconBaseUri + "lua.png";
                break;
            case "r":
                // Note that "R" is not supported for now but keeping the configuration.
                config.type = LanguageMode.Code;
                config.editorConfig = this.getCodemirrorPreferences();
                config.editorConfig.mode = {
                    name: "r",
                };
                config.icon = LanguageIconBaseUri + "r.png";
                break;
            case "latex":
                config.type = LanguageMode.Static;
                break;
            //TODO: Find an icon
            case "html":
                config.type = LanguageMode.Static;
                config.icon = LanguageIconBaseUri + "html.png";
                break;
            case "markdown":
                config.type = LanguageMode.Static;
                break;
            //TODO: Find an icon
            case "json":
                config.type = LanguageMode.Notation;
                config.editorConfig = this.getCodemirrorPreferences();
                config.editorConfig.mode = "application/ld+json";
            default:
        }

        return config;
    };
}