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
  
  _, err := exec.Command(
    "git", "clone", "https://github.com/"+user_repo+".git").Output()

  if err != nil && len(os.Args) >= 2 {
    fmt.Println(user_repo, "not found, onii")
  }
}
