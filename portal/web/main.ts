// main entry point
import {platformBrowserDynamic} from '@angular/platform-browser-dynamic';
import {AppModule} from './app/app.module';
import {enableProdMode} from '@angular/core';
import {environment} from './app/environments/environment';

if(environment.production) {
    enableProdMode();
}
platformBrowserDynamic().bootstrapModule(AppModule);
