import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import {NSSimService} from './nssim.service'
import {NgxDatatableModule} from '@swimlane/ngx-datatable';
import {ResultsComponent} from './results/results.component';

@NgModule({
  declarations: [
    AppComponent,
    ResultsComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    NgxDatatableModule
  ],
  providers: [
    NSSimService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
