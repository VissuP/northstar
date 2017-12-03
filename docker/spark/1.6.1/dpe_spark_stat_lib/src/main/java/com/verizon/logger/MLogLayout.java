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

package com.verizon.logger;

import org.apache.log4j.Layout;
import org.apache.log4j.Level;
import org.apache.log4j.MDC;
import org.apache.log4j.helpers.PatternConverter;
import org.apache.log4j.helpers.PatternParser;
import org.apache.log4j.spi.LoggingEvent;

import java.lang.management.ManagementFactory;
import java.lang.management.RuntimeMXBean;
import java.util.Map;

/**
 * Created by ramakve on 4/6/16.
 * This class is used to format the log output in Mlog format.
 */
public class MLogLayout extends Layout {
    private static final String ALARM_LEVEL = "ALARM";
    private static final String ERROR_LEVEL = "ERROR";
    private static final String INFO_LEVEL = "INFO";
    private static final String DEBUG_LEVEL = "DEBUG";

    private String conversionPattern;
    private PatternConverter originalConverter;
    private PatternConverter alarmConverter;
    private PatternConverter errorConverter;
    private PatternConverter infoConverter;
    private PatternConverter debugConverter;

    private String format(LoggingEvent event, PatternConverter converter) {
        StringBuffer buffer = new StringBuffer();
        while (converter != null) {
            converter.format(buffer, event);
            converter = converter.next;
        }
        return buffer.toString();
    }

    @Override
    public String format(LoggingEvent event) {
        if(MDC.get("PID") == null) {
            initialize();
        }

        if (event.getLevel().equals(Level.TRACE)) {
            return format(event, debugConverter);
        } else if (event.getLevel().equals(Level.DEBUG)) {
            return format(event, debugConverter);
        } else if (event.getLevel().equals(Level.INFO)) {
            return format(event, infoConverter);
        } else if (event.getLevel().equals(Level.WARN)) {
            return format(event, infoConverter);
        } else if (event.getLevel().equals(Level.ERROR)) {
            return format(event, errorConverter);
        } else if (event.getLevel().equals(Level.FATAL)) {
            return format(event, alarmConverter);
        }

        return format(event, originalConverter);
    }

    @Override
    public boolean ignoresThrowable() {
        return false;
    }


    private void initialize() {
        Map<String, String> envs = System.getenv();
        for(Map.Entry<String, String> env : envs.entrySet() ) {
            MDC.put("ENV:" + env.getKey(), env.getValue());
        }
        RuntimeMXBean rt = ManagementFactory.getRuntimeMXBean();
        String pid = rt.getName();
        if (pid.contains("@")) {
            pid = pid.substring(0, pid.indexOf("@"));
        }
        MDC.put("PID", pid);

        String conversionPattern = getConversionPattern();
        if (conversionPattern != null && conversionPattern.contains("%p")) {
            String alarmLayout = conversionPattern.replace("%p", ALARM_LEVEL);
            String errorLayout = conversionPattern.replace("%p", ERROR_LEVEL);
            String infoLayout = conversionPattern.replace("%p", INFO_LEVEL);
            String debugLayout = conversionPattern.replace("%p", DEBUG_LEVEL);
            PatternParser parser = new PatternParser(alarmLayout);
            alarmConverter = parser.parse();
            parser = new PatternParser(errorLayout);
            errorConverter = parser.parse();
            parser = new PatternParser(infoLayout);
            infoConverter = parser.parse();
            parser = new PatternParser(debugLayout);
            debugConverter = parser.parse();
            parser = new PatternParser(conversionPattern);
            originalConverter = parser.parse();
        }
    }

    public String getConversionPattern() {
        return this.conversionPattern;
    }
    public void setConversionPattern(String conversionPattern) {
        this.conversionPattern = conversionPattern;
        initialize();
    }

    @Override
    public void activateOptions() {

    }
}

