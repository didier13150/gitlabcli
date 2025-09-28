package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/didier13150/gitlablib"
)

type GLCliConfig struct {
	GitlabUrl      string
	IdFile         string
	GroupIdFile    string
	VarsFile       string
	GroupVarsFile  string
	GlobalVarsFile string
	EnvsFile       string
	ProjectsFile   string
	DebugFile      string
	TokenFile      string
	RemoteName     string
	DebugMode      bool
	VerboseMode    bool
	DryrunMode     bool
	ExportMode     bool
	DeleteMode     bool
	BootstrapMode  bool
}

type GLCli struct {
	Config     GLCliConfig
	ProjectId  string
	GroupId    string
	RemoteName string
	token      string
	vars       gitlablib.GitlabVar
	envs       gitlablib.GitlabEnv
	projects   gitlablib.GitlabProject
}

func NewGLCli() GLCli {

	glcli := GLCli{}
	if len(os.Getenv("GLCLI_GITLAB_URL")) > 0 {
		glcli.Config.GitlabUrl = os.Getenv("GLCLI_GITLAB_URL")
	} else {
		glcli.Config.GitlabUrl = "https://gitlab.com"
	}
	if len(os.Getenv("GLCLI_ID_FILE")) > 0 {
		glcli.Config.IdFile = os.Getenv("GLCLI_ID_FILE")
	} else {
		glcli.Config.IdFile = ".gitlab.id"
	}
	if len(os.Getenv("GLCLI_GROUP_ID_FILE")) > 0 {
		glcli.Config.GroupIdFile = os.Getenv("GLCLI_GROUP_ID_FILE")
	} else {
		glcli.Config.GroupIdFile = ".gitlab.gid"
	}
	if len(os.Getenv("GLCLI_VAR_FILE")) > 0 {
		glcli.Config.VarsFile = os.Getenv("GLCLI_VAR_FILE")
	} else {
		glcli.Config.VarsFile = ".gitlab-vars.json"
	}
	if len(os.Getenv("GLCLI_GROUP_VAR_FILE")) > 0 {
		glcli.Config.GroupVarsFile = os.Getenv("GLCLI_GROUP_VAR_FILE")
	} else {
		glcli.Config.GroupVarsFile = ".gitlab-groupvars.json"
	}
	if len(os.Getenv("GLCLI_GLOBAL_VAR_FILE")) > 0 {
		glcli.Config.GlobalVarsFile = os.Getenv("GLCLI_GLOBAL_VAR_FILE")
	} else {
		glcli.Config.GlobalVarsFile = ".gitlab-globalvars.json"
	}
	if len(os.Getenv("GLCLI_ENV_FILE")) > 0 {
		glcli.Config.EnvsFile = os.Getenv("GLCLI_ENV_FILE")
	} else {
		glcli.Config.EnvsFile = ".gitlab-envs.json"
	}
	if len(os.Getenv("GLCLI_PROJECT_FILE")) > 0 {
		glcli.Config.ProjectsFile = os.Getenv("GLCLI_PROJECT_FILE")
	} else {
		glcli.Config.ProjectsFile = os.Getenv("HOME") + "/.gitlab-projects.json"
	}
	if len(os.Getenv("GLCLI_TOKEN_FILE")) > 0 {
		glcli.Config.TokenFile = os.Getenv("GLCLI_TOKEN_FILE")
	} else {
		glcli.Config.TokenFile = os.Getenv("HOME") + "/.gitlab.token"
	}
	if len(os.Getenv("GLCLI_DEBUG_FILE")) > 0 {
		glcli.Config.DebugFile = os.Getenv("GLCLI_DEBUG_FILE")
	} else {
		glcli.Config.DebugFile = "debug.txt"
	}
	if len(os.Getenv("GLCLI_REMOTE_NAME")) > 0 {
		glcli.Config.RemoteName = os.Getenv("GLCLI_REMOTE_NAME")
	} else {
		glcli.Config.RemoteName = "origin"
	}

	glcli.Config.DebugMode = false
	glcli.Config.VerboseMode = false
	glcli.Config.DryrunMode = false
	glcli.Config.ExportMode = false
	glcli.Config.DeleteMode = false
	glcli.Config.BootstrapMode = false

	return glcli
}

