package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Result struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`

	Description   string    `json:"description"`
	Url           string    `json:"html_url"`
	Created_at    time.Time `json:"created_at"`
	Udated_at     time.Time `json:"updated_at"`
	Fork          bool      `json:"fork"`
	Allow_forking bool      `json:"allow_forking"`
	Open_issues   int       `json:"open_issues"`
}

func main() {

	app := fiber.New()

	config := NewConfig()

	allReposforUsersOrOrgs := getReposForUserOrOrgs(config)

	fmt.Println("User or Orgs: ", config.UserName)
	fmt.Println("Type of: ", config.TypeOf)
	fmt.Println("Number of repos: ", len(allReposforUsersOrOrgs))
	fmt.Println("List of repos: ", allReposforUsersOrOrgs)

	url2Clone := Csvwriter(allReposforUsersOrOrgs)

	cloneARepo(url2Clone)

	app.Get("/download", downloadRepository)

	defer log.Fatal(app.Listen(":3000"))

}

// it make a get on one of theses ( based on typeOf ): https://api.github.com/users/username/repos ou https://api.github.com/orgs/OrgsName/repos
func getReposForUserOrOrgs(config Config) []Result {

	url := fmt.Sprintf("https://api.github.com/%s/%s/repos", config.TypeOf, config.UserName)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resultBodyParsed := parseResponse(string(resultBody))

	return resultBodyParsed
}

func parseResponse(body string) []Result {
	var result []Result
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Udated_at.After(result[j].Udated_at)
	})
	return result
}

func Csvwriter(response []Result) []string {

	// for the clone part :)
	var urls []string

	file, err := os.Create("repos.csv")
	if err != nil {
		log.Fatal("failed to create the file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	header := []string{
		"Name",
		"Owner",
		"Description",
		"Url",
		"Created_at",
		"Updated_at",
		"Fork",
		"Allow_forking",
		"Open_issues",
	}
	if err := w.Write(header); err != nil {
		log.Fatal("error writing collumn", err)
	}

	for _, record := range response {
		row := []string{
			record.Name,
			record.Owner.Login,
			record.Description,
			record.Url,
			record.Created_at.String(),
			record.Udated_at.String(),
			strconv.FormatBool(record.Fork),
			strconv.FormatBool(record.Allow_forking),
			strconv.Itoa(record.Open_issues),
		}
		if err := w.Write(row); err != nil {
			log.Fatal("error writing the file", err)
		}
		// Append the URL to the slice
		urls = append(urls, record.Url)
	}
	// Return the slice of URLs
	return urls
}
