simple-s3-kafka-consumer
========================

A simple project for dumping data from kafka to S3. Great for streaming large amounts of web data into S3 for use with mapreduce.

```
Usage: ./run-consumer --OPT1=val --OPT2=val ...

Options:
  --s3cmd-dir           (required) absolute path to the directory where s3cmd is stored
  --s3-target-prefix    (required) s3://bucketname/foldername/datafileprefix_ - the s3 directory and file prefix where 
                                   the data will get dumped
  --kafka-bin-dir       (required) absolute path to the directory where the kafka binary files are stored
  --topic               (required) kafka topic to consume
  --blocksize           (optional) amount of data to buffer to local disk before flushing to s3 (default 10 mb). flushing 
                                   to s3 creates a new file in the bucket
  --zookeeper           (optional) host:port for zookeeper. default is localhost:2181
```

Details
-------------------------

This script, written in python, uses the console consumer packaged with kafka to stream data into S3. It does so by buffering blocks of data to the local file system and periodically flushing to S3 once the blocks reach a certain size.

The benefit of this is that the code is independent of the kafka api version (0.7.* vs 0.8.* - and beyond) and you get to use zookeeper for all it's coordination power rather than manually specifying brokers.

The consumer class allows you to supply a preprocessor callback so you can manipulate the data before writing to S3:
```
def callbackExample(line):
    print("line passed to data preprocessor: " + line)

    #note: line includes newline character, so make sure to return that
    return arg
    
consumerInst = consumer.Consumer(
  s3CmdDir = s3CmdDir,
  kafkaBinDir = kafkaBinDir,
  s3TargetPrefix = s3TargetPrefix,
  kafkaTopic = topic,
  dataPreprocessorCallback = callbackExample
)
```

Dependencies
-------------------------

* Kafka binaries
* [s3cmd](http://s3tools.org/s3cmd)
