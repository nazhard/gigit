package main

import (
  "fmt"
  "os"
  "os/exec"
)

func main() {

  var user_repo string

  if len(os.Args) >= 2 {
    user_repo = os.Args[1]
  } else {
    fmt.Println("Onii-chan! anata wa need repository!\n")
    fmt.Println("Example: gigit nazhard/gigit")
  }
  
  out, err := exec.Command(
    "git", "clone", "https://github.com/"+user_repo+".git").Output()

  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(string(out))

}
