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

var userCachePath, _ = os.UserCacheDir()

var CachePath = filepath.Join(userCachePath, "gigit")

var (
	url    string
	inside string
	hash   string
	err    error
)

func fetcher(user, repo, commit, goberr, sub, url, c_url string) {
	fmt.Println("Fetching " + gchalk.Underline(url))

	if c_url != "" {
		_, err := gigit.Get(user+"/"+repo, commit, goberr, c_url)
		if err != nil {
			log.Fatal(
				gchalk.Red("Internal error"))
		}
	}

	file := filepath.Join(goberr, repo+".tar.gz")

	gigit.Extract(file, sub, ".", 2)

	if len(commit) == 7 {
		_ = os.Rename(user+"-"+repo+"-"+commit, repo)
	} else if len(commit) > 7 {
		u := commit[:7]
		_ = os.Rename(user+"-"+repo+"-"+u, repo)
	} else {
		_ = os.Rename(sub, repo)
	}

	fmt.Println(
		gchalk.Green("repo success full downloaded."))
}

func Exec(user, repo, subdir string) error {
	goberr := filepath.Join(CachePath, user, repo)
	err = os.MkdirAll(goberr, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// By default it uses HEAD as commit, which is why GetLatestCommit() is so important.
	// When the download is complete, there will be a tarball,
	// in the tarball there is a pattern of names user-repo-commit_hash
	commit_hash := gigit.LatestCommit(user + "/" + repo)
	if len(commit_hash) > 7 {
		hash = commit_hash[:7]
	}

	if subdir != "" {
		inside = user + "-" + repo + "-" + hash + "/" + subdir
		url, err = gigit.Get(user+"/"+repo, "HEAD", goberr, "")
	} else {
		inside = user + "-" + repo + "-" + hash
		url, err = gigit.Get(os.Args[1], "HEAD", goberr, "")
	}

	if err != nil {
		fmt.Println(err)
		fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))

		return fmt.Errorf("Error")
	}

	if err == nil {
		fetcher(user, repo, commit_hash, goberr, inside, url, "")
	}

	return nil
}

// Why did I do this? Because I was stressed out with the errors that were appearing.
// I think making the code twice is easier to read and maintain.
func SharpExec(user, repo string) error {
	user_repo := user + repo

	index_one := strings.Index(os.Args[1], "#")
	if index_one != -1 {
		hash = os.Args[1][index_one+1:]
	}

	goberr := filepath.Join(CachePath, user, repo)

	url, err := gigit.Get(user_repo, hash, goberr, "")

	if strings.Contains(hash, "v") {
		v := hash[1:]
		c_url := "https://github.com/" + user_repo + "/archive/refs/tags/" + hash + ".tar.gz"
		url, err = gigit.Get(user_repo, v, goberr, c_url)
		version := strings.TrimPrefix(hash, "v")
		inside = repo + "-" + version
	} else {
		inside = user + "-" + repo + "-" + hash
	}

	if err != nil {
		// When commit hash errors.  gigit will check if it's a branch or not.
		url, commit, err := gigit.CommitBranch(user_repo, hash)

		if err != nil {
			fmt.Println(err)
			fmt.Print(
				gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))
		}

		if err == nil {
			fetcher(user, repo, commit, goberr, inside, url, url)
		}
	}

	if err == nil {
		fetcher(user, repo, hash, goberr, inside, url, "")
	}

	return nil
}
