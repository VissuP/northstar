import './rxjs-extensions';

import {NgModule}      from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {FormsModule, ReactiveFormsModule}   from '@angular/forms';
import {HttpModule}    from '@angular/http';

import {AppRoutingModule} from './app.routing';
import {AppComponent} from './app.component';
import {LoginComponent} from './login/login.component';
import {AuthGuard} from './shared/auth/auth-guard.service';
import {PortalModule} from './portal/portal.module';

import {LoggingService, LogLevel} from "./shared/services/logging.service";
import {environment} from "./environments/environment";
import {LocationStrategy, HashLocationStrategy} from '@angular/common';


@NgModule({
    imports: [
        BrowserModule,
        FormsModule,
        HttpModule,
        ReactiveFormsModule,
        PortalModule,
        AppRoutingModule
    ],
    declarations: [
        AppComponent,
        LoginComponent,
    ],
    providers: [
        AuthGuard,
        LoggingService,
        { provide: LocationStrategy, useClass: HashLocationStrategy }
    ],
    bootstrap: [AppComponent]

})

export class AppModule {
    constructor(log: LoggingService) {
        if (!environment.production) {
            log.setLogLevel(LogLevel.Trace);
        }
    }
}
