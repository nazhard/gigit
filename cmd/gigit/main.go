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
	"path/filepath"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/nazhard/gigit"
)

const version = "v0.1.0"

// This function is only to print an error if the user-supplied format is incorrect or inappropriate.
func invalid() {
	fmt.Println(
		gchalk.Red("Error: invalid format"))
	fmt.Println(
		gchalk.Bold("We are unable to accept more than 2 arguments.\nType 'gigit help' for more information."))
	fmt.Println(
		gchalk.Green("Valid format: 'gigit user/repo'"))
}

func exec() {
	usr_cache_path, _ := os.UserCacheDir()
	cache_path := filepath.Join(usr_cache_path, "gigit")
	_, err := os.Stat(cache_path)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(cache_path, os.ModePerm)
	}

	// By default it uses HEAD as commit, which is why GetLatestCommit() is so important.
	// When the download is complete, there will be a tarball,
	// in the tarball there is a pattern of names user-repo-commit_hash
	url, err := gigit.Get(os.Args[1], "HEAD", cache_path)
	commit_hash := gigit.GetLatestCommit(os.Args[1])
	fmt.Println("Fetching " + gchalk.Underline(url))

	if err != nil {
		fmt.Println(err)
		fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))

		gigit.Clone("https://github.com", os.Args[1])
	}

	if err == nil {
		n := strings.Split(os.Args[1], "/")
		sub := n[0] + "-" + n[1] + "-" + commit_hash
		file_name := n[1] + ".tar.gz"
		cache := filepath.Join(cache_path, n[0], n[1])
		_ = os.MkdirAll(cache, os.ModePerm)

		file := filepath.Join(cache, file_name)
		gigit.ExtractGz(file, n[1], sub)

		fmt.Println("repo success full downloaded.")
	}
}

func main() {
	// By default tells the user what version of gigit they are using
	fmt.Println(
		gchalk.Bold("using gigit " + version + "\n"))

	// Handles when the user does not give any commands
	if len(os.Args) == 1 {
		fmt.Println(
			gchalk.Red("Onii-chan! anata wa need repository!"))
		fmt.Println(
			gchalk.Blue("Example: gigit nazhard/gigit"))
	}

	if len(os.Args) == 2 {
		// Checks to see if os.Args[1] (argument) has "/" or not.
		if strings.Contains(os.Args[1], "/") {
			exec()
		}

		// Handles "help" commands
		if os.Args[1] == "help" {
			fmt.Println(`Usage: gigit user/repo
       gigit user/repo/subdir
       gigit host:user/repo
       gigit host:user/repo/subdir

Host: github or gitlab

Examples: gigit nazhard/gigit
          gigit github:nazhard/gigit`)
		}
	}

	// Currently we do not have a feature to handle arguments with more than 2
	if len(os.Args) >= 3 {
		invalid()
	}
}
