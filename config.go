package main

import (
	"log"
	"os"
)

const (
	githubToken = "GITHUB_TOKEN"
	userName    = "GITHUB_USERNAME"
	typeOf      = "GITHUB_TYPEOF"
)

type Config struct {
	GithubToken string
	UserName    string
	TypeOf      string
}

func NewConfig() Config {
	githubtoken, ok := os.LookupEnv(githubToken)
	if !ok || githubtoken == "" {
		log.Fatal(githubToken)
	}

	username, ok := os.LookupEnv(userName)
	if !ok || username == "" {
		log.Fatal(userName)
	}

	typeof, ok := os.LookupEnv(typeOf)
	if !ok || typeof == "" {
		log.Fatal(typeOf)
	}
	return Config{
		GithubToken: githubtoken,
		UserName:    username,
		TypeOf:      typeof,
	}
}
