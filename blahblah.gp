package main
 /* 
  import (
      "flag"
  )
  
  // Example 1: A single string flag called pecies" with default value "gopher".
  var species = flag.String("species", "gopher", "the species we are studying")
  
  // Example 2: Two flags sharing a riable, so we can have a shorthand.
  // The order of initialization is defined, so make sure both use the
  // same default value. They must be set up th an init function.
  var gopherType string
  
  func init() {
      flag.StringVar(&gopherType, "gopher_type", defaultGopher, usage)
      flag.StringVar(&gopherType, "g", defaultGopher, usage+" (shorthand)")
  }

func main(){
    flag.Parse()
    print("hi")
}*/

import (
    "bufio"
    "log"
    "os"
    "os/exec"
)

func main() {
    cmd := exec.Command("echo","-e", "asefasef\naefasef\nasefasefa")
    outPipe, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    errPipe, err := cmd.StderrPipe()
    if err != nil {
        log.Fatal(err)
    }

    if err := cmd.Start(); err != nil {
        log.Fatal(err)
    }

    outBuf := bufio.NewReader(outPipe)
    errBuf := bufio.NewReader(errPipe)

    go func() {
        for {
            outLine, outErr := outBuf.ReadBytes('\n')
            if outErr != nil {
                break
            }
            os.Stdout.Write(outLine)
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

    err = cmd.Wait()
    if err != nil {
        log.Fatal(err)
    }
}
