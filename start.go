package main

import (
    "fmt"
    "os/exec"
    "os/signal"
    "os"
    "runtime"
    //"time"
    "syscall"
    "github.com/pvnick/easy-s3-kafka-consumer/reader"
)

var coresToUtilize int

func testS3Cmd() {
    fmt.Println("Testing that S3Cmd works")
    s3CmdPath, error := exec.LookPath("s3cmd")
    if error != nil {
        panic(error)
    } 
    cmd := exec.Command(s3CmdPath, "ls")
    output, err := cmd.CombinedOutput()
    if err != nil {
        panic("Error testing s3cmd:\n" + string(output))
    } 
    fmt.Println("S3Cmd is working properly")
}

func enableMultiCore() {
    //this is meant to utilize as much of our machine as possible
    runtime.GOMAXPROCS(coresToUtilize)
}

func cleanup() {
    fmt.Println("bye now")
}

func main() {
    fmt.Println("Initializing")
    coresToUtilize = runtime.NumCPU()
    enableMultiCore()
    testS3Cmd()
    defer cleanup()

    readers := make([]*reader.Instance, coresToUtilize)
    for i := 0; i < coresToUtilize; i++ {
        readers[i] = reader.New()
        go readers[i].Start()
    }

    //capture any unexpected exit signals so we can kill all the cli consumers in the cleanup
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
    <-c
}