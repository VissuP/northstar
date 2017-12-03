import { Component, ViewEncapsulation, ViewChild } from '@angular/core';
import {NSSimService} from './nssim.service';
import {Title} from "@angular/platform-browser";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  private title: string;
  private simService: NSSimService;
  private info;
  private modeAuto: boolean;
  @ViewChild('testTable') table: any;

  constructor(simService: NSSimService, title: Title) {
    this.simService = simService;
    simService.getInfo().subscribe((response) => {
        this.info = response.json();
        this.modeAuto = (this.info.mode === 'auto');
        this.title = 'Northstar ' + this.info.environment + ' simulator';

        title.setTitle(this.title);
      },
    );
  }
}
