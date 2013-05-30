package consumeworker

import (
    "os/exec"
    "os"
)

type Config struct {
    S3CmdPath               string
    KafkaCLIConsumerPath    string
    S3TargetPrefix          string
    KafkaTopic              string
    Zookeeper               string
    BlockSize               int
}

type Worker interface {
    Start()
}

type Instance struct {
    Worker
    config Config
}

func (this *Instance) Start(newProcessChan chan *os.Process) {
    cmd := exec.Command(this.config.KafkaCLIConsumerPath, 
        "--zookeeper", this.config.Zookeeper,
        "--topic", this.config.KafkaTopic)
    cmd.Start()
    newProcessChan <- cmd.Process
}

func New(config Config) *Instance {
    return &Instance{
        config: config,
    }
}