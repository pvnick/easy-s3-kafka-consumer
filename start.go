package main

import (
    "fmt"
    "os/exec"
    "runtime"
    //"github.com/pvnick/easy-s3-kafka-consumer/consumer"
)

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
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
    fmt.Println("Initializing")
    enableMultiCore()
    testS3Cmd()
}