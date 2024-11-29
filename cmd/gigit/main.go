/*
Gigit the repository downloader

Gigit is inspired by degit a repository downloader written by Rich Harris in JavaScript.
Gigit has almost the same features as degit. Downloading repositories, caching features, and some still in development.

# Simple usage example:

	gigit <user>/<repo>

In the example above, gigit will download the GitHub respository with the username <user> and the repository that has the name <repo>.

# Spesific branch, commit hash, tag

You can use specific branches, commits, or tags with a `#`

	gigit user/repo#dev
	gigit user/repo#691c0bf

	// on spesific tag, "v" is required
	gigit user/repo#v1.0.0

# Subdirectory

Download sub directory only.

	gigit user/repo/dir

	gigit nazhard/gigit/cmd/gigit

# Commands

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

	"github.com/faulbert/gigit"
	"github.com/faulbert/gigit/internal/cli"
	"github.com/jwalton/gchalk"
)

const version = "v0.1.0"

// invalid prints an error message when the user provides an incorrect or inappropriate format.
func invalid() {
	fmt.Println(gchalk.Red("Error: invalid format"))
	fmt.Println(gchalk.Bold("Type 'gigit help' for more information."))
	fmt.Println(gchalk.Green("Valid format: 'gigit user/repo'"))
}

// is responsible for handling user input for fetching a repository.
// It checks if the input format is valid and then calls the appropriate functions to download the repository.
func handleFetchInput(args int, input string) {
	// Define the cache path for storing downloaded repositories
	cachePath := cli.CachePath

	// Ensure the cache directory exists
	_, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(cachePath, os.ModePerm) // Create cache directory if it does not exist
	}

	// Handle the case where there are exactly two arguments
	if args == 2 {
		if strings.Contains(input, "/") { // Check if input contains '/' (user/repo format)
			if strings.Contains(input, "#") { // Check for commit hash or version
				// Split the input into user/repo and commit/branch part
				parts := strings.Split(input, "#")
				repoParts := strings.Split(parts[0], "/")
				user, repo := repoParts[0], repoParts[1]

				userRepo := user + "/" + repo

				// Try fetching the repository with the specific commit/branch
				err := cli.SharpExec(user, repo)
				if err != nil {
					// If fetching fails, clone the repository
					gigit.Clone("https://github.com", userRepo, false)
				}
			} else {
				// Handle the case where only user/repo is provided
				repoParts := strings.Split(input, "/")
				user, repo := repoParts[0], repoParts[1]

				// Try fetching the repository without any subdir or commit/branch
				err := cli.Exec(user, repo, "")
				if err != nil {
					// If fetching fails, clone the repository
					gigit.Clone("https://github.com", user+"/"+repo, false)
				}
			}
		}
	}

	// Handle the case where there are 3 or more arguments (invalid input)
	if args >= 3 {
		invalid() // Call the invalid function if there are too many arguments
	}
}

func main() {
	// Print the version of the gigit tool that the user is using
	fmt.Println(gchalk.Bold("Using gigit version " + version + "\n"))

	// Get the number of arguments passed to the program
	args := len(os.Args)
	input := ""

	// If there are arguments, capture the first one
	if args > 1 {
		input = os.Args[1]
	}

	// Handle case where no command is given (only the program name)
	if args == 1 {
		fmt.Println(gchalk.Red("Onii-chan! anata wa need repository!"))
		fmt.Println(gchalk.Blue("Example: gigit nazhard/gigit"))
	} else if !strings.Contains(input, "/") { // Check if input doesn't have '/' (not in user/repo format)
		// Handle the 'help' command
		if input == "help" {
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

Examples: 
    gigit nazhard/gigit
    gigit nazhard/gigit/cmd`)
		} else if args >= 2 { // If there are more than 1 argument
			switch input {
			case "clone":
				gigit.Clone("https://github.com", os.Args[2], false)
			case "c1", "1":
				gigit.Clone("https://github.com", os.Args[2], true)
			default:
				invalid() // If the command is not recognized, call invalid() to show an error
			}
		}
	} else { // Handle valid user/repo input
		handleFetchInput(args, input)
	}
}