func (glcli *GLCli) SetProjectParameters(allProjects bool, simpleRequest bool) {
	glcli.projects.SimpleRequest = simpleRequest
	glcli.projects.MembershipOnly = !allProjects
}

func (glcli *GLCli) AddVar() {
	var newvar gitlablib.GitlabVarData
	scanner := bufio.NewScanner(os.Stdin)

	varfile, err := os.OpenFile(glcli.Config.VarsFile, os.O_RDONLY, 0644)
	if err == nil {
		err = varfile.Close()
		if err != nil {
			log.Fatalln("Cannot close var file")
		}
		glcli.vars.ImportVars(glcli.Config.VarsFile)
	}

	fmt.Print("Variable key []: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newvar.Key = strings.TrimSpace(strings.ReplaceAll(scanner.Text(), " ", "_"))
	} else {
		log.Fatal("Key cannot be empty\n")
	}

	fmt.Print("Variable value ['']: ")
	scanner.Scan()
	newvar.Value = scanner.Text()

	fmt.Print("Variable environment ['*']: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newvar.Env = strings.TrimSpace(scanner.Text())
	} else {
		newvar.Env = "*"
	}

	fmt.Print("Variable description [null]: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newvar.Description = scanner.Text()
	}

	fmt.Print("Variable is raw [true]: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newvar.IsRaw, _ = strconv.ParseBool(strings.TrimSpace(string(scanner.Text())))
	} else {
		newvar.IsRaw = true
	}

	fmt.Print("Variable is protected [false]: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newvar.IsProtected, _ = strconv.ParseBool(strings.TrimSpace(string(scanner.Text())))
	} else {
		newvar.IsProtected = false
	}

	fmt.Print("Variable is hidden [false]: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newvar.IsHidden, _ = strconv.ParseBool(strings.TrimSpace(string(scanner.Text())))
	} else {
		newvar.IsHidden = false
	}

	fmt.Print("Variable is masked [false]: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newvar.IsMasked, _ = strconv.ParseBool(strings.TrimSpace(string(scanner.Text())))
	} else {
		newvar.IsMasked = false
	}

	glcli.vars.GitlabData = append(glcli.vars.FileData, newvar)
	glcli.vars.ExportVars(glcli.Config.VarsFile)
	log.Print("Exit now because var is added to vars file")
}

func (glcli *GLCli) AddEnv() {
	var newenv gitlablib.GitlabEnvData
	scanner := bufio.NewScanner(os.Stdin)

	envfile, err := os.OpenFile(glcli.Config.EnvsFile, os.O_RDONLY, 0644)
	if err == nil {
		err = envfile.Close()
		if err != nil {
			log.Fatalln("Cannot close env file")
		}
		glcli.envs.ImportEnvs(glcli.Config.EnvsFile)
	}

	fmt.Print("Environment name []: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newenv.Name = strings.TrimSpace(strings.ReplaceAll(scanner.Text(), " ", "_"))
	} else {
		log.Fatal("Name cannot be empty\n")
	}

	fmt.Print("Environment description [null]: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newenv.Description = scanner.Text()
	}

	fmt.Print("Environment URL [null]: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newenv.Url = scanner.Text()
	}

	fmt.Print("Environment state ['available']: ")
	scanner.Scan()
	if scanner.Text() != "" {
		newenv.State = strings.TrimSpace(scanner.Text())
	} else {
		newenv.State = "available"
	}

	glcli.envs.GitlabData = append(glcli.envs.FileData, newenv)
	glcli.envs.ExportEnvs(glcli.Config.EnvsFile)
	log.Print("Exit now because env is added to envs file")
}

