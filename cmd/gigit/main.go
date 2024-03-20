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
	"log"
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

func exec(user, repo, cache_path string) error {
	goberr := filepath.Join(cache_path, user, repo)

	// By default it uses HEAD as commit, which is why GetLatestCommit() is so important.
	// When the download is complete, there will be a tarball,
	// in the tarball there is a pattern of names user-repo-commit_hash
	url, err := gigit.Get(os.Args[1], "HEAD", goberr, "")

	commit_hash := gigit.GetLatestCommit(os.Args[1])
	sub := user + "-" + repo + "-" + commit_hash
	file_name := repo + ".tar.gz"

	fmt.Println("Fetching " + gchalk.Underline(url))

	if err != nil {
		fmt.Println(err)
		fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))

		return fmt.Errorf("Error")
	}

	if err == nil {
		cache := filepath.Join(cache_path, user, repo)
		err = os.MkdirAll(cache, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		file := filepath.Join(cache, file_name)
		gigit.ExtractGz(file, repo, sub)

		fmt.Println(
			gchalk.Green("repo success full downloaded."))
	}

	return nil
}

// Why did I do this? Because I was stressed out with the errors that were appearing.
// I think making the code twice is easier to read and maintain.
func sharpExec(u_r, user, repo, cache_path string) error {
	var (
		hash string
		url  string
		err  error
	)

	index_one := strings.Index(os.Args[1], "#")
	if index_one != -1 {
		hash = os.Args[1][index_one+1:]
	}

	goberr := filepath.Join(cache_path, user, repo)

	url, err = gigit.Get(u_r, hash, goberr, "")

	if strings.Contains(hash, "v") {
		v := hash[1:]
		c_url := "https://github.com/" + u_r + "/archive/refs/tags/" + hash + ".tar.gz"
		url, err = gigit.Get(u_r, v, goberr, c_url)
	}

	sub := user + "-" + repo + "-" + hash
	file_name := repo + ".tar.gz"

	if err != nil {
		url, commit, err := gigit.GetCommitBranch(u_r, hash)
		if err != nil {
			fmt.Println(err)
			fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))

			return fmt.Errorf("Upps")
		}

		if err == nil {
			fmt.Println("Fetching " + gchalk.Underline(url))

			cache := filepath.Join(cache_path, user, repo)

			err = os.MkdirAll(cache, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			url, err = gigit.Get(u_r, commit, goberr, url)
			if err != nil {
				fmt.Println(
					gchalk.Red("Internal error"))
			}

			file := filepath.Join(cache, file_name)
			gigit.ExtractGz(file, repo, sub)

			fmt.Println(
				gchalk.Green("repo success full downloaded."))
		}
	}

	if err == nil {
		fmt.Println("Fetching " + gchalk.Underline(url))

		cache := filepath.Join(cache_path, user, repo)
		err = os.MkdirAll(cache, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		file := filepath.Join(cache, file_name)
		gigit.ExtractGz(file, repo, sub)

		fmt.Println(
			gchalk.Green("repo success full downloaded."))
	}

	return nil
}

func main() {
	// By default tells the user what version of gigit they are using
	fmt.Println(
		gchalk.Bold("using gigit " + version + "\n"))

	usr_cache_path, _ := os.UserCacheDir()

	cache_path := filepath.Join(usr_cache_path, "gigit")
	_, err := os.Stat(cache_path)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(cache_path, os.ModePerm)
	}

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
			if !strings.Contains(os.Args[1], "#") {
				array := strings.Split(os.Args[1], "/")
				user, repo := array[0], array[1]

				err := exec(user, repo, cache_path)
				if err != nil {
					gigit.Clone("https://github.com", os.Args[1])
				}
			}

			if strings.Contains(os.Args[1], "#") {
				eps := strings.Split(os.Args[1], "#")
				array := strings.Split(eps[0], "/")
				user, repo := array[0], array[1]

				user_repo := user + "/" + repo

				err := sharpExec(eps[0], user, repo, cache_path)
				if err != nil {
					gigit.Clone("https://github.com", user_repo)
				}
			}
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

		if os.Args[1] != "help" && !strings.Contains(os.Args[1], "/") {
			invalid()
		}
	}

	// Currently we do not have a feature to handle arguments with more than 2
	if len(os.Args) >= 3 {
		invalid()
	}
}
