package main

import (
    "fmt"
    "os/exec"
    "github.com/pvnick/easy-s3-kafka-consumer/consumer"
)

func main() {
    s3CmdPath, error := exec.LookPath("s3cmd")
    if error != nil {
        panic(error)
    } 
    cmd := exec.Command(s3CmdPath, "ls")
    output, err := cmd.CombinedOutput()
    if err != nil {
        panic("Error testing s3cmd:\n" + string(output))
    }

    fmt.Println(consumer.Hello_world())
}