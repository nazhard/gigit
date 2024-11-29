package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/faulbert/gigit"
	"github.com/jwalton/gchalk"
)

// Define user cache path and custom cache directory
var userCachePath, _ = os.UserCacheDir()
var CachePath = filepath.Join(userCachePath, "gigit")

// Global variables for URL, commit hash, and error handling
var (
	url    string
	inside string
	hash   string
	err    error
)

// fetcher is responsible for downloading and extracting a Git repository.
func fetcher(user, repo, commit, goberr, sub, url, c_url string) {
	// Inform the user about the repository being fetched
	fmt.Println("Fetching " + gchalk.Underline(url))

	// If a custom URL is provided, fetch the data from it
	if c_url != "" {
		_, err := gigit.Get(user+"/"+repo, commit, goberr, c_url)
		if err != nil {
			log.Fatal(gchalk.Red("Internal error"))
		}
	}

	// Define the file path where the repository archive will be stored
	file := filepath.Join(goberr, repo+".tar.gz")

	// Extract the contents of the archive
	gigit.Extract(file, sub, ".", 2)

	// Rename the extracted directory based on commit hash or repository name
	if len(commit) == 7 {
		_ = os.Rename(user+"-"+repo+"-"+commit, repo)
	} else if len(commit) > 7 {
		u := commit[:7] // Shorten commit hash to 7 characters
		_ = os.Rename(user+"-"+repo+"-"+u, repo)
	} else {
		_ = os.Rename(sub, repo)
	}

	// Print success message
	fmt.Println(gchalk.Green("Repo successfully downloaded."))
}

// Exec is responsible for downloading and extracting the latest commit of a repository.
func Exec(user, repo, subdir string) error {
	// Define the path for storing the repository download
	goberr := filepath.Join(CachePath, user, repo)
	err = os.MkdirAll(goberr, os.ModePerm)
	if err != nil {
		log.Fatal(err) // Exit if unable to create the directory
	}

	// Get the latest commit hash for the repository (defaults to HEAD)
	commit_hash := gigit.LatestCommit(user + "/" + repo)
	if len(commit_hash) > 7 {
		hash = commit_hash[:7] // Use the first 7 characters of the commit hash
	}

	// Define the path inside the downloaded archive based on the provided subdir
	if subdir != "" {
		inside = user + "-" + repo + "-" + hash + "/" + subdir
		url, err = gigit.Get(user+"/"+repo, "HEAD", goberr, "")
	} else {
		inside = user + "-" + repo + "-" + hash
		url, err = gigit.Get(os.Args[1], "HEAD", goberr, "")
	}

	// Handle errors when downloading the repository
	if err != nil {
		fmt.Println(err)
		fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))
		return fmt.Errorf("Error")
	}

	// Once the repository is downloaded, proceed to fetch and extract it
	if err == nil {
		fetcher(user, repo, commit_hash, goberr, inside, url, "")
	}

	return nil
}

// SharpExec handles fetching specific commits or tags for a repository.
func SharpExec(user, repo string) error {
	// Combine the user and repository into a single string
	user_repo := user + repo

	// Extract the commit hash from the arguments if provided
	index_one := strings.Index(os.Args[1], "#")
	if index_one != -1 {
		hash = os.Args[1][index_one+1:]
	}

	// Define the path for storing the repository download
	goberr := filepath.Join(CachePath, user, repo)

	// Fetch the repository using the specified commit or tag
	url, err := gigit.Get(user_repo, hash, goberr, "")

	// If a version tag is provided, handle it separately
	if strings.Contains(hash, "v") {
		v := hash[1:] // Strip the 'v' prefix from the version tag
		c_url := "https://github.com/" + user_repo + "/archive/refs/tags/" + hash + ".tar.gz"
		url, err = gigit.Get(user_repo, v, goberr, c_url)
		version := strings.TrimPrefix(hash, "v")
		inside = repo + "-" + version
	} else {
		inside = user + "-" + repo + "-" + hash
	}

	// Handle errors when downloading the repository
	if err != nil {
		// If there is an error, attempt to fetch the commit or branch
		url, commit, err := gigit.CommitBranch(user_repo, hash)

		if err != nil {
			fmt.Println(err)
			fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))
		}

		// If fetching the commit or branch is successful, proceed to extract it
		if err == nil {
			fetcher(user, repo, commit, goberr, inside, url, url)
		}
	}

	// If no errors, proceed to extract the repository
	if err == nil {
		fetcher(user, repo, hash, goberr, inside, url, "")
	}

	return nil
}
