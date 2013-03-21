simple-s3-kafka-consumer
========================

A simple project for dumping data from kafka to s3. Great for streaming large amounts of web data into s3 for use with mapreduce.

```
Usage: ./consumer --OPT1=val --OPT2=val ...

Options:
  --s3cmd-dir           (required) absolute path to the directory where s3cmd is stored
  --s3-target-dir       (required) s3://bucketname/foldername/ - the s3 directory where the data will get dumped
  --kafka-bin-dir       (required) absolute path to the directory where the kafka binary files are stored
  --topic               (required) kafka topic to consume
  --blocksize           (optional) amount of data to buffer to local disk before flushing to s3 (default 10 mb).
                                   flushing to s3 creates a new file in the bucket
  --zookeeper           (optional) host:port for zookeeper. default is localhost:2181
```

Details
-------------------------

This script, written in python, uses the console consumer packaged with kafka to stream data into s3. It does so by buffering blocks of data to the local file system and periodically flushing to s3 once the blocks reach a certain size.

The benefit of this is that the code is independent of the kafka api version (0.7.* vs 0.8.* - and beyond) and you get to use zookeeper for all it's coordination power rather than manually specifying brokers.

Dependencies
-------------------------

* Kafka binaries
* [s3cmd](http://s3tools.org/s3cmd)
