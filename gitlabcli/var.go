package gitlabcli

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type GitlabVar struct {
	UrlBase    string
	Token      string
	Verbose    bool
	DryrunMode bool
	ProjectId  string
	Glapi      GLApi
	GitlabData []GitlabVarData
	FileData   []GitlabVarData
}

func NewGitlabVar(UrlBase string, Token string, Verbose bool) GitlabVar {
	glvar := GitlabVar{}
	glvar.UrlBase = UrlBase
	glvar.Token = Token
	glvar.Verbose = Verbose
	glvar.DryrunMode = false
	glvar.Glapi = NewGLApi(UrlBase, Token, Verbose)
	glvar.GitlabData = []GitlabVarData{}
	glvar.FileData = []GitlabVarData{}
	return glvar
}

func (glvar *GitlabVar) GetVarsFromGitlab() error {
	var varsByPage = []GitlabVarData{}

	currentPage := 1
	nbPage := 1
	perPage := 50
	var err error
	var resp []byte

	for currentPage <= nbPage {
		resp, nbPage, err = glvar.Glapi.CallGitlabApi("/api/v4/projects/"+glvar.ProjectId+"/variables?per_page="+strconv.Itoa(perPage)+"&page="+strconv.Itoa(currentPage), http.MethodGet, nil)
		if err != nil {
			log.Fatalln(err)
			return err
		}
		err := json.Unmarshal([]byte(resp), &varsByPage)
		if err != nil {
			log.Println("Decoding var json", err)
			continue
		}
		currentPage++
		glvar.GitlabData = append(glvar.GitlabData, varsByPage...)
	}
	return nil
}

func (glvar *GitlabVar) ExportVars(filename string) {
	json, err := json.MarshalIndent(glvar.GitlabData, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(filename, json, glvar.Verbose)
}

func (glvar *GitlabVar) ImportVars(filename string) {
	resp := readFromFile(filename, "vars", glvar.Verbose)
	err := json.Unmarshal([]byte(resp), &glvar.FileData)
	if err != nil {
		log.Println("Importing var json", err)
	}
}

func (glvar *GitlabVar) InsertVar(variable GitlabVarData) error {
	if glvar.DryrunMode {
		log.Print("Cannot inserting var in dryrun mode")
	}
	urlapi := "/api/v4/projects/" + glvar.ProjectId + "/variables"
	log.Printf("Use URL %s to insert var", urlapi)
	json, err := json.Marshal(variable)
	if err != nil {
		log.Println(err)
		return err
	}
	if glvar.Verbose {
		log.Println(string(json))
	}
	log.Printf("Insert var %s in %s env", variable.Key, variable.Env)
	resp, _, err := glvar.Glapi.CallGitlabApi(urlapi, http.MethodPost, json)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if glvar.Verbose {
		log.Println(string(resp))
	}
	return nil
}

func (glvar *GitlabVar) UpdateVar(variable GitlabVarData) error {
	if glvar.DryrunMode {
		log.Print("Cannot updating var in dryrun mode")
	}
	urlapi := "/api/v4/projects/" + glvar.ProjectId + "/variables/" + variable.Key + "?filter[environment_scope]=" + variable.Env
	log.Printf("Use URL %s to update var", urlapi)
	json, err := json.Marshal(variable)
	if err != nil {
		log.Println(err)
		return err
	}
	if glvar.Verbose {
		log.Println(string(json))
	}
	log.Printf("Update var %s in %s env", variable.Key, variable.Env)
	resp, _, err := glvar.Glapi.CallGitlabApi(urlapi, http.MethodPut, json)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if glvar.Verbose {
		log.Println(string(resp))
	}
	return nil
}

func (glvar *GitlabVar) DeleteVar(variable GitlabVarData) error {
	if glvar.DryrunMode {
		log.Print("Cannot deleting var in dryrun mode")
	}
	urlapi := "/api/v4/projects/" + glvar.ProjectId + "/variables/" + variable.Key + "?filter[environment_scope]=" + variable.Env
	log.Printf("Use URL %s to delete var", urlapi)
	log.Printf("Delete var %s in %s env", variable.Key, variable.Env)
	resp, _, err := glvar.Glapi.CallGitlabApi(urlapi, http.MethodDelete, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if glvar.Verbose {
		log.Println(string(resp))
	}
	return nil
}

func (glvar *GitlabVar) GetEnvsFromVars() []string {
	var envs []string
	var found bool
	for _, item := range glvar.FileData {
		found = false
		for _, env := range envs {
			if env == item.Env {
				found = true
				break
			}
		}
		if !found {
			envs = append(envs, item.Env)
		}
	}
	return envs
}

func (glvar *GitlabVar) CompareVar() ([]GitlabVarData, []GitlabVarData, []GitlabVarData) {
	var found bool
	var varsToAdd []GitlabVarData
	var varsToDelete []GitlabVarData
	var varsToUpdate []GitlabVarData
	var missingKey = make(map[string]GitlabVarData)

	for _, item1 := range glvar.GitlabData {
		found = false
		for _, item2 := range glvar.FileData {
			if item1 == item2 {
				found = true
				break
			}
		}
		if !found {
			missingKey[item1.Env+"/"+item1.Key] = item1
		}
	}
	for _, item2 := range glvar.FileData {
		found = false
		for _, item1 := range glvar.GitlabData {
			if item1 == item2 {
				found = true
				break
			}
		}
		if !found {
			_, itemExists := missingKey[item2.Env+"/"+item2.Key]
			if itemExists {
				delete(missingKey, item2.Env+"/"+item2.Key)
				varsToUpdate = append(varsToUpdate, item2)
				log.Printf("Var %s[%s] should be updated", item2.Key, item2.Env)
			} else {
				varsToAdd = append(varsToAdd, item2)
				log.Printf("Var %s[%s] should be added", item2.Key, item2.Env)
			}
		}
	}
	for _, item := range missingKey {
		varsToDelete = append(varsToDelete, item)
		log.Printf("Var %s[%s] should be deleted", item.Key, item.Env)
	}
	return varsToAdd, varsToDelete, varsToUpdate
}
