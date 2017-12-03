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

import org.apache.log4j.Level;
import org.apache.log4j.Priority;
import org.apache.log4j.spi.LocationInfo;
import org.apache.log4j.spi.LoggingEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.management.*;
import java.lang.management.*;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.*;
import java.util.concurrent.TimeUnit;

/**
 * This class is responsible for printing dpespark statistics.
 */
public class StatCollector extends Thread {
    private Logger logger = LoggerFactory.getLogger(getClass());

    // This is only used for output formatting.
    private org.apache.log4j.Logger log4JLogger = org.apache.log4j.Logger.getLogger(getClass());

    private int pollIntervalSecs;

    private static final String STDOUT_APPENDER = "stdout";

    private MemoryMXBean memory = ManagementFactory.getMemoryMXBean();
    private List<MemoryPoolMXBean> memoryPools = ManagementFactory.getMemoryPoolMXBeans();
    private ThreadMXBean threads = ManagementFactory.getThreadMXBean();
    private RuntimeMXBean runtime = ManagementFactory.getRuntimeMXBean();
    private List<GarbageCollectorMXBean> garbageCollectors = ManagementFactory.getGarbageCollectorMXBeans();
    private OperatingSystemMXBean os = ManagementFactory.getOperatingSystemMXBean();

    public StatCollector(int pollIntervalSecs) {
        this.pollIntervalSecs = pollIntervalSecs;
    }

    private double heapUsage() {
        final MemoryUsage usage = memory.getHeapMemoryUsage();
        return usage.getUsed() / (double) usage.getMax();
    }

    private double nonHeapUsage() {
        final MemoryUsage usage = memory.getNonHeapMemoryUsage();
        return usage.getUsed() / (double) usage.getMax();
    }

    private double heapUsed() {
        return memory.getHeapMemoryUsage().getUsed();
    }

    private Map<String, Double> memoryPoolUsage() {
        final Map<String, Double> pools = new TreeMap<String, Double>();
        for (MemoryPoolMXBean pool : memoryPools) {
            final double max = pool.getUsage().getMax() == -1 ?
                    pool.getUsage().getCommitted() :
                    pool.getUsage().getMax();
            pools.put(pool.getName(), pool.getUsage().getUsed() / max);
        }
        return Collections.unmodifiableMap(pools);
    }

    private double fileDescriptorUsage() {
        try {
            final Method getOpenFileDescriptorCount = os.getClass().getDeclaredMethod("getOpenFileDescriptorCount");
            getOpenFileDescriptorCount.setAccessible(true);
            final Long openFds = (Long) getOpenFileDescriptorCount.invoke(os);
            final Method getMaxFileDescriptorCount = os.getClass().getDeclaredMethod("getMaxFileDescriptorCount");
            getMaxFileDescriptorCount.setAccessible(true);
            final Long maxFds = (Long) getMaxFileDescriptorCount.invoke(os);
            return openFds.doubleValue() / maxFds.doubleValue();
        } catch (NoSuchMethodException e) {
            return Double.NaN;
        } catch (IllegalAccessException e) {
            return Double.NaN;
        } catch (InvocationTargetException e) {
            return Double.NaN;
        }
    }

    private Map<State, Double> threadStatePercentages() {
        final Map<State, Double> conditions = new HashMap<State, Double>();
        for (State state : State.values()) {
            conditions.put(state, 0.0);
        }

        final long[] allThreadIds = threads.getAllThreadIds();
        final ThreadInfo[] allThreads = threads.getThreadInfo(allThreadIds);
        int liveCount = 0;
        for (ThreadInfo info : allThreads) {
            if (info != null) {
                final State state = info.getThreadState();
                conditions.put(state, conditions.get(state) + 1);
                liveCount++;
            }
        }
        for (State state : new ArrayList<State>(conditions.keySet())) {
            conditions.put(state, conditions.get(state) / liveCount);
        }

        return Collections.unmodifiableMap(conditions);
    }


    private void sleep() {
        try {
            Thread.sleep(1000 * pollIntervalSecs);
        } catch (InterruptedException ie) {
            logger.error("Polling sleep interrupted", ie);
        }
    }

    public void run() {
      while (true) {
          VZMetricContext context = new VZMetricContext(pollIntervalSecs);
          printJvmMetrics(context);
          endContext(context);
          sleep();
      }
    }

    public static int getLineNumber() {
        return Thread.currentThread().getStackTrace()[2].getLineNumber();
    }
    public static String getMethod() {
        return Thread.currentThread().getStackTrace()[2].getMethodName();
    }

