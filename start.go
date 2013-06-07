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

func cleanup(workerSemaphore chan *consumeworker.Instance) {
    log.Println("Attempting to exit gracefully")
    log.Println("Killing child processes")
    for worker := range workerSemaphore {
        log.Println("Killing worker")
        worker.Apoptosis("Host sent kill signal")
        <-worker.CleanedUp 
    }
    log.Println("Cleanup completed successfully")
}

func main() {
    log.Println("Initializing")
    coresToUtilize = runtime.NumCPU()
    enableMultiCore()
    
    config := consumeworker.Config{
        TempDir: "/tmp/easys3",
        S3CmdPath: "yes", //"/usr/local/bin/s3cmd",
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

    alive := true
    exitSignalChan := make(chan os.Signal, 1)
    workerSemaphore := make(chan *consumeworker.Instance, coresToUtilize)
    signal.Notify(exitSignalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
    
    go func(){
        <-exitSignalChan
        alive = false
        log.Println("Caught exit signal")
        cleanup(workerSemaphore)
    }()

    log.Println("Launching workers")
    for alive {
        worker := consumeworker.New(config)
        workerSemaphore <- worker
        go func() {
            workers.Start()
            //Start() returns when the worker dies, so keep the semaphore full with live ones
            <-workerSemaphore
        }
    }
}




