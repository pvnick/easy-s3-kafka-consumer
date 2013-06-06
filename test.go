package main
  
func main(){
    workers := []int{1,2,3,4}
    for worker := range workers {
        print(worker)
    }
}