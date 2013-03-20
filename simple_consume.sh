#!/bin/bash

PATH_TO_KAFKA_BIN_DIR="/home/ubuntu/kafka/bin"
PATH_TO_S3CMD_BIN_DIR="/home/ubuntu/s3cmd-1.5.0-alpha1"
S3_TARGET_DIR="s3://datatests/twitter"
BLOCK_SIZE=10000000

$PATH_TO_KAFKA_BIN_DIR/kafka-console-consumer.sh --zookeeper localhost:2181 --topic twitter --from-beginning \
    | python handle_consumption.py --blocksize=$BLOCK_SIZE --s3cmd-dir="$PATH_TO_S3CMD_BIN_DIR" --s3target-dir="$S3_TARGET_DIR"
