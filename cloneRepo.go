package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func cloneARepo(urls []string) {
	baseDestinationPath := "./reposCloned"

	for _, url := range urls {
		repoName := getRepoNameFromURL(url)

		destinationPath := fmt.Sprintf("%s/%s", baseDestinationPath, repoName)

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

	createZipFile("./reposCloned.zip", baseDestinationPath)
}

func createZipFile(zipFilePath, sourceDir string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(filePath) == ".zip" {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		zipHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		zipHeader.Name, err = filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		zipEntry, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getRepoNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}
