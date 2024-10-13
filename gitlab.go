package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/0x4f53/textsubs"
)

const baseURL = "https://gitlab.com/api/v4"

var perPage *int
var maxCommits *int
var jsonOutput *bool

type Namespace struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Namespace Namespace `json:"namespace"`
	WebURL    string    `json:"web_url"`
}

type Commit struct {
	ID         string    `json:"id"`
	ShortID    string    `json:"short_id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	AuthorName string    `json:"author_name"`
	PatchURL   string    `json:"patch_url"`
}

type MergedOutput struct {
	ProjectID         int       `json:"project_id"`
	ProjectName       string    `json:"project_name"`
	ProjectPath       string    `json:"project_path"`
	ProjectNamespace  string    `json:"project_namespace"`
	ProjectWebURL     string    `json:"project_web_url"`
	Kind              string    `json:"kind"`
	CommitID          string    `json:"commit_id"`
	CommitShortID     string    `json:"commit_short_id"`
	CommitTitle       string    `json:"commit_title"`
	CommitMessage     string    `json:"commit_message"`
	CommitCreatedAt   time.Time `json:"commit_created_at"`
	CommitAuthorName  string    `json:"commit_author_name"`
	CommitPatchURL    string    `json:"commit_patch_url"`
	AssociatedDomains []string  `json:"associated_domains"`
}

func getProjects(perPage int, page int) ([]Project, error) {
	url := fmt.Sprintf("%s/projects?visibility=public&per_page=%d&page=%d", baseURL, perPage, page)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch projects: %s", resp.Status)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func getCommits(perPage int, projectID int) ([]Commit, error) {
	url := fmt.Sprintf("%s/projects/%d/repository/commits?per_page=%d", baseURL, projectID, perPage)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch commits for project %d: %s", projectID, resp.Status)
	}

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	for i := range commits {
		commits[i].PatchURL = fmt.Sprintf("%s/projects/%d/repository/commits/%s.patch", baseURL, projectID, commits[i].ID)
	}
	return commits, nil
}

func curl(url string) string {
	response, err := http.Get(url)
	if err != nil {
		//log.Fatalf("Error fetching the URL: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		//log.Fatalf("Error: received status code %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//log.Fatalf("Error reading response body: %v", err)
	}

	return string(body)
}

func GetGitlabCommits(perPage int, maxCommits int) []MergedOutput {
	var allMergedOutputs []MergedOutput

	page := 1
	for len(allMergedOutputs) < maxCommits {
		projects, err := getProjects(perPage, page)
		if err != nil {
			log.Fatalf("Error fetching projects: %v", err)
		}
		if len(projects) == 0 {
			break
		}

		for _, project := range projects {
			commits, err := getCommits(perPage, project.ID)
			if err != nil {
				log.Printf("Error fetching commits for project %d: %v", project.ID, err)
				continue
			}

			for _, commit := range commits {
				commit.PatchURL = fmt.Sprintf("%s/-/commit/%s.patch", project.WebURL, commit.ID)

				// Read patch data
				commitContents := curl(commit.PatchURL)
				associatedDomains, _ := textsubs.DomainsOnly(commitContents, false)
				associatedDomains = textsubs.Resolve(associatedDomains)

				data := MergedOutput{
					ProjectID:         project.ID,
					ProjectName:       project.Name,
					ProjectPath:       project.Path,
					ProjectNamespace:  project.Namespace.Name,
					ProjectWebURL:     project.WebURL,
					CommitID:          commit.ID,
					Kind:              project.Namespace.Kind,
					CommitShortID:     commit.ShortID,
					CommitTitle:       commit.Title,
					CommitMessage:     commit.Message,
					CommitCreatedAt:   commit.CreatedAt,
					CommitAuthorName:  commit.AuthorName,
					CommitPatchURL:    commit.PatchURL,
					AssociatedDomains: associatedDomains,
				}

				allMergedOutputs = append(allMergedOutputs, data)

				if *jsonOutput {
					output, err := json.Marshal(data)
					if err != nil {
						log.Fatalf("Error marshalling JSON: %v", err)
					}
					appendToFile(timestamp()+".json", string(output)+"\n")
					fmt.Println(string(output))
				} else {
					output, err := json.MarshalIndent(data, "", "  ")
					if err != nil {
						log.Fatalf("Error marshalling JSON: %v", err)
					}
					fmt.Println(string(output))
				}

				if len(allMergedOutputs) >= maxCommits {
					break
				}
			}

			if len(allMergedOutputs) >= maxCommits {
				break
			}
		}

		page++
	}

	return allMergedOutputs
}

func timestamp() string {
	now := time.Now()
	return fmt.Sprintf("%02d-%02d-%04d-%d-%d", now.Month(), now.Day(), now.Year(), now.Hour(), now.Minute())
}

func appendToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		return err
	}

	return nil
}

func main() {
	perPage = flag.Int("per", 100, "Results to grab per page from GitLab API (max: 100)")
	maxCommits = flag.Int("max", 100, "Maximum number of commits to grab from GitLab API (max: 1000)")
	jsonOutput = flag.Bool("json", false, "Save as line-separated JSON in a file (filename format: <01-01-2024-0>.json)")

	flag.Parse()

	_ = GetGitlabCommits(*perPage, *maxCommits)
}
