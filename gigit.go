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

type (
	Commit struct {
		Sha string `json:"sha"`
	}

	Branch struct {
		Commit struct {
			Sha string `json:"sha"`
		} `json:"commit"`
	}
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

// Get downloads a specific commit or branch tarball from a GitHub repository and saves it locally.
//
// Parameters:
//   - name: The GitHub repository in the format "owner/repo" (e.g., "nazhard/gigit").
//   - commit: The specific commit hash or branch name to fetch (e.g., "main" or a hash like "abc123").
//   - destination: The directory where the tarball will be saved.
//   - url: (Optional) A custom URL to fetch the tarball from. If empty, the default GitHub API tarball URL is used.
//
// Returns:
//   - A string containing the URL used for the fetch.
//   - An error if the fetch fails, the directory cannot be created, or the tarball cannot be saved.
//
// Behavior:
//   - Constructs the default GitHub tarball URL if no custom URL is provided.
//   - Sends a GET request to fetch the tarball from the URL.
//   - Extracts the repository name from the `name` parameter and generates a tarball filename.
//   - Creates the destination directory if it does not exist.
//   - Writes the tarball content to a file in the destination directory.
//   - Returns the URL used and any encountered errors.
func Get(name, commit, destination, url string) (string, error) {
	// Default to constructing Github Tarball URL if not provided
	u := url
	if u == "" {
		u = "https://api.github.com/repos/" + name + "/tarball" + "/" + commit
	}

	// Fetch content from URL
	res, err := client.Get(u)
	if err != nil {
		return u, errors.New(gchalk.Red("Error when fetching " + u))
	}
	defer res.Body.Close()

	// Extract file name from repository name
	fileName := name
	lastSlashIndex := strings.LastIndex(fileName, "/")

	if lastSlashIndex != -1 {
		fileName = fileName[lastSlashIndex+1:]
	}

	// Check if status code is not 200 (OK)
	if res.StatusCode != http.StatusOK {
		return u, errors.New(gchalk.Red("Error " + name + " not found"))
	}

	// Construct the outout file path
	outPath := filepath.Join(destination, fileName+".tar.gz")

	// Create destination directory if it doesn't exists
	if _, err = os.Stat(destination); os.IsNotExist(err) {
		if err := os.MkdirAll(destination, os.ModePerm); err != nil {
			return u, errors.New(gchalk.Red("Error creating destination directory: " + err.Error()))
		}
	}

	// Create the output file
	f, err := os.Create(outPath)
	if err != nil {
		return u, errors.New(gchalk.Red("Error creating file: " + err.Error()))
	}
	defer f.Close()

	// Write the response body to the output file
	if _, err = f.ReadFrom(res.Body); err != nil {
		return u, errors.New(gchalk.Red("Error writing to file: " + err.Error()))
	}

	return u, nil
}

// LatestCommit fetches the latest commit hash from a specified GitHub repository.
//
// Parameters:
//   - userRepo: The GitHub repository in the format "owner/repo" (e.g., "nazhard/gigit").
//
// Returns:
//   - The latest commit hash as a string.
//   - If the request fails or there are no commits, the function returns an empty string.
//
// Behavior:
//   - Sends a GET request to the GitHub API to retrieve the list of commits for the given repository.
//   - Parses the JSON response to extract the hash of the most recent commit.
//   - If there is an error (e.g., invalid repository, no internet connection), it logs the issue and exits the program.
func LatestCommit(userRepo string) string {
	// Construct the API URL to fetch the commits
	url := "https://api.github.com/repos/" + userRepo + "/commits"

	// Send the GET request
	res, err := client.Get(url)
	if err != nil {
		log.Fatal("Error fetching commits. Check your internet connection or the repository name.")
	}
	defer res.Body.Close()

	// Read the response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	// Parse the JSON response into a slice of Commit structs
	var commits []Commit
	if err := json.Unmarshal(data, &commits); err != nil {
		log.Println("Error parsing JSON: ", err)
		return ""
	}

	// Extract and return the SHA of the latest commit if available
	if len(commits) > 0 {
		return commits[0].Sha
	}

	// Return an empty string if no commits are found
	log.Println("No commits found for repository: ", userRepo)
	return ""
}

// CommitBranch fetches the commit hash of a specific branch in a GitHub repository.
// It also constructs the URL to download the branch as a tar.gz archive.
//
// Parameters:
//   - userRepo: The GitHub repository in the format "owner/repo" (e.g., "nazhard/gigit").
//   - branch: The branch name (e.g., "main").
//
// Returns:
//   - The URL to download the branch archive.
//   - The commit hash of the branch.
//   - An error if the branch does not exist or the request fails.
func CommitBranch(userRepo, branch string) (string, string, error) {
	// URLs for API and archive download
	apiURL := "https://api.github.com/repos/" + userRepo + "/branches/" + branch
	archiveURL := "https://github.com/" + userRepo + "/archive/refs/heads/" + branch + ".tar.gz"

	// Fetch branch data from GitHub API
	res, err := client.Get(apiURL)
	if err != nil {
		return "", "", fmt.Errorf("Error fetching branch: %w", err)
	}
	defer res.Body.Close()

	// Read the response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("Error reading response body: %w", err)
	}

	// Parse JSON response into a Branch struct
	var b Branch
	if err := json.Unmarshal(data, &b); err != nil {
		return "", "", fmt.Errorf("Error parsing branch data: %w", err)
	}

	// Extract the commit SHA
	commit := b.Commit.Sha

	return archiveURL, commit, nil
}