func (glcli *GLCli) CopyVars(envfrom string, envto string) {
	var newvar gitlablib.GitlabVarData

	varfile, err := os.OpenFile(glcli.Config.VarsFile, os.O_RDONLY, 0644)
	if err == nil {
		err = varfile.Close()
		if err != nil {
			log.Fatalln("Cannot close var file")
		}
		glcli.vars.ImportVars(glcli.Config.VarsFile)
	}
	envfile, err := os.OpenFile(glcli.Config.EnvsFile, os.O_RDONLY, 0644)
	if err == nil {
		err = envfile.Close()
		if err != nil {
			log.Fatalln("Cannot close env file")
		}
		glcli.envs.ImportEnvs(glcli.Config.EnvsFile)
	}

	// Check if envto exists or create it (in file) before processing
	found := false
	for _, env := range glcli.envs.FileData {
		if env.Name == envto {
			found = true
		}
	}
	if !found {
		// Add env to
		var newenv gitlablib.GitlabEnvData
		newenv.Name = envto
		glcli.envs.GitlabData = append(glcli.envs.FileData, newenv)
		glcli.envs.ExportEnvs(glcli.Config.EnvsFile)
	}

	var toAdd []gitlablib.GitlabVarData
	for _, variable := range glcli.vars.FileData {
		if variable.Env == envfrom {
			log.Printf("Found %s (%s)", variable.Key, variable.Env)
			found := false
			for _, var2 := range glcli.vars.FileData {
				if var2.Env == envto && var2.Key == variable.Key {
					found = true
				}
			}
			if !found {
				newvar = variable
				newvar.Env = envto
				toAdd = append(toAdd, newvar)
			}
		}
	}

	glcli.vars.GitlabData = append(glcli.vars.FileData, toAdd...)
	glcli.vars.ExportVars(glcli.Config.VarsFile)
	log.Printf("Exit now because vars from %s env are copied to %s env", envfrom, envto)
}

func (glcli *GLCli) Bootstrap() {
	var varExample gitlablib.GitlabVarData
	var envExample gitlablib.GitlabEnvData
	varExample.Key = "VAR_KEY"
	varExample.Value = "VAR_VALUE"
	varExample.Env = "*"
	varExample.IsRaw = true
	varExample.Description = "Description of VAR_KEY"
	glcli.vars.GitlabData = append(glcli.vars.GitlabData, varExample)
	envExample.Name = "ENV_NAME"
	envExample.Description = "Description of ENV_NAME"
	envExample.State = "available"
	glcli.envs.GitlabData = append(glcli.envs.GitlabData, envExample)
	f, err := os.OpenFile(glcli.Config.VarsFile, os.O_RDONLY, 0644)
	if err == nil {
		log.Fatal("Cannot bootstrap because var file exists.")
	} else {
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatalln("Cannot close file", err)
			}
		}()
	}
	f, err = os.OpenFile(glcli.Config.EnvsFile, os.O_RDONLY, 0644)
	if err == nil {
		log.Fatal("Cannot bootstrap because env file exists.")
	} else {
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatalln("Cannot close file", err)
			}
		}()
	}

	glcli.vars.ExportVars(glcli.Config.VarsFile)
	glcli.envs.ExportEnvs(glcli.Config.EnvsFile)
	log.Print("Exit now because bootstrap is done")
}

func (glcli *GLCli) Setup() {
	glcli.token = gitlablib.ReadFromFile(glcli.Config.TokenFile, "token", glcli.Config.VerboseMode)
	glcli.vars = gitlablib.NewGitlabVar(glcli.Config.GitlabUrl, glcli.token, glcli.Config.VerboseMode)
	glcli.envs = gitlablib.NewGitlabEnv(glcli.Config.GitlabUrl, glcli.token, glcli.Config.VerboseMode)
	glcli.projects = gitlablib.NewGitlabProject(glcli.Config.GitlabUrl, glcli.token, glcli.Config.VerboseMode)

	if glcli.Config.VerboseMode {
		log.Printf("Get token from: %s", glcli.Config.TokenFile)
	}
	if glcli.Config.DryrunMode {
		glcli.projects.DryrunMode = glcli.Config.DryrunMode
		glcli.vars.DryrunMode = glcli.Config.DryrunMode
		glcli.envs.DryrunMode = glcli.Config.DryrunMode
	}
}

func (glcli *GLCli) ExportProjects() {
	log.Printf("Export current Gitlab projects to %s file", glcli.Config.ProjectsFile)
	err := glcli.projects.GetProjectsFromGitlab()
	if err != nil {
		log.Fatal("Cannot fetch projects from gitlab")
	}
	glcli.projects.ExportProjects(glcli.Config.ProjectsFile)
	log.Print("Exit now because project export is done")
}

