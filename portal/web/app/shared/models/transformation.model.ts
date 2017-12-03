const NoneTransformationType: string = "None";

// Defines a transformation model object.
export class Transformation {
    public static Scheduled: string = "Scheduled";
    public static Unscheduled: string = "Unscheduled";
    public static DefaultEntryPoint: string = "main";
    public static DefaultLanguage: string = "lua";

    public code: Code;
    public description: string;
    public entryPoint: string;
    public id: string;
    public language: string;
    public lastUpdated: string;
    public name: string;
    public scheduled: boolean;
    public state: string;
    public timeout: number;
    public version: string;
    public arguments: Object;
    public schedule: Schedule;

    // JSON decoding in typescript only works if you set the fields manually otherwise it will be a generic object.
    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as Transformation);
        } else {
            this.code = new Code();
            this.entryPoint = Transformation.DefaultEntryPoint;
            this.language = Transformation.DefaultLanguage;
            this.schedule = new Schedule();
            this.timeout = 10;
        }
    }

    // Encode makes a copy and packages up a transformation so that it can be sent over the backend.
    public encode(): Transformation {
        let encodedTransformation: Transformation = new Transformation();

        encodedTransformation.arguments = this.arguments;
        encodedTransformation.code = new Code(this.code);
        encodedTransformation.description = this.description;
        encodedTransformation.entryPoint = this.entryPoint;
        encodedTransformation.id = this.id;
        encodedTransformation.language = this.language;
        encodedTransformation.lastUpdated = this.lastUpdated;
        encodedTransformation.name = this.name;
        encodedTransformation.schedule = new Schedule(this.schedule);
        encodedTransformation.scheduled = this.scheduled;
        encodedTransformation.setState();
        encodedTransformation.timeout = this.timeout;
        encodedTransformation.version = this.version;

        // Cron service requires code to be base64 encoded.
        encodedTransformation.code.value = btoa(this.code.value);
        return encodedTransformation;
    }

    public setState(): void {
        this.state = this.isScheduled() ? Transformation.Scheduled : Transformation.Unscheduled;
    }

    public setScheduled(): void {
        this.scheduled = this.isScheduled();
    }

    public getScheduleDescription(): string {
        if (this.isScheduled()) {
            return "This transformation operates on " + this.schedule.event.name.toLocaleLowerCase();
        }
        return "This transformation is not scheduled.";
    }

    private getDescription(): string {
        if (this.description) {
            return this.description;
        }
        return "No description set.";
    }

    private unmarshal(obj: Transformation) {
        this.arguments = obj.arguments;
        this.code = new Code(obj.code);
        this.decode(obj.code.value);
        this.description = obj.description;
        this.entryPoint = obj.entryPoint;
        this.id = obj.id;
        this.language = obj.language;
        this.lastUpdated = obj.lastUpdated;
        this.name = obj.name;
        this.schedule = new Schedule(obj.schedule);
        this.scheduled = obj.scheduled;
        this.setState();
        this.timeout = obj.timeout;
        this.version = obj.version;

    }

    // Decode does the required processing to receive a transformation from the backend.
    // Note that this is called automatically by the constructor.
    private decode(code: string): void {
        this.code.value = atob(code);
    }

    // Helper method used to determine if transformation is scheduled.
    private isScheduled(): boolean {
        return this.schedule && this.schedule.event && this.schedule.event.type !== NoneTransformationType;
    }
}

export class Code {
    public type: string;
    public value: string;

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as Code);
        }
    }

    private unmarshal(obj: Code) {
        this.type = obj.type;
        this.value = obj.value;
    }
}

export class EventSchema {
    public category: string;
    public deviceKind: string;
    public name: string;
    public description: string;
    public type: string;
    public semantic: string;
    public fields: EventSchemaField[];

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as EventSchema);
        } else {
            this.fields = new Array<EventSchemaField>();
            this.category = "None";
        }
    };

    private unmarshal(obj: EventSchema) {
        this.category = obj.category;
        this.deviceKind = obj.deviceKind;
        this.name = obj.name;
        this.description = obj.description;
        this.type = obj.type;
        this.semantic = obj.semantic;
        this.fields = obj.fields;
    }
}

export class EventSchemaField {
    public name: string;
    public type: string;
    public constant: boolean;
    public value: string;

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as EventSchemaField);
        }
    };

    public getName(): string {
        return this.name;
    }

    public getType(): string {
        return this.type;
    }

    public getValue(): string {
        return this.value;
    }

    private unmarshal(obj: EventSchemaField) {
        this.name = obj.name;
        this.type = obj.type;
        this.value = obj.value;
        this.constant = obj.constant;
    }
}

export class Schedule {
    public kind: string;
    public id: string;
    public version: string;
    public createdOn: string;
    public lastUpdated: string;
    public event: TransformationEvent;

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as Schedule);
        } else {
            this.event = new TransformationEvent();
        }
    }

    private unmarshal(obj: Schedule) {
        this.kind = obj.kind;
        this.id = obj.id;
        this.version = obj.version;
        this.createdOn = obj.createdOn;
        this.lastUpdated = obj.createdOn;
        this.event = new TransformationEvent(obj.event);
    }
}

export class TransformationEvent {
    public type: string;
    public name: string;
    public value: string;

    constructor(obj?: Object) {
        if (obj) {
            this.unmarshal(obj as TransformationEvent);
        } else {
            this.type = NoneTransformationType;
        }
    };

    private unmarshal(obj: TransformationEvent) {
        this.type = obj.type;
        this.name = obj.name;
        this.value = obj.value;
    }
}
