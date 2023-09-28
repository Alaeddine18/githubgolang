package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func cloneARepo(urls []string) {
	baseDestinationPath := "./reposCloned"

	for _, url := range urls {
		// Extract the repository name from the URL
		repoName := getRepoNameFromURL(url)

		// Create a subdirectory for each repository
		destinationPath := fmt.Sprintf("%s/%s", baseDestinationPath, repoName)

		// Check if the destination directory already exists
		if _, err := os.Stat(destinationPath); !os.IsNotExist(err) {
			fmt.Printf("Repo %s:  already exists. We Skip \n", repoName)
			continue
		}

		cmd := exec.Command("git", "clone", url, destinationPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to clone %s: %v\n", url, err)
			continue
		}

		fmt.Printf("Cloned successful %s.\n", repoName)
	}
}

func getRepoNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}
