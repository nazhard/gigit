package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
)

const err_clone_msg = "Something error when clone "

var user_repo string

func main() {
  if len(os.Args) >= 2 {
    if strings.Contains(os.Args[1], "-1") {
      user_repo = os.Args[2]

      _, err := exec.Command(
        "git", "clone", "--depth=1", "https://github.com/"+user_repo+".git").Output()

      handleErr(err, err_clone_msg + user_repo)
      handleErr(err, user_repo + " Not Found, Onii")
    } else {
      user_repo = os.Args[1]

      _, err := exec.Command(
        "git", "clone", "https://github.com/"+user_repo+".git").Output()

      handleErr(err, err_clone_msg + user_repo)
      handleErr(err, user_repo + " Not Found, Onii")
    }
  } else {
    fmt.Println("Onii-chan! anata wa need repository!\n")
    fmt.Println("Example: gigit nazhard/gigit")
  }
}

func handleErr(err error, printText any) {
  if err != nil {
    fmt.Println(printText)
  }
}
