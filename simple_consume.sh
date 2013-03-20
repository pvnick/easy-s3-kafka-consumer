#!/bin/bash

cd $(dirname "$0")

PATH_TO_KAFKA_BIN_DIR="/home/ubuntu/kafka/bin"
PATH_TO_S3CMD_BIN_DIR="/home/ubuntu/s3cmd-1.5.0-alpha1"
S3_TARGET_DIR="s3://datatests/twitter"
ZOOKEEPER_HOST="localhost"
ZOOKEEPER_PORT=2181

#topic to consume
KAFKA_TOPIC="twitter"

#BLCOK_SIZE refers to the amount of incomming data to buffer to the local file system before flushing to s3
BLOCK_SIZE=10485760 #10 megabytes

$PATH_TO_KAFKA_BIN_DIR/kafka-console-consumer.sh --zookeeper "$ZOOKEEPER_HOST:$ZOOKEEPER_PORT" --topic "$KAFKA_TOPIC" \
    | python handle_consumption.py --blocksize=$BLOCK_SIZE --s3cmd-dir="$PATH_TO_S3CMD_BIN_DIR" --s3target-dir="$S3_TARGET_DIR"