// CheckCache checks if a specific cache file or directory exists.
//
// Parameters:
//   - path: The base directory where caches are stored.
//   - name: The name of the repository or resource (e.g., "nazhard/gigit").
//   - commitHash: The commit hash or unique identifier for the cache.
//
// Returns:
//   - `true` if the cache exists.
//   - `false` if the cache does not exist.
//
// Behavior:
//   - Constructs the full cache path by combining the base path, name, and commit hash.
//   - Uses `os.Stat` to check if the specified path exists.
//   - Returns `false` if the path does not exist or there is an error.
func CheckCache(path, name, commitHash string) bool {
	// Construct the full path to the cache
	cachePath := filepath.Join(path, name, commitHash)

	// Check if the path exists
	_, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// Extract extracts a tar.gz archive to a specified destination.
//
// Parameters:
//   - source: Path to the tar.gz archive.
//   - target: Specific file/directory to extract (leave empty to extract all).
//   - destination: Directory where the extracted content will be placed.
//   - strip: Number of leading path components to remove from file paths.
func Extract(source, target, destination string, strip int) {
	archive := archiver.NewTarGz()

	archive.StripComponents = strip
	archive.Extract(source, target, destination)
}

// Clone clones a Git repository from a specified host and repository name.
//
// Parameters:
//   - host: The Git server's base URL (e.g., "https://github.com").
//   - userRepo: The repository in the format "owner/repo" (e.g., "nazhard/gigit").
//   - depth: A boolean value indicating whether to perform a shallow clone.
//   - `true`: Clones only the latest commit (`--depth=1`).
//   - `false`: Clones the entire repository history.
//
// Behavior:
//   - Constructs a `git clone` command based on the provided parameters.
//   - Executes the command and streams its output (stdout and stderr) to the console.
//   - Ensures that the command is properly started and waited on to finish.
func Clone(host, userRepo string, depth bool) {
	// Determine the git clone command based on the depth parameter
	var cmd *exec.Cmd
	repoURL := host + "/" + userRepo + ".git"

	if depth {
		cmd = exec.Command("git", "clone", "--depth=1", repoURL)
	} else {
		cmd = exec.Command("git", "clone", repoURL)
	}

	// Set command output to be streamed to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting git clone:", err)
		return
	}

	// Wait for the command to finish
	defer cmd.Wait()
}
