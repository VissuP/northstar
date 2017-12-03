#!/bin/sh
#set -x
servicePid=0
service=$1
debug=
loggerPid=
loggerBinary="/usr/local/bin/logger"
sleepDuration=2
isTermReceived=

precheck() {
    if [ -f "$loggerBinary" ] &&  [  -f "$service" ]
    then
        debugLog "File[$loggerBinary $service] found "
    else
        errorLog "Either of file[$loggerBinary $service] not found..."
        sleep 10
        exit 1
    fi
}

debugLog() {
    if [ "$debug" != "" ]
    then
        echo "DEBUG::$*"
    fi
}

infoLog() {
    echo "INFO::$*"
}

errorLog() {
    echo "ERROR::$*"
}


getPid() {
    process=$1
    while true 
    do
        lPid=`ps -eaf | grep $process | grep -v startup | grep -v grep | awk {'printf ("%s ", $1)'}`
        if [ -z "$lPid" ] 
        then
            sleep $sleepDuration
            continue
        else
            echo "$lPid"
            break
        fi
    done
}

waitFor() {
    totalProcess=$#
    killedCount=0
    infoLog "Going to track $totalProcess process"
    while [ $killedCount != $totalProcess ]
    do
        killedCount=0
        sleep $sleepDuration
        for i in $*
        do 
            if [ ! -d "/proc/$i" ] 
            then 
                killedCount=$(($killedCount + 1))
                # Generate SIGTERM if any of the process found to be killed internally
                if  [ -z $isTermReceived ]
                then
                    infoLog "Process [$i] found to be killed, raising SIGTERM signal for all"
                    term_handler
                fi
            fi
        done
        debugLog "Total killed process $killedCount"
    done
    infoLog "Tracked $totalProcess process killed "
}


validation() {
    if [ -z "$loggerPid" ] || [ -z "$servicePid" ]
    then 
        errorLog "Error::Process Id could not be found for mlogger or Service[$service]"
        return
    fi
}

term_handler() {
    validation
    # Sending SIGTERM to logger first so that they start dumping message on
    # stdout which may get collected by platform kafka
    for i in $loggerPid
    do
        infoLog "Sending SIGTERM to process id $i"
        kill -SIGTERM $i
    done

    if [ $servicePid -ne 0 ]; then
        infoLog "Sending SIGTERM to process id $servicePid"
        kill -SIGTERM $servicePid
    fi
    isTermReceived=1
}

precheck
trap 'term_handler' SIGTERM

debugLog "Going to start $service with mlogger"
{ exec $service  2>&1 1>&3 3>&- | $loggerBinary -st=tcp -ost=false; }  3>&1 1>&2 |  $loggerBinary -st=tcp &

servicePid=`getPid $service`
loggerPid=`getPid $loggerBinary`

infoLog "Process id of [$service] is [$servicePid] and for [$loggerBinary] are [$loggerPid]"
debugLog "Wait for service & logger binary to finish"
waitFor $servicePid $loggerPid
infoLog "Exiting"
