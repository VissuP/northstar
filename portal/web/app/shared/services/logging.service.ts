import {Injectable} from '@angular/core';

 export const LogLevel = {
    Production: 1,
    Error:      2,
    Warn:       3,
    Debug:      4,
    Trace:      5
};

@Injectable()
export class LoggingService {
    private loglevel: number;

    constructor() {
        this.loglevel = LogLevel.Production;
    }

    public setLogLevel(logLevel: number) {
       this.loglevel = logLevel;
    }

    public error(message?: string, ...optionalParams: any[]): void {
        if (this.loglevel >= LogLevel.Error) {
            console.error(message, ...optionalParams)
        }
    }

    public warn(message?: string, ...optionalParams: any[]): void {
        if (this.loglevel >= LogLevel.Warn) {
            console.warn(message, ...optionalParams)
        }
    }

    public debug(message?: string, ...optionalParams: any[]): void {
        if (this.loglevel >= LogLevel.Debug) {
            console.debug(message, ...optionalParams)
        }
    }

    // Note: This loglevel is special compared to the rest. Not only will it print your message, but it will print a stacktrace.
    public trace(message?: string, ...optionalParams: any[]): void {
        if (this.loglevel >= LogLevel.Trace) {
            console.trace(message, ...optionalParams)
        }
    }

}