package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type GitlabVar struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description any    `json:"description"`
	Env         string `json:"environment_scope"`
	IsRaw       bool   `json:"raw"`
	IsHidden    bool   `json:"hidden"`
	IsProtected bool   `json:"protected"`
	IsMasked    bool   `json:"masked"`
}

func getVarsFromGitlab(gitlabUrl string, token string, projectId string) ([]GitlabVar, error) {
	var vars = []GitlabVar{}
	var varsByPage = []GitlabVar{}

	currentPage := 1
	nbPage := 1
	perPage := 50
	var err error
	var resp []byte

	for currentPage <= nbPage {
		resp, nbPage, err = callGitlabApi(gitlabUrl+"/api/v4/projects/"+projectId+"/variables?per_page="+strconv.Itoa(perPage)+"&page="+strconv.Itoa(currentPage), token, http.MethodGet, nil)
		if err != nil {
			log.Fatalln(err)
			return []GitlabVar{}, err
		}
		err := json.Unmarshal([]byte(resp), &varsByPage)
		if err != nil {
			log.Println("Decoding var json", err)
			continue
		}
		currentPage++
		vars = append(vars, varsByPage...)
	}
	return vars, nil
}

func exportVars(filename string, vars []GitlabVar) {
	json, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(filename, json)
}

func importVars(filename string) []GitlabVar {
	var result []GitlabVar
	resp := readFromFile(filename, "vars")
	err := json.Unmarshal([]byte(resp), &result)
	if err != nil {
		log.Println("Importing var json", err)
	}
	return result
}

func compareVarSet(set1 []GitlabVar, set2 []GitlabVar) ([]GitlabVar, []GitlabVar, []GitlabVar) {
	var found bool
	var varsToAdd []GitlabVar
	var varsToDelete []GitlabVar
	var varsToUpdate []GitlabVar
	var missingKey = make(map[string]GitlabVar)

	for _, item1 := range set1 {
		found = false
		for _, item2 := range set2 {
			if item1 == item2 {
				found = true
				break
			}
		}
		if !found {
			missingKey[item1.Env+"/"+item1.Key] = item1
		}
	}
	for _, item2 := range set2 {
		found = false
		for _, item1 := range set1 {
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
				log.Println("Var", item2, "should be updated")
			} else {
				varsToAdd = append(varsToAdd, item2)
				log.Println("Var", item2, "should be added")
			}
		}
	}
	for _, item := range missingKey {
		varsToDelete = append(varsToDelete, item)
		log.Println("Var", item, "should be deleted")
	}
	return varsToAdd, varsToDelete, varsToUpdate
}

func insertVar(url string, token string, projectId string, variable GitlabVar) error {
	urlapi := url + "/api/v4/projects/" + projectId + "/variables"
	log.Printf("Use URL %s to insert var", urlapi)
	json, err := json.Marshal(variable)
	if err != nil {
		log.Println(err)
		return err
	}
	if verbose {
		log.Println(string(json))
	}
	log.Printf("Insert var %s in %s env", variable.Key, variable.Env)
	resp, _, err := callGitlabApi(urlapi, token, http.MethodPost, json)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if verbose {
		log.Println(string(resp))
	}
	return nil
}

func updateVar(gitlabUrl string, token string, projectId string, variable GitlabVar) error {
	urlapi := gitlabUrl + "/api/v4/projects/" + projectId + "/variables/" + variable.Key + "?filter[environment_scope]=" + variable.Env
	log.Printf("Use URL %s to update var", urlapi)
	json, err := json.Marshal(variable)
	if err != nil {
		log.Println(err)
		return err
	}
	if verbose {
		log.Println(string(json))
	}
	log.Printf("Update var %s in %s env", variable.Key, variable.Env)
	resp, _, err := callGitlabApi(urlapi, token, http.MethodPut, json)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if verbose {
		log.Println(string(resp))
	}
	return nil
}

func deleteVar(gitlabUrl string, token string, projectId string, variable GitlabVar) error {
	urlapi := gitlabUrl + "/api/v4/projects/" + projectId + "/variables/" + variable.Key + "?filter[environment_scope]=" + variable.Env
	log.Printf("Use URL %s to delete var", urlapi)
	log.Printf("Delete var %s in %s env", variable.Key, variable.Env)
	resp, _, err := callGitlabApi(urlapi, token, http.MethodDelete, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if verbose {
		log.Println(string(resp))
	}
	return nil
}

func getEnvsFromVars(vars []GitlabVar) []string {
	var envs []string
	var found bool
	for _, item := range vars {
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
