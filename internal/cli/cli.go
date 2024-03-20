package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/nazhard/gigit"
)

func fetcher(user, repo, cache_path, commit, file_name, goberr, sub, url, c_url string, do bool) {
	user_repo := user + "/" + repo

	fmt.Println("Fetching " + gchalk.Underline(url))

	cache := filepath.Join(cache_path, user, repo)

	err := os.MkdirAll(cache, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if do == true {
		_, err = gigit.Get(user_repo, commit, goberr, c_url)
		if err != nil {
			log.Fatal(
				gchalk.Red("Internal error"))
		}
	}

	file := filepath.Join(cache, file_name)
	gigit.ExtractGz(file, repo, sub)

	fmt.Println(
		gchalk.Green("repo success full downloaded."))
}

func Exec(user, repo, cache_path, subdir string) error {
	var (
		url       string
		sub       string
		file_name string
		err       error
	)

	goberr := filepath.Join(cache_path, user, repo)

	// By default it uses HEAD as commit, which is why GetLatestCommit() is so important.
	// When the download is complete, there will be a tarball,
	// in the tarball there is a pattern of names user-repo-commit_hash
	commit_hash := gigit.GetLatestCommit(user + "/" + repo)

	if subdir != "" {
		file_name = repo + ".tar.gz"
		sub = user + "-" + repo + "-" + commit_hash + "/" + subdir
		url, err = gigit.Get(user+"/"+repo, "HEAD", goberr, "")
	} else {
		file_name = subdir + ".tar.gz"
		sub = user + "-" + repo + "-" + commit_hash
		url, err = gigit.Get(os.Args[1], "HEAD", goberr, "")
		fmt.Println("Fetching " + gchalk.Underline(url))
	}

	if err != nil {
		fmt.Println(err)
		fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))

		return fmt.Errorf("Error")
	}

	if err == nil {
		fetcher(user, repo, cache_path, commit_hash, file_name, goberr, sub, url, "", false)
	}

	return nil
}

// Why did I do this? Because I was stressed out with the errors that were appearing.
// I think making the code twice is easier to read and maintain.
func SharpExec(u_r, user, repo, cache_path string) error {
	var (
		hash string
		err  error
	)

	index_one := strings.Index(os.Args[1], "#")
	if index_one != -1 {
		hash = os.Args[1][index_one+1:]
	}

	goberr := filepath.Join(cache_path, user, repo)

	url, err := gigit.Get(u_r, hash, goberr, "")

	if strings.Contains(hash, "v") {
		v := hash[1:]
		c_url := "https://github.com/" + u_r + "/archive/refs/tags/" + hash + ".tar.gz"
		url, err = gigit.Get(u_r, v, goberr, c_url)
	}

	sub := user + "-" + repo + "-" + hash
	file_name := repo + ".tar.gz"

	if err != nil {
		// When commit hash errors.  gigit will check if it's a branch or not.
		url, commit, err := gigit.GetCommitBranch(u_r, hash)

		if err != nil {
			fmt.Println(err)
			fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))

			return fmt.Errorf("Upps")
		}

		if err == nil {
			fetcher(user, repo, cache_path, commit, file_name, goberr, sub, url, url, true)
		}
	}

	if err == nil {
		fetcher(user, repo, cache_path, hash, file_name, goberr, sub, url, "", false)
	}

	return nil
}
