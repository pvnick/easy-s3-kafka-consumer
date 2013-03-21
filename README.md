simple-s3-kafka-consumer
========================

A simple project for dumping data from kafka to s3. Great for streaming large amounts of web data into s3 for use with mapreduce.

Usage: ./consumer --OPT1=val --OPT2=val ...

Options:
  --s3cmd-dir           (required) absolute path to the directory where s3cmd is stored
  --s3-target-dir       (required) s3://bucketname/foldername/ - the s3 directory where the data will get dumped
  --kafka-bin-dir       (required) absolute path to the directory where the kafka binary files are stored
  --topic               (required) kafka topic to consume
  --blocksize           (optional) amount of data to buffer to local disk before flushing to s3 (default 10 mb). flushing 
                                   to s3 creates a new file in the bucket
  --zookeeper           (optional) host:port for zookeeper. default is localhost:2181

Details
========================
