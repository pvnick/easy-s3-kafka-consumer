package main

import (
    "fmt"
    "os/exec"
    "os/signal"
    "os"
    "runtime"
    //"time"
    "syscall"
    "github.com/pvnick/easy-s3-kafka-consumer/consumeworker"
)

var coresToUtilize int
var childProcs []*os.Process

func testS3Cmd(s3CmdPath string) {
    fmt.Println("Testing that S3Cmd works")
    testCmd := exec.Command(s3CmdPath, "ls")
    output, error := testCmd.CombinedOutput()
    if error != nil {
        panic("Error testing s3cmd:\n" + string(output))
    } 
    fmt.Println("S3Cmd is working properly")
}

func enableMultiCore() {
    runtime.GOMAXPROCS(coresToUtilize)
}

func cleanup() {
    fmt.Println("Attempting to exit gracefully")
    fmt.Println("Killing child processes")
    for _, proc := range childProcs {
        proc.Kill()
    }
    fmt.Println("Cleanup completed successfully")
}

func main() {
    fmt.Println("Initializing")
    defer cleanup()
    coresToUtilize = runtime.NumCPU()
    enableMultiCore()
    
    fmt.Println("skipping s3cmd check")
    if false {
        s3CmdPath, error := exec.LookPath("s3cmd")
        if error != nil {
            panic(error)
        } 
        testS3Cmd(s3CmdPath)
    }
    
    workerConfig := consumeworker.Config{
        S3CmdPath: "",
        KafkaCLIConsumerPath: "yes",
        S3TargetPrefix: "",
        KafkaTopic: "",
        Zookeeper: "localhost:2181",
        BlockSize: 1024 * 1024 * 10,
    }

    fmt.Println("Launching workers")
    workers := make([]*consumeworker.Instance, coresToUtilize)
    newChildProc := make(chan *os.Process)
    childProcs = make([]*os.Process, coresToUtilize)
    for i := 0; i < coresToUtilize; i++ {
        workers[i] = consumeworker.New(workerConfig)
        //each worker launches an instance of the kafka cli consumer
        //we receive the new child process  through the newChildProc channel
        go workers[i].Start(newChildProc)
        childProcs[i] = <-newChildProc
    }

    fmt.Println("Consumer started successfully")
    exitSignalChan := make(chan os.Signal, 1)
    signal.Notify(exitSignalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
    <-exitSignalChan
    fmt.Println("Caught exit signal")
}