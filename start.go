package main

import (
    "os/exec"
    "os/signal"
    "os"
    "log"
    "runtime"
    "syscall"
    "github.com/pvnick/easy-s3-kafka-consumer/consumeworker"
)

var coresToUtilize int
var workers []*consumeworker.Instance

func testS3Cmd(s3CmdPath string) {
    log.Println("Testing that S3Cmd works")
    testCmd := exec.Command(s3CmdPath, "ls")
    output, err := testCmd.CombinedOutput()
    if err != nil {
        log.Fatal("Error testing s3cmd:\n", string(output))
    } 
    log.Println("S3Cmd is working properly")
}

func enableMultiCore() {
    runtime.GOMAXPROCS(coresToUtilize)
}

func cleanup() {
    log.Println("Attempting to exit gracefully")
    log.Println("Killing child processes")
    for _, worker := range workers {
        go proc.Apoptosis()
    }
    //todo: wait for all to complete
    log.Println("Cleanup completed successfully")
}

func main() {
    log.Println("Initializing")
    defer cleanup()
    coresToUtilize = runtime.NumCPU()
    enableMultiCore()
    
    config := consumeworker.Config{
        TempDir: "/tmp/easys3",
        S3CmdPath: "/usr/local/bin/s3cmd",
        KafkaCLIConsumerPath: "/usr/local/kafka_install/bin/kafka-console-consumer.sh",
        S3TargetPrefix: "s3://pvnick_kafka_output/",
        KafkaTopic: "test-topic",
        Zookeeper: "localhost:2181",
        BlockSize: 1024 * 1024 * 10,
    }

    log.Println("skipping s3cmd check")
    if false {
        testS3Cmd(config.S3CmdPath)
        //TODO: test kafka cli consumer
    }
    
    log.Println("Launching workers")
    //TODO: use buffered channel sempahore to relaunch dead workers
    workers = make([]*consumeworker.Instance, coresToUtilize)
    for i := 0; i < coresToUtilize; i++ {
        workers[i] = consumeworker.New(config)
        go workers[i].Start()
    }

    log.Println("Consumer started successfully")
    exitSignalChan := make(chan os.Signal)
    signal.Notify(exitSignalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
    //<-exitSignalChan
    log.Println("Caught exit signal")
}




