// Define supported alert types.
export const AlertType = {
    Error: "error",
    Success: "success",
    Warning: "warning",
};

// Defines the supported time units.
export const Time = {
    Second: 1000,
    Minute: 60 * 1000,
    Hour: 60 * 60 * 1000,
    Day: 24 * 60 * 60 * 1000,
};

// Defines alert type.
export class Alert {
    public timeout: number;    // Time at which the alert banner expires.
    public message: string;
    public href: string;
    public type: string;

    constructor(type: string, message: string, timeout?: number, href?: string) {
        this.type = type;
        this.message = message;

        if (!timeout) {
            timeout = 10 * Time.Second; // Arbitrary default. May be a better value.ß
        }
        this.timeout = timeout;

        this.href = href;
    }
}
