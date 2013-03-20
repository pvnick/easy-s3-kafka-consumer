#!/bin/bash

PATH_TO_KAFKA_BIN_DIR="/home/ubuntu/kafka/bin"
TMP_FILE=""
MAX_SIZE=1000000

function GET_TEMP_FILE {
    STARTED=0
    TMP_FILE=""
    while [ $STARTED -eq 0 -o -e "$TMP_FILE" ]; do
        STARTED=1
        CURRENT_TS=$(date +%s)
        TMP_FILE="/tmp/simple_s3_consumer_"$CURRENT_TS"_"$RANDOM$RANDOM
    done
    touch $TMP_FILE
}

function KILL_PID_SUBPROCESSES {
    MY_PID_GROUP=$(ps x -o "%p %r" | grep -r "^$MY_PID" | cut -d" " -f2)
    kill -TERM -$MY_PID_GROUP
}

function KILL_MY_SUBPROCESSES {
    MY_PID=$$
    #kill all subprocesses to clean up background processes
    KILL_PID_SUBPROCESSES $MY_PID
}

trap KILL_SUBPROCESSES INT 

GET_TEMP_FILE

$PATH_TO_KAFKA_BIN_DIR/kafka-console-consumer.sh --zookeeper localhost:2181 --topic twitter --from-beginning > $TMP_FILE &
CONSUMER_PID=$!

while  [ 1 ]; do
    if [ $(stat -c %s $TMP_FILE) -gt $MAX_SIZE ]; then
        echo "Chunk finished"

        #reset the temp file locaiton so we can work with the current one before deleting it
        PROCESSING_FILE_PATH=$TMP_FILE
        GET_TEMP_FILE

        #kill the consumer 
        KILL_PID_SUBPROCESSES
        
        echo "Uploading to S3"

        echo "Done"
        rm -f $PROCESSING_FILE_PATH
    fi
done
