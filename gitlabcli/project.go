package gitlabcli

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type GitlabProjectData struct {
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

type GitlabProject struct {
	UrlBase string
	Token   string
	Verbose bool
	Glapi   GLApi
	Data    []GitlabProjectData
}

func NewGitlabProject(UrlBase string, Token string, Verbose bool) GitlabProject {
	glproj := GitlabProject{}
	glproj.UrlBase = UrlBase
	glproj.Token = Token
	glproj.Verbose = Verbose
	glproj.Data = []GitlabProjectData{}
	glproj.Glapi = NewGLApi(UrlBase, Token, Verbose)
	return glproj
}

func (glproj *GitlabProject) GetProjectsFromGitlab() error {
	var projectsByPage = []GitlabProjectData{}

	currentPage := 1
	nbPage := 1
	perPage := 50
	var err error
	var resp []byte

	for currentPage <= nbPage {
		resp, nbPage, err = glproj.Glapi.CallGitlabApi("/api/v4/projects?per_page="+strconv.Itoa(perPage)+"&page="+strconv.Itoa(currentPage), http.MethodGet, nil)
		if err != nil {
			log.Fatalln(err)
			return err
		}
		err := json.Unmarshal([]byte(resp), &projectsByPage)
		if err != nil {
			log.Println("Decoding project json", err)
			continue
		}
		currentPage++
		glproj.Data = append(glproj.Data, projectsByPage...)
	}
	return nil
}

func (glproj *GitlabProject) ExportProjects(filename string) {
	json, err := json.MarshalIndent(glproj.Data, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(filename, json, glproj.Verbose)
}

func (glproj *GitlabProject) ImportProjects(filename string) {
	resp := readFromFile(filename, "projects", glproj.Verbose)
	err := json.Unmarshal([]byte(resp), &glproj.Data)
	if err != nil {
		log.Println("Importing projects json", err)
	}
}

func (glproj *GitlabProject) GetProjectIdByRepoUrl(url string) int {
	for _, project := range glproj.Data {
		if project.SshUrlToRepo == url || project.HttpUrlToRepo == url {
			return project.Id
		}
	}
	return 0
}
