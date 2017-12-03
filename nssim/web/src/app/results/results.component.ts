import {Component, Input, ViewChild} from "@angular/core";
import {NSSimService} from "../nssim.service";
import {ExecutionResult, Summary} from "../summary.model";

@Component({
  selector: 'results-list',
  templateUrl: './results.component.html',
  styleUrls: ['./results.component.css']
})
export class ResultsComponent {
  @ViewChild('testTable') table: any;
  private summary: Summary;
  private simService: NSSimService;
  @Input() modeAuto: boolean;

  constructor(simService: NSSimService) {
    this.simService = simService;
    simService.getTestSummary().subscribe((response) => {
      this.summary = response;
      console.log('SUMMARY:', this.summary);
    });
  }

  public toggleExpandLogs(row) {
    this.table.rowDetail.toggleExpandRow(row);
  }

  public selectLogs(row, result: ExecutionResult) {
    row.result = result.logMessages;
  }

  public executeTest(testName: string) {
    this.simService.executeTest(testName).subscribe((response) => {
      console.log("RESPONSE:", response);
    });

  }

}
