package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type GitlabProject struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Description       any    `json:"description"`
	Path              string `json:"path"`
	NameWithNamespace string `json:"name_with_namespace"`
	PathWithNamespace string `json:"path_with_namespace"`
	SshUrlToRepo      string `json:"ssh_url_to_repo"`
	HttpUrlToRepo     string `json:"http_url_to_repo"`
	WebUrl            string `json:"web_url"`
	Visibility        string `json:"visibility"`
}

func getProjectsFromGitlab(gitlabUrl string, token string) ([]GitlabProject, error) {
	var projects = []GitlabProject{}
	var projectsByPage = []GitlabProject{}

	currentPage := 1
	nbPage := 1
	perPage := 50
	var err error
	var resp []byte

	for currentPage <= nbPage {
		resp, nbPage, err = callGitlabApi(gitlabUrl+"/api/v4/projects?per_page="+strconv.Itoa(perPage)+"&page="+strconv.Itoa(currentPage), token, http.MethodGet, nil)
		if err != nil {
			log.Fatalln(err)
			return []GitlabProject{}, err
		}
		err := json.Unmarshal([]byte(resp), &projectsByPage)
		if err != nil {
			log.Println("Decoding project json", err)
			continue
		}
		currentPage++
		projects = append(projects, projectsByPage...)
	}
	return projects, nil
}

func exportProjects(filename string, projects []GitlabProject) {
	json, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(filename, json)
}

func importProjects(filename string) []GitlabProject {
	var result []GitlabProject
	resp := readFromFile(filename, "projects")
	err := json.Unmarshal([]byte(resp), &result)
	if err != nil {
		log.Println("Importing projects json", err)
	}
	return result
}

func getProjectIdByRepoUrl(url string, projects []GitlabProject) int {
	for _, project := range projects {
		if project.SshUrlToRepo == url || project.HttpUrlToRepo == url {
			return project.Id
		}
	}
	return 0
}
