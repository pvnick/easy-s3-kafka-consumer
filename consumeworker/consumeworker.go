package consumeworker

import (
    "os/exec"
    "os"
    "bufio"
)

type Config struct {
    S3CmdPath               string
    KafkaCLIConsumerPath    string
    S3TargetPrefix          string
    KafkaTopic              string
    Zookeeper               string
    BlockSize               uint32
}

type Worker interface {
    Start()
    bufferBlock()
    flushToS3()
}

type Instance struct {
    Worker
    CLIConsumerProc *os.Process
    config Config
    currentBufferSize uint32

}

func (this *Instance) Start() {
    cmd := exec.Command(this.config.KafkaCLIConsumerPath, 
        "--zookeeper", this.config.Zookeeper,
        "--topic", this.config.KafkaTopic)
    cmd.Start()
    this.CLIConsumerProc = cmd.Process
    this.currentBufferSize = 0xFFFFFFFF
    for this.currentBufferSize < this.config.BlockSize { 
        this.bufferBlock(cmd)
    }
    this.flushToS3()
}

func (this *Instance) bufferBlock(cliProcCmd *exec.Cmd) {
    outBuf := bufio.NewReader(cliProcCmd.StdoutPipe)
    errBuf := bufio.NewReader(cliProcCmd.StderrPipe)

    go func() {
        for {
            outLine, outErr := outBuf.ReadBytes('\n')
            if outErr != nil {
                break
            }

        }
    }()

    go func() {
        for {
            errLine, errErr := errBuf.ReadBytes('\n')
            if errErr != nil {
                break
            }
            // is this supposed to be os.Stderr ?
            os.Stdout.Write(errLine)
        }
    }()
}

func (this *Instance) flushToS3() {

}

func New(config Config) *Instance {
    return &Instance{
        config: config,
    }
}