    private void endContext(VZMetricContext context) {
        context.close();

        LocationInfo info = new LocationInfo("", getClass().getCanonicalName(), getMethod(), ""+ getLineNumber());
        LoggingEvent event = new LoggingEvent(getClass().getCanonicalName(), log4JLogger, System.currentTimeMillis(), Level.INFO,  context.getMetricsString(), Thread.currentThread().getName(),null, null, info, null);
        String logStr = org.apache.log4j.Logger.getRootLogger().getAppender(STDOUT_APPENDER).getLayout().format(event);
        System.err.print(logStr.replaceFirst("INFO", "STAT"));
    }

    private void printJvmMetrics(VZMetricContext context) {
        context.addMetric("jvm.memory.heap.usage", "set", heapUsage());
        context.addMetric("jvm.memory.heap.used", "set", heapUsed());
        context.addMetric("jvm.memory.non_heap.usage", "set", nonHeapUsage());

        for (Map.Entry<String, Double> pool : memoryPoolUsage().entrySet()) {
            String gaugeName = String.format("jvm.memory.pool.%s.usage", pool.getKey());
            context.addMetric(gaugeName, "set", pool.getValue());
        }

        context.addMetric("jvm.daemon_threads.count", "set", threads.getDaemonThreadCount());
        context.addMetric("jvm.threads.count", "set", threads.getThreadCount());
        context.addMetric("jvm.uptime", "set", TimeUnit.MILLISECONDS.toSeconds(runtime.getUptime()));
        context.addMetric("jvm.fd_usage", "set", fileDescriptorUsage());

        for (Map.Entry<Thread.State, Double> entry : threadStatePercentages()
                .entrySet()) {
            String name = String.format("jvm.threads.state.%s",
                    entry.getKey());
            context.addMetric(name, "set", entry.getValue());
        }

        for (GarbageCollectorMXBean gc: garbageCollectors) {

            String p = String.format("jvm.gc.%s", gc.getName());
            context.addMetric(p + ".time", "set", gc.getCollectionTime());
            context.addMetric(p + ".runs", "set", gc.getCollectionCount());
        }
    }
    static class VZMetricContext {
        private StringBuilder stringBuilder;
        private long formattedDate;
        private int interval;
        private boolean first;
        private static int builderLength;

        public VZMetricContext(int interval) {
            this.interval = interval;
            stringBuilder = new StringBuilder(builderLength);
            stringBuilder.append("{");
            first = true;
            formattedDate = System.currentTimeMillis() / 1000;
        }

        public void close() {
            stringBuilder.append(",");
            add(stringBuilder, "module", System.getenv("MON_APP"));
            stringBuilder.append(",");
            add(stringBuilder, "etime", formattedDate);
            stringBuilder.append(",");
            add(stringBuilder, "interval", interval);
            stringBuilder.append("}");
            if(builderLength < stringBuilder.length()) {
                builderLength = stringBuilder.length();
            }
        }

        public String getMetricsString() {
            return stringBuilder.toString();
        }

        private void add(StringBuilder builder, String name, long value) {
            builder.append("\"" + name + "\"");
            builder.append(": ");
            builder.append(value);
        }

        private void add(StringBuilder builder, String name, int value) {
            builder.append("\"" + name + "\"");
            builder.append(": ");
            builder.append(value);
        }

        private void add(StringBuilder builder, String name, float value) {
            if(Float.isNaN(value) || Float.isInfinite(value)) {
                add(builder, name, ((Float)value).toString());
                return;
            }

            builder.append("\"" + name + "\"");
            builder.append(": ");
            builder.append(value);
        }

        private void add(StringBuilder builder, String name, double value) {
            if(Double.isNaN(value) || Double.isInfinite(value)) {
                add(builder, name, ((Double)value).toString());
                return;
            }

            builder.append("\"" + name + "\"");
            builder.append(": ");
            builder.append(value);
        }

        private void add(StringBuilder builder, String name, String value) {
            builder.append("\"" + name + "\"");
            builder.append(": ");
            builder.append("\"" + value + "\"");
        }

        public void addMetric(String name, String type, Object value) {
            if(value == null) {
                return;
            }

            if(!first) {
                stringBuilder.append(",");
            }

            if(value instanceof String) {
                add(stringBuilder, name, value.toString());
            } else if(value instanceof Long) {
                add(stringBuilder, name, (Long) value);
            } else if(value instanceof Integer) {
                add(stringBuilder, name, (Integer) value);
            } else if(value instanceof Float) {
                add(stringBuilder, name, (Float) value);
            } else if(value instanceof Double) {
                add(stringBuilder, name, (Double) value);
            } else {
                add(stringBuilder, name, value.toString());
            }
            first = false;
        }
    }
}