func (glcli *GLCli) AdminRun() {
	log.Printf("Fetching global vars from gitlab with URL %s", glcli.Config.GitlabUrl)
	err := glcli.vars.GetGlobalVarsFromGitlab()
	if err != nil {
		log.Fatal("Cannot fetch global vars from gitlab project")
	}
	// log.Printf("Fetching envs from gitlab with URL %s", glcli.Config.GitlabUrl)
	// err = glcli.envs.GetEnvsFromGitlab()
	// if err != nil {
	// 	log.Fatal("Cannot fetch envs from gitlab")
	// }

	if glcli.Config.DebugMode {
		glcli.debug()
	}
	if glcli.Config.ExportMode {
		log.Printf("Export current Gitlab global vars to %s file", glcli.Config.GlobalVarsFile)
		glcli.vars.ExportGlobalVars(glcli.Config.GlobalVarsFile)
		// log.Printf("Export current Gitlab envs to %s file", glcli.Config.EnvsFile)
		// glcli.envs.ExportEnvs(glcli.Config.EnvsFile)
		log.Print("Exit now because export is done")
		return
	}
	globalvarfile, err := os.OpenFile(glcli.Config.GlobalVarsFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Nothing to do because global var file cannot be found.")
		os.Exit(1)
	}
	err = globalvarfile.Close()
	if err != nil {
		log.Fatalln("Cannot close global var file (test)")
	}

	glcli.vars.ImportGlobalVars(glcli.Config.GlobalVarsFile)

	toAdd, toDelete, toUpdate := glcli.vars.CompareGlobalVar()
	if glcli.Config.VerboseMode {
		log.Print("Compare the group variables between those present on GitLab and those in variable file")
	}
	if glcli.Config.DryrunMode {
		log.Print("Exit now because dryrun mode is active")
		return
	}
	for _, item := range toAdd {
		err = glcli.vars.InsertGlobalVar(item)
		if err != nil {
			log.Fatalf("Cannot insert global var %s", item.Key)
		}
	}
	if len(toAdd) == 0 {
		log.Print("No global var to insert")
	}
	for _, item := range toUpdate {
		err = glcli.vars.UpdateGlobalVar(item)
		if err != nil {
			log.Fatalf("Cannot update global var %s", item.Key)
		}
	}
	if len(toUpdate) == 0 {
		log.Print("No global var to update")
	}
	if len(toDelete) == 0 {
		log.Print("No global var to delete")
	}
	if glcli.Config.DeleteMode {
		for _, item := range toDelete {
			err = glcli.vars.DeleteGlobalVar(item)
			if err != nil {
				log.Fatalf("Cannot delete global var %s", item.Key)
			}
		}
	} else {
		if len(toDelete) > 0 {
			log.Printf("%d global var(s) may be deleted, but delete flag in command line is not set", len(toDelete))
		}
	}
}

