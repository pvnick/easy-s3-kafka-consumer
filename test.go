package main
  
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
}