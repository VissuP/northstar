import {ErrorHandler, Injectable} from "@angular/core";
import {LoggingService} from "../services/logging.service";

@Injectable()
export class GlobalErrorHandler implements ErrorHandler {
    private log: LoggingService;

    constructor(log: LoggingService) {
        this.log = log;
    }

    public handleError(error: any): void {
        // Note that we're at the bottom, many angular components depend on us which means we can't depend on them.
        this.log.error("Recovered from uncaught exception: ", error);
        window.location.href = "/northstar/#/portal/error";
    }
}

// These are the options to override the root error handler that angular uses for uncaught exceptions.
export let GlobalErrorHandlerProvider = [
    {
        provide: ErrorHandler,
        useClass: GlobalErrorHandler,
    }
];