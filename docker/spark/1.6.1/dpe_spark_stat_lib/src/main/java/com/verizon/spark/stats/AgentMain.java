/*
 * Copyright (C) 2017 Verizon. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.verizon.spark.stats;

import java.lang.instrument.Instrumentation;
import java.util.Map;

/*
  Agent entry class.
 */
public class AgentMain {
    private static final int DEFAULT_POLL_INTERVAL_SEC = 10;

    public static void premain(String args, Instrumentation ins) {
      for (Map.Entry<String, String> e : System.getenv().entrySet()) {
        System.out.println("Name : " + e.getKey() + ", value: "  + e.getValue());
      }

      try {
        int pollIntervalSecs = DEFAULT_POLL_INTERVAL_SEC;
        String statsInterval = System.getenv("STATS_INTERVAL");
        if (statsInterval != null && !statsInterval.isEmpty()) {
          pollIntervalSecs = Integer.parseInt(statsInterval);
        }

        System.out.println("Starting with pollInterval=" + pollIntervalSecs);
        new StatCollector(pollIntervalSecs).start();
      } catch (Exception e) {
          e.printStackTrace();
      }
    }
}
