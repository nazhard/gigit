package gigit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/mholt/archiver/v3"
)

type commitT struct {
	Sha string `json:"sha"`
}

type branchT struct {
	Commit struct {
		Sha string `json:"sha"`
	} `json:"commit"`
}

var (
	user_repo string
	file_name string
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

// Get is used to download the repository.
// Get requires the repository name and its commit and output path for downloaded repository
// and optionally, a custom url.
// Get returns url used to download repo as string and returns an error when an error occurs.
func Get(name, commit, destination, c_url string) (string, error) {
	url := c_url
	if c_url == "" {
		url = "https://api.github.com/repos/" + name + "/tarball" + "/" + commit
	}

	res, err := client.Get(url)
	if err != nil {
		return url, errors.New(gchalk.Red("Error when fetching " + url))
	}
	defer res.Body.Close()

	fullPath := name
	lastSlashIndex := strings.LastIndex(fullPath, "/")

	if lastSlashIndex != -1 {
		file_name = fullPath[lastSlashIndex+1:]
	}

	if res.StatusCode != http.StatusOK {
		return url, errors.New(gchalk.Red("Error " + name + " not found"))
	}

	if res.StatusCode == http.StatusOK {
		out_path := filepath.Join(destination, file_name+".tar.gz")
		_, err := os.Stat(destination)
		if os.IsNotExist(err) {
			_ = os.MkdirAll(destination, os.ModePerm)
		}

		f, _ := os.Create(out_path)
		defer f.Close()

		_, err = f.ReadFrom(res.Body)
		if err != nil {
			fmt.Println(err)
		}
	}

	return url, nil
}

// This function will fetch the latest commit hash from the repository for the extractor's use.
// GetLatestCommit is only used when the user does not provide a specific commit hash.
//
// Example: "nazhard/gigit"
func LatestCommit(user_repo string) string {
	url := "https://api.github.com/repos/" + user_repo + "/commits"
	res, err := client.Get(url)
	if err != nil {
		log.Fatal("Error, maybe your internet connection is bad")
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	var commits []commitT

	err = json.Unmarshal(data, &commits)
	if err != nil {
		return ""
	}

	var c string
	if len(commits) > 1 {
		commit := commits[0]
		c = commit.Sha
	}

	return c
}

func CommitBranch(user_repo, branch string) (string, string, error) {
	url_to_fetch := "https://api.github.com/repos/" + user_repo + "/branches/" + branch
	url := "https://github.com/" + user_repo + "/archive/refs/tags/" + branch + ".tar.gz"

	res, err := client.Get(url_to_fetch)
	if err != nil {
		return "", "", fmt.Errorf("Upps branch not found!")
	}
	res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	var bran branchT

	err = json.Unmarshal(data, &bran)
	if err != nil {
		return "", "", fmt.Errorf("upps, error")
	}

	commit := bran.Commit.Sha

	return url, commit, nil
}

// Checks the cache stored in the default cache directory for gigit.
//
// CheckCache returns a boolean. If the cache exists, it returns true. If it doesn't exist, it returns false.
//
// CheckCache requires path for the default cache path, name for the repository name, and commit hash.
func CheckCache(path, name, commit_hash string) bool {
	cache_path := filepath.Join(path, name, commit_hash)

	_, err := os.Stat(cache_path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// Extract downloaded repository.
//
// Extract file .tar.gz as "source",
// then an output directory for the extracted contents of the .tar.gz file as "destination",
// and a specific directory/path inside the .tar.gz file to extract its contents as "target".
func Extract(source, target, destination string, strip int) {
	archive := archiver.NewTarGz()

	archive.StripComponents = strip
	archive.Extract(source, target, destination)
}

// Clone repositories using git instead.
//
// This will automatically be used when gigit does not find the intended repository.
//
// This is especially useful when you want to type "gigit user/repo" instead of "git clone https...".
// In simple terms, it is meant to clone a private repository.
// Set depth to true if you want using "--depth 1"
//
// user_repo here refers to a string containing "user/repo", not "user" or "repo" only!
func Clone(host, user_repo string, depth bool) {
	var cmd *exec.Cmd

	if depth == true {
		cmd = exec.Command("git", "clone", "--depth=1", host+"/"+user_repo+".git")
	} else {
		cmd = exec.Command("git", "clone", host+"/"+user_repo+".git")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Start()

	defer cmd.Wait()
}
