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

type Worker interface {
    //public
    Start()
    Apoptosis()

    //private
    bufferBlock()
    flushToS3()
}

type Instance struct {
    //public
    Worker

    //private
    config Config
    CLIConsumerProc *os.Process
    isAlive bool
}

func (this *Instance) Start() {
    tmpOutputFileHandle, err := ioutil.TempFile(this.TempDir, "consumer-cli-output")
    defer func() {
        tmpOutputFileHandle.Close()

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
        this.CLIConsumerProc = cmd.Process
        this.isAlive = true
        for this.isAlive == true {
            this.bufferBlock(outPipe)
            this.flushToS3()
            tmpOutputFileHandle.Truncate(0)
        }
    }
}

    // make a buffer to keep chunks that are read
    buf := make([]byte, 1024)
    for {
        // read a chunk
        n, err := fi.Read(buf)
        if err != nil && err != io.EOF { panic(err) }
        if n == 0 { break }

        // write a chunk
        if _, err := fo.Write(buf[:n]); err != nil {
            panic(err)
        }
    }
}

func (this *Instance) Apoptosis(tombstoneMessage string) {
    log.Println("S3Cmd is working properly")
    this.isAlive = false
}


func (this *Instance) bufferBlock(cliProcStdOut *io.Reader) {
    outBuf := bufio.NewReader(cliProcStdOut)
    approximateMaxBlockSize := this.config.BlockSize

    for currBlockSize := 0; currBlockSize < approximateMaxBlockSize {
        //FIXME: allow custom line separators
        outLine, err := outBuf.ReadBytes('\n')
        if err != nil {
            
            currBlockSize += len(outLine)
        }
    }
}

func (this *Instance) flushToS3() {

}

func New(config Config) *Instance {
    return &Instance{
        config: config,
    }
}