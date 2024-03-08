package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
)

var (
  user_repo string
  others []string
)

func main() {
  if len(os.Args) >= 2 {
    if strings.Contains(os.Args[1], "-1") {
      user_repo = os.Args[2]
      
      for _, i := range os.Args[3:] {
        others = append(others, i)
      }

      cmd := exec.Command("git", "clone", "--depth=1", "https://github.com/"+user_repo+".git", strings.Join(others, ", "))
      cmd.Stdout = os.Stdout
      cmd.Stderr = os.Stderr

      _ = cmd.Start()

      defer cmd.Wait()
    } else {
      user_repo = os.Args[1]

      for _, i := range os.Args[2:] {
        others = append(others, i)
      }

      cmd := exec.Command("git", "clone", "https://github.com/"+user_repo+".git", strings.Join(others, ", "))
      cmd.Stdout = os.Stdout
      cmd.Stderr = os.Stderr

      _ = cmd.Start()

      defer cmd.Wait()
    }
  } else {
    fmt.Println("Onii-chan! anata wa need repository!\n")
    fmt.Println("Example: gigit nazhard/gigit")
  }
}
