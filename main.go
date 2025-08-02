package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

var verbose bool
var debug bool

func main() {

	var projectId string
	var varsFile string
	var envsFile string
	var projectsFile string
	var gitlabUrl string
	var gitlabTokenFile string
	var token string

	defaultUrl := "https://gitlab.com"
	defaultIdFile := ".gitlab.id"
	defaultVarFile := ".gitlab-vars.json"
	defaultEnvFile := ".gitlab-envs.json"
	defaultProjectsFile := os.Getenv("HOME") + "/.gitlab-projects.json"
	defaultTokenFile := os.Getenv("HOME") + "/.gitlab.token"
	defaultDebugFile := "debug.txt"

	if len(os.Getenv("GLVARS_GITLAB_URL")) > 0 {
		defaultUrl = os.Getenv("GLVARS_GITLAB_URL")
	}
	if len(os.Getenv("GLVARS_TOKEN_FILE")) > 0 {
		defaultTokenFile = os.Getenv("GLVARS_TOKEN_FILE")
	}
	if len(os.Getenv("GLVARS_VAR_FILE")) > 0 {
		defaultVarFile = os.Getenv("GLVARS_VAR_FILE")
	}
	if len(os.Getenv("GLVARS_ENV_FILE")) > 0 {
		defaultEnvFile = os.Getenv("GLVARS_ENV_FILE")
	}
	if len(os.Getenv("GLVARS_PROJECT_FILE")) > 0 {
		defaultProjectsFile = os.Getenv("GLVARS_PROJECT_FILE")
	}
	if len(os.Getenv("GLVARS_ID_FILE")) > 0 {
		defaultIdFile = os.Getenv("GLVARS_ID_FILE")
	}
	if len(os.Getenv("GLVARS_DEBUG_FILE")) > 0 {
		defaultDebugFile = os.Getenv("GLVARS_DEBUG_FILE")
	}
	flag.StringVar(&projectId, "id", "", "Gitlab project identifiant.")
	flag.StringVar(&varsFile, "varfile", defaultVarFile, "File which contains vars.")
	flag.StringVar(&envsFile, "envfile", defaultEnvFile, "File which contains envs.")
	flag.StringVar(&projectsFile, "projectfile", defaultProjectsFile, "File which contains projects.")
	flag.StringVar(&gitlabUrl, "url", defaultUrl, "Gitlab URL.")
	flag.StringVar(&gitlabTokenFile, "token", defaultTokenFile, "File which contains token to access Gitlab API.")
	flag.BoolVar(&verbose, "verbose", false, "Make application more talkative.")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	var dryrun = flag.Bool("dryrun", false, "Run in dry-run mode (read only).")
	var export = flag.Bool("export", false, "Export current variables in var file.")
	var exportProjectsOnly = flag.Bool("export-projects", false, "Export current projects in project file.")
	var deleteIsActive = flag.Bool("delete", false, "Delete Gitlab var if not present in var file.")

	flag.Usage = func() {
		fmt.Printf("Usage: " + os.Args[0] + " [--id <Poject ID>] [--varfile <VAR FILE>] [--envfile <ENV FILE>] [--projectfile <PROJECT FILE>] [--token <TOKEN FILE>] [--dryrun] [--export] [--export-projects] [--delete]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if verbose {
		log.Print("Verbose mode is active")
	}
	if *dryrun {
		log.Print("Dry run mode is active")
	}
	if *export {
		log.Print("Export requested")
	}
	if *exportProjectsOnly {
		log.Print("Export projects requested")
	}

	token = readFromFile(gitlabTokenFile, "token")
	if verbose {
		log.Printf("Get token from: %s", gitlabTokenFile)
	}

	if projectsFile == "" {
		projectsFile = defaultProjectsFile
	}
	if *exportProjectsOnly {
		log.Printf("Export current Gitlab projects to %s file", projectsFile)
		projectsOnGitlab, _ := getProjectsFromGitlab(gitlabUrl, token)
		exportProjects(projectsFile, projectsOnGitlab)
		log.Print("Exit now because export is done")
		os.Exit(0)
	}

	projectfile, err := os.OpenFile(projectsFile, os.O_RDONLY, 0644)
	if err == nil {
		projects := importProjects(projectsFile)
		err = projectfile.Close()
		if err != nil {
			log.Fatalln("Cannot close project file")
		}
		repoUrl := getGitUrl("origin")
		if verbose {
			log.Printf("Get git repository url: %s", repoUrl)
		}
		id := getProjectIdByRepoUrl(repoUrl, projects)
		if id > 0 {
			projectId = strconv.Itoa(id)
			if verbose {
				log.Printf("Get projectId: %s from git repository URL %s", projectId, repoUrl)
			}

		}
	} else {
		if verbose {
			log.Printf("Cannot open %s file", projectsFile)
		}
	}

	if projectId == "" {
		projectId = readFromFile(defaultIdFile, "project Id")
		if verbose {
			log.Printf("Get projectId: %s from %s file", projectId, defaultIdFile)
		}
	}
	if verbose {
		log.Printf("Using projectId: %s", projectId)
	}

	log.Printf("Fetching envs from gitlab with URL %s", gitlabUrl)
	envsOnGitlab, _ := getEnvsFromGitlab(gitlabUrl, token, projectId)
	log.Printf("Fetching vars from gitlab with URL %s", gitlabUrl)
	varsOnGitlab, _ := getVarsFromGitlab(gitlabUrl, token, projectId)
	if debug {
		file, err := os.OpenFile(defaultDebugFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Println("Debug file does not exists and cannot be created")
			os.Exit(1)
		}
		w := bufio.NewWriter(file)
		_, err = fmt.Fprint(w, "================= Envs =================\n")
		if err != nil {
			log.Fatalln("Cannot write into debug file")
		}
		_, err = fmt.Fprintf(w, "%v\n", envsOnGitlab)
		if err != nil {
			log.Fatalln("Cannot write into debug file")
		}
		_, err = fmt.Fprint(w, "================= Vars =================\n")
		if err != nil {
			log.Fatalln("Cannot write into debug file")
		}
		_, err = fmt.Fprintf(w, "%v\n", varsOnGitlab)
		if err != nil {
			log.Fatalln("Cannot write into debug file")
		}
		_, err = fmt.Fprint(w, "================= End ==================\n")
		if err != nil {
			log.Fatalln("Cannot write into debug file")
		}
		err = w.Flush()
		if err != nil {
			log.Fatalln("Cannot flush debug file")
		}
		err = file.Close()
		if err != nil {
			log.Fatalln("Cannot close debug file")
		}
		log.Printf("Debug file is written: %s", defaultDebugFile)
	}

	if *export {
		log.Printf("Export current Gitlab vars to %s file", varsFile)
		exportVars(varsFile, varsOnGitlab)
		log.Printf("Export current Gitlab envs to %s file", envsFile)
		exportEnvs(envsFile, envsOnGitlab)
		log.Print("Exit now because export is done")
		os.Exit(0)
	}
	if verbose {
		log.Print("Compare the environments between those present on GitLab and those in environment file")
	}
	envfile, err := os.OpenFile(envsFile, os.O_RDONLY, 0644)
	if err == nil {
		err = envfile.Close()
		if err != nil {
			log.Fatalln("Cannot close env file (test)")
		}

		envToExport := importEnvs(envsFile)
		envToAdd, envToDelete, envToUpdate := compareEnvSet(envsOnGitlab, envToExport)
		for _, item := range envToAdd {
			if !*dryrun {
				err = insertEnv(gitlabUrl, token, projectId, item)
				if err != nil {
					log.Fatalf("Cannot insert env %s", item.Name)
				}
			}
		}
		if len(envToAdd) == 0 {
			log.Print("No env to insert")
		}
		for _, item := range envToUpdate {
			if !*dryrun {
				err = updateEnv(gitlabUrl, token, projectId, item)
				if err != nil {
					log.Fatalf("Cannot update env %s", item.Name)
				}
			}
		}
		if len(envToUpdate) == 0 {
			log.Print("No env to update")
		}
		if len(envToDelete) == 0 {
			log.Print("No env to delete")
		}
		if *deleteIsActive && !*dryrun {
			for _, item := range envToDelete {
				err = deleteEnv(gitlabUrl, token, projectId, item)
				if err != nil {
					log.Fatalf("Cannot delete env %s", item.Name)
				}
			}
		} else {
			if len(envToDelete) > 0 {
				log.Printf("%d env(s) may be deleted, but delete flag in command line is not set", len(envToDelete))
			}
		}
	}
	varfile, err := os.OpenFile(varsFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Nothing to do because var file cannot be found. You may create it with the export flag in command line.")
		os.Exit(1)
	}
	err = varfile.Close()
	if err != nil {
		log.Fatalln("Cannot close var file (test)")
	}
	if verbose {
		log.Print("Compare the environments between those present on GitLab and those in variable file")
	}
	missingEnvs := getMissingEnvs(getEnvsFromVars(varsOnGitlab), envsOnGitlab)
	for _, env := range missingEnvs {
		var newenv GitlabEnv
		newenv.Name = env
		err = insertEnv(gitlabUrl, token, projectId, newenv)
		if err != nil {
			log.Fatalf("Cannot insert env %s", newenv.Name)
		}
	}
	if len(missingEnvs) == 0 && verbose {
		log.Print("All required envs for vars are present")
	}
	if verbose {
		log.Print("Compare the variables between those present on GitLab and those in variable file")
	}
	varToExport := importVars(varsFile)

	toAdd, toDelete, toUpdate := compareVarSet(varsOnGitlab, varToExport)
	if *dryrun {
		log.Print("Exit now because dryrun mode is active")
		return
	}
	for _, item := range toAdd {
		err = insertVar(gitlabUrl, token, projectId, item)
		if err != nil {
			log.Fatalf("Cannot insert var %s", item.Key)
		}
	}
	if len(toAdd) == 0 {
		log.Print("No var to insert")
	}
	for _, item := range toUpdate {
		err = updateVar(gitlabUrl, token, projectId, item)
		if err != nil {
			log.Fatalf("Cannot update var %s", item.Key)
		}
	}
	if len(toUpdate) == 0 {
		log.Print("No var to update")
	}
	if len(toDelete) == 0 {
		log.Print("No var to delete")
	}
	if *deleteIsActive {
		for _, item := range toDelete {
			err = deleteVar(gitlabUrl, token, projectId, item)
			if err != nil {
				log.Fatalf("Cannot delete var %s", item.Key)
			}
		}
	} else {
		if len(toDelete) > 0 {
			log.Printf("%d var(s) may be deleted, but delete flag in command line is not set", len(toDelete))
		}
	}
	log.Print("Exit")
}

func readFromFile(filename string, kind string) string {
	if verbose {
		if len(kind) > 0 {
			log.Printf("Try to read %s from %s file", kind, filename)
		} else {
			log.Printf("Try to read %s file", filename)
		}
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		if len(kind) > 0 {
			log.Fatal(err, " and no "+kind+" specified in command line")
		} else {
			log.Fatal(err)
		}
		return ""
	}
	return strings.TrimSpace(string(content))
}

func callGitlabApi(gitlabUrl string, token string, method string, data []byte) ([]byte, int, error) {
	nbPage := 0
	if verbose {
		log.Printf("Querying URL: \"%s\"", gitlabUrl)
	}
	req, err := http.NewRequest(method, gitlabUrl, bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Private-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, 0, err
	}
	if len(res.Header["X-Total-Pages"]) > 0 {
		nbPage, _ = strconv.Atoi(res.Header["X-Total-Pages"][0])
		currrentPage, _ := strconv.Atoi(res.Header["X-Page"][0])
		nbItems, _ := strconv.Atoi(res.Header["X-Total"][0])
		if verbose {
			log.Printf("Pagination - Reading page %d/%d - There is %d item(s) to download", currrentPage, nbPage, nbItems)
		}
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln("Cannot close body", err)
		}
	}()

	json, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and body: %s\n", res.StatusCode, json)
	}
	if err != nil {
		log.Fatal(err)
	}
	return json, nbPage, nil
}

func writeFile(filename string, json []byte) {
	err := os.WriteFile(filename, json, 0644)
	if err != nil {
		log.Println("Export to file", filename, err)
		return
	}
}

func getGitUrl(remote string) string {
	inidata, err := ini.Load(".git/config")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	section := inidata.Section("remote \"" + remote + "\"")
	url := section.Key("url").String()
	if verbose {
		log.Printf("On remote %s, found url %s", remote, url)
	}
	return url
}
