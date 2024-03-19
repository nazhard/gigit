package gigit

import (
	"context"
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

	"github.com/codeclysm/extract/v3"
	"github.com/jwalton/gchalk"
)

type Commit struct {
	Sha string `json:"sha"`
}

var (
	user_repo string
	file_name string
)

// Get is used to download the repository.
// Get requires the repository name and its commit and output path for downloaded repository
// Get returns url used to download repo as string and returns an error when an error occurs.
func Get(name, commit, out_path string) (string, error) {
	url := "https://api.github.com/repos/" + name + "/tarball" + "/" + commit

	res, err := http.Get(url)
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
		output_goberr := filepath.Join(out_path, file_name+".tar.gz")
		_, err := os.Stat(out_path)
		if os.IsNotExist(err) {
			_ = os.MkdirAll(out_path, os.ModePerm)
		}
		f, _ := os.Create(output_goberr)

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
// Example: "gigit nazhard/gigit"
func GetLatestCommit(user_repo string) string {
	url := "https://api.github.com/repos/" + user_repo + "/commits"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	var commits []Commit

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
// Extract requires a file or path to a file in .tar.gz format as "in",
// then an output directory for the extracted contents of the .tar.gz file as "out"
// and a specific directory/path inside the .tar.gz file to extract its contents as "path".
func ExtractGz(in, out, path string) {
	file, _ := os.Open(in)

	var shift = func(path string) string {
		parts := strings.Split(path, string(filepath.Separator))
		parts = parts[1:]
		return strings.Join(parts, string(filepath.Separator))
	}

	extract.Gz(context.TODO(), file, out, shift)
}

// Clone repositories using git instead.
//
// This will automatically be used when gigit does not find the intended repository.
//
// This is especially useful when you want to type "gigit user/repo" instead of "git clone https...".
// In simple terms, it is meant to clone a private repository.
//
// user_repo here refers to a string containing "user/repo", not "user" or "repo" only!
func Clone(host, user_repo string) {
	cmd := exec.Command("git", "clone", host+"/"+user_repo+".git")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Start()

	defer cmd.Wait()
}
