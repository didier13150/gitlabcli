package gitlabcli

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type GitlabEnv struct {
	UrlBase    string
	Token      string
	Verbose    bool
	DryrunMode bool
	ProjectId  string
	Glapi      GLApi
	GitlabData []GitlabEnvData
	FileData   []GitlabEnvData
}

func NewGitlabEnv(UrlBase string, Token string, Verbose bool) GitlabEnv {
	glenv := GitlabEnv{}
	glenv.UrlBase = UrlBase
	glenv.Token = Token
	glenv.Verbose = Verbose
	glenv.DryrunMode = false
	glenv.Glapi = NewGLApi(UrlBase, Token, Verbose)
	glenv.GitlabData = []GitlabEnvData{}
	glenv.FileData = []GitlabEnvData{}
	return glenv
}

func (glenv *GitlabEnv) GetEnvsFromGitlab() error {
	var envsByPage = []GitlabEnvData{}

	currentPage := 1
	nbPage := 1
	perPage := 50
	var err error
	var resp []byte

	for currentPage <= nbPage {
		resp, nbPage, err = glenv.Glapi.CallGitlabApi("/api/v4/projects/"+glenv.ProjectId+"/environments?per_page="+strconv.Itoa(perPage)+"&page="+strconv.Itoa(currentPage), http.MethodGet, nil)
		if err != nil {
			log.Fatalln(err)
			return err
		}
		err := json.Unmarshal([]byte(resp), &envsByPage)
		if err != nil {
			log.Println("Decoding env json", err)
			continue
		}
		currentPage++
		glenv.GitlabData = append(glenv.GitlabData, envsByPage...)
	}
	return nil
}

func (glenv *GitlabEnv) ExportEnvs(filename string) {
	json, err := json.MarshalIndent(glenv.GitlabData, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(filename, json, glenv.Verbose)
}

func (glenv *GitlabEnv) ImportEnvs(filename string) {
	resp := readFromFile(filename, "envs", glenv.Verbose)
	err := json.Unmarshal([]byte(resp), &glenv.FileData)
	if err != nil {
		log.Println("Importing env json", err)
	}
}

func (glenv *GitlabEnv) InsertEnv(env GitlabEnvData) error {
	if glenv.DryrunMode {
		log.Print("Cannot inserting env in dryrun mode")
	}
	urlapi := "/api/v4/projects/" + glenv.ProjectId + "/environments"
	log.Printf("Use URL %s to insert env", urlapi)
	json, err := json.Marshal(env)
	if err != nil {
		log.Println(err)
		return err
	}
	if glenv.Verbose {
		log.Println(string(json))
	}
	log.Printf("Insert env %s (%d)", env.Name, env.Id)
	resp, _, err := glenv.Glapi.CallGitlabApi(urlapi, http.MethodPost, json)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if glenv.Verbose {
		log.Println(string(resp))
	}
	return nil
}

func (glenv *GitlabEnv) UpdateEnv(env GitlabEnvData) error {
	if glenv.DryrunMode {
		log.Print("Cannot updating env in dryrun mode")
	}
	urlapi := "/api/v4/projects/" + glenv.ProjectId + "/environments/" + strconv.Itoa(env.Id)
	log.Printf("Use URL %s to update env", urlapi)
	json, err := json.Marshal(env)
	if err != nil {
		log.Println(err)
		return err
	}
	if glenv.Verbose {
		log.Println(string(json))
	}
	log.Printf("Update env %s (%d)", env.Name, env.Id)
	resp, _, err := glenv.Glapi.CallGitlabApi(urlapi, http.MethodPut, json)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if glenv.Verbose {
		log.Println(string(resp))
	}
	return nil
}

func (glenv *GitlabEnv) DeleteEnv(env GitlabEnvData) error {
	if glenv.DryrunMode {
		log.Print("Cannot deleting env in dryrun mode")
	}
	urlapi := "/api/v4/projects/" + glenv.ProjectId + "/environments/" + strconv.Itoa(env.Id)
	log.Printf("Use URL %s to delete env", urlapi)
	log.Printf("Stop env %s (%d)", env.Name, env.Id)
	resp, _, err := glenv.Glapi.CallGitlabApi(urlapi+"/stop", http.MethodPost, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if glenv.Verbose {
		log.Println(string(resp))
	}
	log.Printf("Delete env %s (%d)", env.Name, env.Id)
	resp, _, err = glenv.Glapi.CallGitlabApi(urlapi, http.MethodDelete, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if glenv.Verbose {
		log.Println(string(resp))
	}
	return nil
}

func (glenv *GitlabEnv) GetMissingEnvs(envsFromVars []string) []string {
	var envs []string
	var found bool
	for _, envNeeded := range envsFromVars {
		if envNeeded == "*" {
			continue
		}
		found = false
		for _, env := range glenv.GitlabData {
			if env.Name == envNeeded {
				found = true
				break
			}
		}
		if !found {
			log.Printf("Env %s should be added", envNeeded)
			envs = append(envs, envNeeded)
		}
	}
	return envs
}

func (glenv *GitlabEnv) CompareEnv() ([]GitlabEnvData, []GitlabEnvData, []GitlabEnvData) {
	var found bool
	var envsToAdd []GitlabEnvData
	var envsToDelete []GitlabEnvData
	var envsToUpdate []GitlabEnvData
	var missingKey = make(map[string]GitlabEnvData)

	for _, item1 := range glenv.GitlabData {
		found = false
		for _, item2 := range glenv.FileData {
			if item1.Name == item2.Name && item1.Description == item2.Description && item1.Url == item2.Url {
				found = true
				break
			}
		}
		if !found {
			missingKey[item1.Name] = item1
		}
	}
	for _, item2 := range glenv.FileData {
		found = false
		for _, item1 := range glenv.GitlabData {
			if item1.Name == item2.Name && item1.Description == item2.Description && item1.Url == item2.Url {
				found = true
				break
			}
		}
		if !found {
			_, itemExists := missingKey[item2.Name]
			if itemExists {
				delete(missingKey, item2.Name)
				envsToUpdate = append(envsToUpdate, item2)
				log.Printf("Env %s should be updated", item2.Name)
			} else {
				envsToAdd = append(envsToAdd, item2)
				log.Printf("Env %s should be added", item2.Name)
			}
		}
	}
	for _, item := range missingKey {
		envsToDelete = append(envsToDelete, item)
		log.Printf("Env %s should be deleted", item.Name)
	}
	return envsToAdd, envsToDelete, envsToUpdate
}
