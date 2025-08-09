package gitlabcli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type GLCliConfig struct {
	GitlabUrl         string
	IdFile            string
	VarsFile          string
	EnvsFile          string
	ProjectsFile      string
	DebugFile         string
	TokenFile         string
	RemoteName        string
	DebugMode         bool
	VerboseMode       bool
	DryrunMode        bool
	ExportMode        bool
	ExportProjectMode bool
	DeleteMode        bool
}

type GLCli struct {
	Config     GLCliConfig
	ProjectId  string
	RemoteName string
	token      string
	vars       GitlabVar
	envs       GitlabEnv
	projects   GitlabProject
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
	if len(os.Getenv("GLCLI_VAR_FILE")) > 0 {
		glcli.Config.VarsFile = os.Getenv("GLCLI_VAR_FILE")
	} else {
		glcli.Config.VarsFile = ".gitlab-vars.json"
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
	glcli.Config.ExportProjectMode = false
	glcli.Config.DeleteMode = false

	return glcli
}

func (glcli *GLCli) Run() {

	glcli.token = readFromFile(glcli.Config.TokenFile, "token", glcli.Config.VerboseMode)
	glcli.vars = NewGitlabVar(glcli.Config.GitlabUrl, glcli.token, glcli.Config.VerboseMode)
	glcli.envs = NewGitlabEnv(glcli.Config.GitlabUrl, glcli.token, glcli.Config.VerboseMode)
	glcli.projects = NewGitlabProject(glcli.Config.GitlabUrl, glcli.token, glcli.Config.VerboseMode)

	if glcli.Config.VerboseMode {
		log.Printf("Get token from: %s", glcli.Config.TokenFile)
	}

	if glcli.Config.ExportProjectMode {
		log.Printf("Export current Gitlab projects to %s file", glcli.Config.ProjectsFile)
		glcli.projects.GetProjectsFromGitlab()
		glcli.projects.ExportProjects(glcli.Config.ProjectsFile)
		log.Print("Exit now because project export is done")
		os.Exit(0)
	}

	projectfile, err := os.OpenFile(glcli.Config.ProjectsFile, os.O_RDONLY, 0644)
	if err == nil {
		glcli.projects.ImportProjects(glcli.Config.ProjectsFile)
		err = projectfile.Close()
		if err != nil {
			log.Fatalln("Cannot close project file")
		}
		repoUrl := getGitUrl(glcli.RemoteName, glcli.Config.VerboseMode)
		if glcli.Config.VerboseMode {
			log.Printf("Get git repository url for remote %s: %s", glcli.RemoteName, repoUrl)
		}
		id := glcli.projects.GetProjectIdByRepoUrl(repoUrl)
		if id > 0 {
			glcli.ProjectId = strconv.Itoa(id)
			if glcli.Config.VerboseMode {
				log.Printf("Get projectId: %s from git repository URL %s", glcli.ProjectId, repoUrl)
			}
		}
	} else {
		if glcli.Config.VerboseMode {
			log.Printf("Cannot open %s file", glcli.Config.ProjectsFile)
		}
	}
	if glcli.ProjectId == "" {
		glcli.ProjectId = readFromFile(glcli.Config.IdFile, "project Id", glcli.Config.VerboseMode)
		if glcli.Config.VerboseMode {
			log.Printf("Get projectId: %s from %s file", glcli.ProjectId, glcli.Config.IdFile)
		}
	}
	if glcli.Config.VerboseMode {
		log.Printf("Using projectId: %s", glcli.ProjectId)
	}
	glcli.envs.ProjectId = glcli.ProjectId
	glcli.vars.ProjectId = glcli.ProjectId

	log.Printf("Fetching envs from gitlab with URL %s", glcli.Config.GitlabUrl)
	glcli.envs.GetEnvsFromGitlab()
	log.Printf("Fetching vars from gitlab with URL %s", glcli.Config.GitlabUrl)
	glcli.vars.GetVarsFromGitlab()
	if glcli.Config.DebugMode {
		glcli.debug()
	}

	if glcli.Config.ExportMode {
		log.Printf("Export current Gitlab vars to %s file", glcli.Config.VarsFile)
		glcli.vars.ExportVars(glcli.Config.VarsFile)
		log.Printf("Export current Gitlab envs to %s file", glcli.Config.EnvsFile)
		glcli.envs.ExportEnvs(glcli.Config.EnvsFile)
		log.Print("Exit now because export is done")
		os.Exit(0)
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
		log.Print("Compare the environments between those present on GitLab and those in variable file")
	}
	missingEnvs := glcli.envs.GetMissingEnvs(glcli.vars.GetEnvsFromVars())
	for _, env := range missingEnvs {
		var newenv GitlabEnvData
		newenv.Name = env
		err = glcli.envs.InsertEnv(newenv)
		if err != nil {
			log.Fatalf("Cannot insert env %s", newenv.Name)
		}
	}
	if len(missingEnvs) == 0 && glcli.Config.VerboseMode {
		log.Print("All required envs for vars are present")
	}
	if glcli.Config.VerboseMode {
		log.Print("Compare the variables between those present on GitLab and those in variable file")
	}
	glcli.vars.ImportVars(glcli.Config.VarsFile)

	toAdd, toDelete, toUpdate := glcli.vars.CompareVar()
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
