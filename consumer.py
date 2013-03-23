#!/usr/bin/python

import tempfile as tf
import os
import atexit
import subprocess
import datetime
import time

class Consumer:
    tmpFile = None
    tmpFilePath = ""
    s3CmdDir = ""
    kafkaBinDir = ""
    s3TargetPrefix = ""
    kafkaTopic = ""
    blockSize = 0
    zookeeper = ""
    consoleConsumerProc = None
    dataPreprocessorCallback = None

    def __init__(self, s3CmdDir, kafkaBinDir, s3TargetPrefix, kafkaTopic, blockSize = 1024 * 1024 * 10, zookeeper = "localhost:2181", dataPreprocessorCallback = None):
        self.s3CmdDir = s3CmdDir.rstrip("/")
        self.kafkaBinDir = kafkaBinDir.rstrip("/")
        self.s3TargetPrefix = s3TargetPrefix.rstrip("/")
        self.kafkaTopic = kafkaTopic
        self.blockSize = blockSize
        self.zookeeper = zookeeper
        self.dataPreprocessorCallback = dataPreprocessorCallback

        #take some liberties with exit behavior to ensure we always flush temporary files from disk to s3
        self.registerExitHandler()

    def run(self):
        print("Starting console-based consumer") 
        self.consoleConsumerProc = subprocess.Popen([
                self.kafkaBinDir + "/kafka-console-consumer.sh",
                "--zookeeper", self.zookeeper,
                "--topic", self.kafkaTopic
            ], stdout = subprocess.PIPE)

        print("Starting consumption loop")
        while True:
            #temporarily store data locally on disk
            self.bufferBlock()
            
            #flush to s3 once files get large enough (or if the application exits - see registerExitHandler)
            self.storeBlock()

    def bufferBlock(self):
        print("Consuming data and buffering to local disk")

        fileSize = 0
        self.tmpFile = tf.NamedTemporaryFile(delete=False)
        self.tmpFilePath = self.tmpFile.name

        #consumed data comes through from the console consumer stdout
        consumedData = self.consoleConsumerProc.stdout

        for line in iter(consumedData.readline, ""):
            #if a data preprocessor has been specified then pass the whole line to them and write what it passes back
            if (self.dataPreprocessorCallback != None):
                line = self.dataPreprocessorCallback(line)

            self.tmpFile.write(line)
            fileSize += len(line)

            if (fileSize >= self.blockSize):
                break

    def storeBlock(self):
        print("Flushing buffered data to s3 and deleting temporary local file") 

        if (self.tmpFile):
            self.tmpFile.close()

        now = datetime.datetime.now()

        filePath = self.tmpFilePath
        tmpFileBaseName = os.path.basename(filePath)

        #target file name formatted to sort alphabetacally by newest file and be safe from duplicates
        targetFileName = "%04d-%02d-%02d-%d-%07d-%s" % (now.year, now.month, now.day, time.mktime(now.timetuple()), now.microsecond, tmpFileBaseName)
        targetFullPath = self.s3TargetPrefix + targetFileName

        s3CmdOutput = subprocess.check_output([self.s3CmdDir + "/s3cmd", "put", filePath, targetFullPath])
        #todo: error handling here

        os.unlink(filePath)

    def registerExitHandler(self):
        print("Registering exit handler")
        atexit.register(self.beforeExit)

    def beforeExit(self):
        print("Exit signal detected; cleaning up")
        self.consoleConsumerProc.terminate()
        self.storeBlock()
