package consumeworker

import (
    "os/exec"
    "os"
    "bufio"
    "io/ioutil"
)

type Config struct {
    TempDir                 string
    S3CmdPath               string
    KafkaCLIConsumerPath    string
    S3TargetPrefix          string
    KafkaTopic              string
    Zookeeper               string
    BlockSize               uint32
}

type Instance struct {
    CleanedUp chan int
    config Config
    cliConsumerProc *os.Process
    alive bool
}

func (this *Instance) Start() {
    tmpOutputFileHandle, err := ioutil.TempFile(this.TempDir, "consumer-cli-output")
    defer func() {
        tmpOutputFileHandle.Close()
        os.Remove(tmpOutputFileHandle.Name())
    }
    if err == nil { 
        this.tmpOutputFileHandle = tmpOutputFileHandle
        cmd := exec.Command(this.config.KafkaCLIConsumerPath, 
            "--zookeeper", this.config.Zookeeper,
            "--topic", this.config.KafkaTopic)
        outPipe, err := cmd.StdoutPipe()
        if err != nil {
            log.Fatal(err)
        }
        cmd.Start()
        this.cliConsumerProc = cmd.Process
        this.alive = true
        for this.alive {
            this.bufferBlock(outPipe)
            this.flushToS3()
            tmpOutputFileHandle.Truncate(0)
        }
    }
    this.CleanedUp <- 1
}

func (this *Instance) Apoptosis(tombstoneMessage string) {
    if this.alive == true {
        this.alive = false
        this.cliConsumerProc.Kill()
        log.Println("Worker killed with message: ", tombstoneMessage)
    }
}

func (this *Instance) bufferBlock(cliProcStdOut *io.Reader, tempFile *os.File) {
    outBuf := bufio.NewReader(cliProcStdOut)
    approximateMaxBlockSize := this.config.BlockSize
    for currBlockSize := 0; currBlockSize < approximateMaxBlockSize; this.alive {
        //FIXME: allow custom line separators
        outLine, err := outBuf.ReadBytes('\n')
        if err != nil {
            this.Apoptosis("Error reading line from consumer cli process")
        } else {
            _, err := tempFile.Write(outLine)
            if err != nil {
                this.Apoptosis("Error writing line to temporary file")
            } else {
                currBlockSize += len(outLine)
            }
        }
    }
}

func (this *Instance) flushToS3() {

}

func New(config Config) *Instance {
    return &Instance{
        CleanedUp: make(chan int, 1)
        config: config,
    }
}