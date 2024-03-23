/*
Gigit the repository downloader

Gigit is inspired by degit a repository downloader written by Rich Harris in JavaScript.
Gigit has almost the same features as degit. Downloading repositories, caching features, and some still in development.

Simple usage example:

	gigit <user>/<repo>

In the example above, gigit will download the GitHub respository with the username <user> and the repository that has the name <repo>.
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/nazhard/gigit"
	"github.com/nazhard/gigit/internal/cli"
)

const version = "v0.1.0"

// This function is only to print an error if the user-supplied format is incorrect or inappropriate.
func invalid() {
	fmt.Println(
		gchalk.Red("Error: invalid format"))
	fmt.Println(
		gchalk.Bold("Type 'gigit help' for more information."))
	fmt.Println(
		gchalk.Green("Valid format: 'gigit user/repo'"))
}

func main() {
	args := len(os.Args)
	one := ""
	if args >= 2 {
		one = os.Args[1]
	}

	// By default tells the user what version of gigit they are using
	fmt.Println(
		gchalk.Bold("using gigit " + version + "\n"))

	cache_path := cli.CachePath
	_, err := os.Stat(cache_path)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(cache_path, os.ModePerm)
	}

	// Handles when the user does not give any commands
	if args == 1 {
		fmt.Println(
			gchalk.Red("Onii-chan! anata wa need repository!"))
		fmt.Println(
			gchalk.Blue("Example: gigit nazhard/gigit"))
	}

	if args == 2 {
		// Checks to see if os.Args[1] (argument) has "/" or not.
		if strings.Contains(one, "/") {
			if !strings.Contains(one, "#") {
				count := strings.Count(one, "/")
				if count == 1 {
					array := strings.Split(one, "/")
					user, repo := array[0], array[1]

					err := cli.Exec(user, repo, "")
					if err != nil {
						gigit.Clone("https://github.com", user+"/"+repo)
					}
				}

				if count >= 2 {
					array := strings.SplitN(one, "/", 3)
					user, repo, dir := array[0], array[1], array[2]

					err := cli.Exec(user, repo, dir)
					if err != nil {
						gigit.Clone("https://github.com", user+"/"+repo)
					}
				}
			}

			if strings.Contains(one, "#") {
				eps := strings.Split(one, "#")
				array := strings.Split(eps[0], "/")
				user, repo := array[0], array[1]

				user_repo := user + "/" + repo

				err := cli.SharpExec(eps[0], user, repo)
				if err != nil {
					gigit.Clone("https://github.com", user_repo)
				}
			}
		}

		// Handles "help" commands
		if !strings.Contains(one, "/") {
			switch one {
			case "help":
				fmt.Println(`Usage: gigit user/repo
       gigit user/repo/subdir
       gigit host:user/repo
       gigit host:user/repo/subdir

Host: github or gitlab

Examples: gigit nazhard/gigit
          gigit github:nazhard/gigit`)
			case "clone":
				// just_clone(false)
			case "c1", "1":
				// just_clone(true)
			default:
				invalid()
			}
		}
	}

	// Currently we do not have a feature to handle arguments with more than 2
	if args >= 3 {
		invalid()
	}
}
