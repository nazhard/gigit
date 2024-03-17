package main

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

const version = "v0.1.0"

func main() {
	fmt.Println(gchalk.Bold("using gigit " + version))

	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".gigit")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}

	if len(os.Args) > 1 {
		url, err := fetch(os.Args[1], "HEAD")
		c := getLatestCommit()
		fmt.Println("Fetching " + gchalk.Underline(url))

		if err != nil {
			fmt.Println(err)
			fmt.Print(gchalk.BrightBlack("\nRetry with cloning repository...\n\n"))

			clone()
		}

		if err == nil {
			fmt.Println("repo success full downloaded.")
			n := strings.Split(os.Args[1], "/")
			sub := n[0] + "-" + n[1] + "-" + c
			file_name := n[1] + ".tar.gz"
			extractGz(file_name, n[1], sub)

			cache := filepath.Join(path, n[0], n[1])
			_ = os.MkdirAll(cache, os.ModePerm)
			_ = os.Rename(file_name, cache+"/"+file_name)
		}
	} else {
		fmt.Println(gchalk.Red("Onii-chan! anata wa need repository!"))
		fmt.Println(gchalk.Blue("Example: gigit nazhard/gigit"))
	}
}

func fetch(name, commit string) (string, error) {
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

	f, _ := os.Create(file_name + ".tar.gz")

	defer f.Close()

	_, err = f.ReadFrom(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode != http.StatusOK {
		return url, errors.New(gchalk.Red("Error " + name + " not found"))
	}

	return url, nil
}

// get latest commit
func getLatestCommit() string {
	url := "https://api.github.com/repos/nazhard/do/commits"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	var commits []Commit

	err = json.Unmarshal(data, &commits)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	var c string
	if len(commits) > 1 {
		first := commits[0]
		c = first.Sha
	}

	return c
}

// extract downloaded repository
func extractGz(in, out, path string) {
	file, _ := os.Open(in)
	var shift = func(path string) string {
		parts := strings.Split(path, string(filepath.Separator))
		parts = parts[1:]
		return strings.Join(parts, string(filepath.Separator))
	}
	extract.Gz(context.TODO(), file, out, shift)
}

// this will be useful for private repositories
func clone() {
	if len(os.Args) >= 2 {
		if strings.Contains(os.Args[1], "-1") {
			user_repo = os.Args[2]

			cmd := exec.Command("git", "clone", "--depth=1", "https://github.com/"+user_repo+".git")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			_ = cmd.Start()

			defer cmd.Wait()
		} else {
			user_repo = os.Args[1]

			cmd := exec.Command("git", "clone", "https://github.com/"+user_repo+".git")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			_ = cmd.Start()

			defer cmd.Wait()
		}
	}
}
