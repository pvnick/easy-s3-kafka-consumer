package main

func main(){
    sem := make(chan int, 4)
    sem <- 1
    sem <- 2
    sem <- 3
    sem <- 4
    blah := <-sem
    print(blah)
    blah = <-sem
    print(blah)
    blah = <-sem
    print(blah)
    blah = <-sem
    print(blah)
}
