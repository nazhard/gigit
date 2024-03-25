/*
Gigit the repository downloader

Gigit is inspired by degit a repository downloader written by Rich Harris in JavaScript.
Gigit has almost the same features as degit. Downloading repositories, caching features, and some still in development.

## Simple usage example:

	gigit <user>/<repo>

In the example above, gigit will download the GitHub respository with the username <user> and the repository that has the name <repo>.

## Spesific branch, commit hash, tag

You can use specific branches, commits, or tags with a `#`

	gigit user/repo#dev
	gigit user/repo#691c0bf

	// on spesific tag, "v" is required
	gigit user/repo#v1.0.0

## Subdirectory

Download sub directory only.

	gigit user/repo/dir

	gigit nazhard/gigit/cmd/gigit

## Commands

Clone instead of download. With cloning, you will get a .git folder

	gigit clone user/repo

	// Clone with `--depth=1` if you just want to fix typo
	gigit c1 user/repo
	// or
	gigit 1 user/repo
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

func cute(args int, one string) {
	cache_path := cli.CachePath
	_, err := os.Stat(cache_path)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(cache_path, os.ModePerm)
	}

	if args == 2 {
		// Checks to see if os.Args[1] (argument) has "/" or not.
		if strings.Contains(one, "/") {
			if strings.Contains(one, "#") {
				eps := strings.Split(one, "#")
				array := strings.Split(eps[0], "/")
				user, repo := array[0], array[1]

				user_repo := user + "/" + repo

				err := cli.SharpExec(user, repo)
				if err != nil {
					gigit.Clone("https://github.com", user_repo, false)
				}
			} else {
				slash := strings.Count(one, "/")
				if slash == 1 {
					array := strings.Split(one, "/")
					user, repo := array[0], array[1]

					err := cli.Exec(user, repo, "")
					if err != nil {
						gigit.Clone("https://github.com", user+"/"+repo, false)
					}
				}

				if slash >= 2 {
					array := strings.SplitN(one, "/", 3)
					user, repo, dir := array[0], array[1], array[2]

					err := cli.Exec(user, repo, dir)
					if err != nil {
						gigit.Clone("https://github.com", user+"/"+repo, false)
					}
				}
			}
		}
	}

	if args >= 3 {
		invalid()
	}
}

func main() {
	// By default tells the user what version of gigit they are using
	fmt.Println(
		gchalk.Bold("using gigit " + version + "\n"))

	args := len(os.Args)
	one := ""
	if args > 1 {
		one = os.Args[1]
	}

	// Handles when the user does not give any commands
	if args == 1 {
		fmt.Println(
			gchalk.Red("Onii-chan! anata wa need repository!"))
		fmt.Println(
			gchalk.Blue("Example: gigit nazhard/gigit"))
	} else if !strings.Contains(one, "/") {
		if one == "help" {
			fmt.Println(`Usage: gigit [command]
       gigit user/repo
       gigit user/repo/subdir
       gigit user/repo#dev
       gigit user/repo#v1.0.0
       gigit user/repo#7k3b2kw

Commands:
    help
        print this help message.
    clone
        clone repository instead.
    c1, 1
        same as clone but use "--depth=1"
  

Examples: gigit nazhard/gigit
          gigit nazhard/gigit/cmd`)
		} else if args >= 2 {
			switch one {
			case "clone":
				gigit.Clone("https://github.com", os.Args[2], false)
			case "c1", "1":
				gigit.Clone("https://github.com", os.Args[2], true)
			default:
				invalid()
			}
		}
	} else {
		cute(args, one)
	}
}
