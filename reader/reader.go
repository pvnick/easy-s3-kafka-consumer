package reader

import (
    "fmt"
    "os/exec"
)

type ReadInit interface {
    Start()     string
}

type Instance struct {
    Foobar      string
    ReadInit
}

func (r *Instance) Start() {
    fmt.Println(r.Foobar)
    cmd := exec.Command("yes")
    cmd.Start()
}

func New() *Instance {
    blah := new(Instance)
    blah.Foobar="something"
    return blah
}