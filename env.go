package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// type GitlabEnvNamespace struct {
//	Id int `json:"id"`
//}

// type GitlabEnvProject struct {
//	Id        int                `json:"id"`
//	Namespace GitlabEnvNamespace `json:"namespace"`
//}

type GitlabEnv struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	State       any    `json:"state"`
	Url         any    `json:"external_url"`
	Description any    `json:"description"`
	//Project     GitlabEnvProject `json:"project"`
}

func getEnvsFromGitlab(gitlabUrl string, token string, projectId string) ([]GitlabEnv, error) {
	var envs = []GitlabEnv{}
	var envsByPage = []GitlabEnv{}

	currentPage := 1
	nbPage := 1
	perPage := 50
	var err error
	var resp []byte

	for currentPage <= nbPage {
		resp, nbPage, err = callGitlabApi(gitlabUrl+"/api/v4/projects/"+projectId+"/environments?per_page="+strconv.Itoa(perPage)+"&page="+strconv.Itoa(currentPage), token, http.MethodGet, nil)
		if err != nil {
			log.Fatalln(err)
			return []GitlabEnv{}, err
		}
		err := json.Unmarshal([]byte(resp), &envsByPage)
		if err != nil {
			log.Println("Decoding env json", err)
			continue
		}
		currentPage++
		envs = append(envs, envsByPage...)
	}
	return envs, nil
}

func exportEnvs(filename string, envs []GitlabEnv) {
	json, err := json.MarshalIndent(envs, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(filename, json)
}

func importEnvs(filename string) []GitlabEnv {
	var result []GitlabEnv
	resp := readFromFile(filename, "envs")
	err := json.Unmarshal([]byte(resp), &result)
	if err != nil {
		log.Println("Importing env json", err)
	}
	return result
}

func compareEnvSet(set1 []GitlabEnv, set2 []GitlabEnv) ([]GitlabEnv, []GitlabEnv, []GitlabEnv) {
	var found bool
	var envsToAdd []GitlabEnv
	var envsToDelete []GitlabEnv
	var envsToUpdate []GitlabEnv
	var missingKey = make(map[string]GitlabEnv)

	for _, item1 := range set1 {
		found = false
		for _, item2 := range set2 {
			if item1.Name == item2.Name && item1.Description == item2.Description && item1.Url == item2.Url {
				found = true
				break
			}
		}
		if !found {
			missingKey[item1.Name] = item1
		}
	}
	for _, item2 := range set2 {
		found = false
		for _, item1 := range set1 {
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
				log.Println("Env", item2, "should be updated")
			} else {
				envsToAdd = append(envsToAdd, item2)
				log.Println("Env", item2, "should be added")
			}
		}
	}
	for _, item := range missingKey {
		envsToDelete = append(envsToDelete, item)
		log.Println("Env", item, "should be deleted")
	}
	return envsToAdd, envsToDelete, envsToUpdate
}

func insertEnv(gitlabUrl string, token string, projectId string, env GitlabEnv) error {
	urlapi := gitlabUrl + "/api/v4/projects/" + projectId + "/environments"
	log.Printf("Use URL %s to insert env", urlapi)
	json, err := json.Marshal(env)
	if err != nil {
		log.Println(err)
		return err
	}
	if verbose {
		log.Println(string(json))
	}
	log.Printf("Insert env %s (%d)", env.Name, env.Id)
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

func updateEnv(gitlabUrl string, token string, projectId string, env GitlabEnv) error {
	urlapi := gitlabUrl + "/api/v4/projects/" + projectId + "/environments/" + strconv.Itoa(env.Id)
	log.Printf("Use URL %s to update env", urlapi)
	json, err := json.Marshal(env)
	if err != nil {
		log.Println(err)
		return err
	}
	if verbose {
		log.Println(string(json))
	}
	log.Printf("Update env %s (%d)", env.Name, env.Id)
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

func deleteEnv(gitlabUrl string, token string, projectId string, env GitlabEnv) error {
	urlapi := gitlabUrl + "/api/v4/projects/" + projectId + "/environments/" + strconv.Itoa(env.Id)
	log.Printf("Use URL %s to delete env", urlapi)
	log.Printf("Stop env %s (%d)", env.Name, env.Id)
	resp, _, err := callGitlabApi(urlapi+"/stop", token, http.MethodPost, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if verbose {
		log.Println(string(resp))
	}
	log.Printf("Delete env %s (%d)", env.Name, env.Id)
	resp, _, err = callGitlabApi(urlapi, token, http.MethodDelete, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if verbose {
		log.Println(string(resp))
	}
	return nil
}

func getMissingEnvs(envsFromVars []string, envsFromGitlab []GitlabEnv) []string {
	var envs []string
	var found bool
	for _, envNeeded := range envsFromVars {
		if envNeeded == "*" {
			continue
		}
		found = false
		for _, env := range envsFromGitlab {
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