func (glcli *GLCli) Run() {

	projectfile, err := os.OpenFile(glcli.Config.ProjectsFile, os.O_RDONLY, 0644)
	if err == nil {
		glcli.projects.ImportProjects(glcli.Config.ProjectsFile)
		err = projectfile.Close()
		if err != nil {
			log.Fatalln("Cannot close project file")
		}
		repoUrl := gitlablib.GetGitUrl(glcli.RemoteName, glcli.Config.VerboseMode)
		if glcli.Config.VerboseMode {
			log.Printf("Get git repository url for remote %s: %s", glcli.RemoteName, repoUrl)
		}
		if glcli.ProjectId == "" {
			id := glcli.projects.GetProjectIdByRepoUrl(repoUrl)
			if id > 0 {
				glcli.ProjectId = strconv.Itoa(id)
				if glcli.Config.VerboseMode {
					log.Printf("Get projectId: %s from git repository URL %s", glcli.ProjectId, repoUrl)
				}
				gid := glcli.projects.GetGroupIdByProjectId(id)
				if gid > 0 {
					glcli.GroupId = strconv.Itoa(gid)
					if glcli.Config.VerboseMode {
						log.Printf("Get groupId: %s from id %s", glcli.GroupId, glcli.ProjectId)
					}
				}
			}
		}
	} else {
		if glcli.Config.VerboseMode {
			log.Printf("Cannot open %s file", glcli.Config.ProjectsFile)
		}
	}
	if glcli.ProjectId == "" {
		glcli.ProjectId = gitlablib.ReadFromFile(glcli.Config.IdFile, "project Id", glcli.Config.VerboseMode)
		if glcli.Config.VerboseMode {
			log.Printf("Get projectId: %s from %s file", glcli.ProjectId, glcli.Config.IdFile)
		}
	}
	if glcli.GroupId == "" {
		glcli.GroupId = gitlablib.ReadFromFile(glcli.Config.GroupIdFile, "group Id", glcli.Config.VerboseMode)
		if glcli.Config.VerboseMode {
			log.Printf("Get GroupId: %s from %s file", glcli.GroupId, glcli.Config.GroupIdFile)
		}
	}
	if glcli.Config.VerboseMode {
		log.Printf("Using projectId: %s, groupId: %s", glcli.ProjectId, glcli.GroupId)
	}
	glcli.envs.ProjectId = glcli.ProjectId
	glcli.vars.ProjectId = glcli.ProjectId
	glcli.vars.GroupId = glcli.GroupId

	log.Printf("Fetching envs from gitlab with URL %s", glcli.Config.GitlabUrl)
	err = glcli.envs.GetEnvsFromGitlab()
	if err != nil {
		log.Fatal("Cannot fetch envs from gitlab")
	}
	log.Printf("Fetching vars from gitlab with URL %s", glcli.Config.GitlabUrl)
	err = glcli.vars.GetVarsFromGitlab()
	if err != nil {
		log.Fatal("Cannot fetch vars from gitlab project")
	}
	if glcli.GroupId != "" {
		log.Printf("Fetching group vars from gitlab with URL %s", glcli.Config.GitlabUrl)
		err = glcli.vars.GetGroupVarsFromGitlab()
		if err != nil {
			log.Print("Cannot fetch vars from gitlab group")
		}
	}
	if glcli.Config.DebugMode {
		glcli.debug()
	}

	if glcli.Config.ExportMode {
		log.Printf("Export current Gitlab vars to %s file", glcli.Config.VarsFile)
		glcli.vars.ExportVars(glcli.Config.VarsFile)
		log.Printf("Export current Gitlab group vars to %s file", glcli.Config.GroupVarsFile)
		glcli.vars.ExportGroupVars(glcli.Config.GroupVarsFile)
		log.Printf("Export current Gitlab envs to %s file", glcli.Config.EnvsFile)
		glcli.envs.ExportEnvs(glcli.Config.EnvsFile)
		log.Print("Exit now because export is done")
		return
	}

	if glcli.Config.VerboseMode {
		log.Print("Compare the environments between those present on GitLab and those in environment file")
	}

	envfile, err := os.OpenFile(glcli.Config.EnvsFile, os.O_RDONLY, 0644)
	if err == nil {
		err = envfile.Close()
		if err != nil {
			log.Fatalln("Cannot close env file (test)")
		}

		glcli.envs.ImportEnvs(glcli.Config.EnvsFile)
		envToAdd, envToDelete, envToUpdate := glcli.envs.CompareEnv()
		for _, item := range envToAdd {
			if !glcli.Config.DryrunMode {
				err = glcli.envs.InsertEnv(item)
				if err != nil {
					log.Fatalf("Cannot insert env %s", item.Name)
				}
			}
		}
		if len(envToAdd) == 0 {
			log.Print("No env to insert")
		}
		for _, item := range envToUpdate {
			if !glcli.Config.DryrunMode {
				err = glcli.envs.UpdateEnv(item)
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
		if glcli.Config.DeleteMode && !glcli.Config.DryrunMode {
			for _, item := range envToDelete {
				err = glcli.envs.DeleteEnv(item)
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
	varfile, err := os.OpenFile(glcli.Config.VarsFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Nothing to do because var file cannot be found. You may create it with the export flag in command line.")
		os.Exit(1)
	}
	err = varfile.Close()
	if err != nil {
		log.Fatalln("Cannot close var file (test)")
	}
	if glcli.Config.VerboseMode {
		log.Print("Compare the environments between those present on GitLab and those in variable files")
	}

	glcli.vars.ImportVars(glcli.Config.VarsFile)
	glcli.vars.ImportGroupVars(glcli.Config.GroupVarsFile)

	missingEnvs := glcli.envs.GetMissingEnvs(glcli.vars.GetEnvsFromVars())
	for _, env := range missingEnvs {
		if !glcli.Config.DryrunMode {
			var newenv gitlablib.GitlabEnvData
			newenv.Name = env
			err = glcli.envs.InsertEnv(newenv)
			if err != nil {
				log.Fatalf("Cannot insert env %s", newenv.Name)
			}
		}
	}
	if len(missingEnvs) == 0 && glcli.Config.VerboseMode {
		log.Print("All required envs for vars are present")
	}

	if glcli.Config.VerboseMode {
		log.Print("Compare the variables between those present on GitLab and those in variable file")
	}
	toAdd, toDelete, toUpdate := glcli.vars.CompareVar()
	if glcli.Config.VerboseMode {
		log.Print("Compare the group variables between those present on GitLab and those in variable file")
	}
	toGroupAdd, toGroupDelete, toGroupUpdate := glcli.vars.CompareGroupVar()
	if glcli.Config.DryrunMode {
		log.Print("Exit now because dryrun mode is active")
		return
	}
	for _, item := range toAdd {
		err = glcli.vars.InsertVar(item)
		if err != nil {
			log.Fatalf("Cannot insert var %s", item.Key)
		}
	}
	if len(toAdd) == 0 {
		log.Print("No var to insert")
	}
	for _, item := range toUpdate {
		err = glcli.vars.UpdateVar(item)
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
	if glcli.Config.DeleteMode {
		for _, item := range toDelete {
			err = glcli.vars.DeleteVar(item)
			if err != nil {
				log.Fatalf("Cannot delete var %s", item.Key)
			}
		}
	} else {
		if len(toDelete) > 0 {
			log.Printf("%d var(s) may be deleted, but delete flag in command line is not set", len(toDelete))
		}
	}
	for _, item := range toGroupAdd {
		err = glcli.vars.InsertGroupVar(item)
		if err != nil {
			log.Fatalf("Cannot insert group var %s", item.Key)
		}
	}
	if len(toGroupAdd) == 0 {
		log.Print("No group var to insert")
	}
	for _, item := range toGroupUpdate {
		err = glcli.vars.UpdateGroupVar(item)
		if err != nil {
			log.Fatalf("Cannot update group var %s", item.Key)
		}
	}
	if len(toGroupUpdate) == 0 {
		log.Print("No group var to update")
	}
	if len(toGroupDelete) == 0 {
		log.Print("No group var to delete")
	}
	if glcli.Config.DeleteMode {
		for _, item := range toGroupDelete {
			err = glcli.vars.DeleteGroupVar(item)
			if err != nil {
				log.Fatalf("Cannot delete group var %s", item.Key)
			}
		}
	} else {
		if len(toGroupDelete) > 0 {
			log.Printf("%d group var(s) may be deleted, but delete flag in command line is not set", len(toGroupDelete))
		}
	}
	log.Print("Exit")
}

func (glcli GLCli) debug() {
	file, err := os.OpenFile(glcli.Config.DebugFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Debug file does not exists and cannot be created")
		os.Exit(1)
	}
	w := bufio.NewWriter(file)
	_, err = fmt.Fprint(w, "================= Envs =================\n")
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprintf(w, "%v\n", glcli.envs.GitlabData)
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprint(w, "================= Vars =================\n")
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprintf(w, "%v\n", glcli.vars.GitlabData)
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprint(w, "============== Group Vars ==============\n")
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprintf(w, "%v\n", glcli.vars.GitlabGroupData)
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprint(w, "============== Global Vars =============\n")
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprintf(w, "%v\n", glcli.vars.GitlabGlobalData)
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprint(w, "=============== Projects ===============\n")
	if err != nil {
		log.Fatalln("Cannot write into debug file")
	}
	_, err = fmt.Fprintf(w, "%v\n", glcli.projects.Data)
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
	log.Printf("Debug file is written: %s", glcli.Config.DebugFile)
}
