
/**
 * Copyright 2016 Verizon Laboratories. All rights reserved.
 * See provided LICENSE file for use of this source code.
 */

export class Summary {
  public startTime: string;
  public environment: string;
  public totalRuns: number;
  public tests: Test[];

  constructor(obj?: Object) {
    this.tests = [];
    if (obj) {
      this.unmarshal(obj as Summary);
    }
  }

  protected unmarshal(obj: Summary) {
    this.startTime = obj.startTime;
    this.environment = obj.environment;
    this.totalRuns = obj.totalRuns;
    for (let test of obj.tests) {
      this.tests.push(new Test(test));
    }
  }
}

export class Test {
  public id: number;
  public name: string;
  public type: string;
  public group: string;
  public state: string;
  public concurrencyIndex: number;
  public concurrency: number;
  public totalExecutions: Number;
  public totalErrors: Number;
  public lastStatus: string;
  public lastLatency: string;
  public lastErrorMessage: string;
  public executionResults: ExecutionResult[];
  public failureStats: Map<string, Stats>;
  public steps: string[];

  constructor(obj?: Object) {
    this.executionResults = [];
    if (obj) {
      this.unmarshal(obj as Test);
    }
  }

  protected unmarshal(obj: Test) {
    this.id = obj.id;
    this.name = obj.name;
    this.type = obj.type;
    this.group = obj.group;
    this.state = obj.state;
    this.concurrencyIndex = obj.concurrencyIndex;
    this.concurrency = obj.concurrency;
    this.totalExecutions = obj.totalExecutions;
    this.totalErrors = obj.totalErrors;
    this.lastStatus = obj.lastStatus;
    this.lastLatency = obj.lastLatency;
    this.lastErrorMessage = obj.lastErrorMessage;

    if (obj.executionResults) {
        for (let result of obj.executionResults) {
         this.executionResults.push(new ExecutionResult(result));
        }
    }
    this.failureStats = new Map<string, Stats>();
    if (obj.failureStats) {
      for(let stat in obj.failureStats) {
        this.failureStats[stat] = new Stats(obj.failureStats[stat])
      }
      this.steps = Object.keys(this.failureStats)
    }
  }
}

export class ExecutionResult {
  public status: string;
  public concurrencyIndex: number;
  public logMessages: LogMessage[];
  public latency: string;
  public finished: string;

  constructor(obj?: Object) {
    this.logMessages = [];
    if (obj) {
      this.unmarshal(obj as ExecutionResult);
    }
  }

  protected unmarshal(obj: ExecutionResult) {
    this.status = obj.status;
    this.concurrencyIndex = obj.concurrencyIndex;
    this.latency = obj.latency;
    this.finished =  getLocalTime(obj.finished);

    for (let message of obj.logMessages) {
      this.logMessages.push(new LogMessage(message));
    }
  }
}

export class LogMessage {
  public message: string;
  public step: string;
  public success: boolean;
  public time: string;

  constructor(obj?: Object) {
    if (obj) {
      this.unmarshal(obj as LogMessage);
    }
  }

  protected unmarshal(obj: LogMessage) {
    this.message = obj.message;
    this.step = obj.step;
    this.success = obj.success;
    this.time = getLocalTime(obj.time)
  }
}

function getLocalTime(utcTime):string {
  let time = new Date(utcTime);
  //Pad the minutes, seconds, and milliseconds with the zeroes we require.
  let minutes = (time.getMinutes() < 10 ? '0': '') + time.getMinutes();
  let seconds = (time.getSeconds() < 10 ? '0': '') + time.getSeconds();
  let milliseconds = (time.getMilliseconds() < 10 ? '0': '')?(time.getMilliseconds() < 100 ? '0': '') + time.getMilliseconds():time.getMilliseconds();
  return time.getHours() + ":" + minutes + ":" + seconds + "." + milliseconds;
}

export class Stats {
  public successes: number;
  public failures: number;

  constructor(obj?: Object) {
    if(obj) {
      this.unmarshal(obj as Stats)
    }
  }
  protected unmarshal(obj: Stats) {
    this.failures = obj.failures;
    this.successes = obj.successes;
  }
}

