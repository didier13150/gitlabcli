package main

import (
	"os"
	"testing"

	"github.com/didier13150/gitlablib"
)

func TestGLCliConfig(t *testing.T) {
	url := "https://gitlab.com"
	idfile := "project.id"
	gidfile := "project.gid"
	varsfile := "project.vars"
	groupvarsfile := "project.groupvars"
	envsfile := "project.envs"
	projectsfile := "projects.list"
	debugfile := "debug.log"
	tokenfile := "project.token"
	remoteName := "origin"
	debug := true
	verbose := true
	dryrun := true
	export := true
	delete := true

	config := GLCliConfig{}
	config.GitlabUrl = url
	config.IdFile = idfile
	config.GroupIdFile = gidfile
	config.GroupVarsFile = groupvarsfile
	config.VarsFile = varsfile
	config.EnvsFile = envsfile
	config.ProjectsFile = projectsfile
	config.DebugFile = debugfile
	config.TokenFile = tokenfile
	config.RemoteName = remoteName
	config.DebugMode = debug
	config.VerboseMode = verbose
	config.DryrunMode = dryrun
	config.ExportMode = export
	config.DeleteMode = delete

	if config.GitlabUrl != url {
		t.Errorf(`TestGLCliConfig(GitlabUrl) = %s, want %s`, config.GitlabUrl, url)
	}
	if config.IdFile != idfile {
		t.Errorf(`TestGLCliConfig(IdFile) = %s, want %s`, config.IdFile, idfile)
	}
	if config.GroupIdFile != gidfile {
		t.Errorf(`TestGLCliConfig(IdFile) = %s, want %s`, config.GroupIdFile, gidfile)
	}
	if config.VarsFile != varsfile {
		t.Errorf(`TestGLCliConfig(VarsFile) = %s, want %s`, config.VarsFile, varsfile)
	}
	if config.GroupVarsFile != groupvarsfile {
		t.Errorf(`TestGLCliConfig(VarsFile) = %s, want %s`, config.GroupVarsFile, groupvarsfile)
	}
	if config.EnvsFile != envsfile {
		t.Errorf(`TestGLCliConfig(EnvsFile) = %s, want %s`, config.EnvsFile, envsfile)
	}
	if config.ProjectsFile != projectsfile {
		t.Errorf(`TestGLCliConfig(ProjectsFile) = %s, want %s`, config.ProjectsFile, projectsfile)
	}
	if config.DebugFile != debugfile {
		t.Errorf(`TestGLCliConfig(DebugFile) = %s, want %s`, config.DebugFile, debugfile)
	}
	if config.TokenFile != tokenfile {
		t.Errorf(`TestGLCliConfig(TokenFile) = %s, want %s`, config.TokenFile, tokenfile)
	}
	if config.RemoteName != remoteName {
		t.Errorf(`TestGLCliConfig(RemoteName) = %s, want %s`, config.RemoteName, remoteName)
	}
	if config.DebugMode != debug {
		t.Errorf(`TestGLCliConfig(DebugMode) = %t, want %t`, config.DebugMode, debug)
	}
	if config.VerboseMode != verbose {
		t.Errorf(`TestGLCliConfig(VerboseMode) = %t, want %t`, config.VerboseMode, verbose)
	}
	if config.DryrunMode != dryrun {
		t.Errorf(`TestGLCliConfig(DryrunMode) = %t, want %t`, config.DryrunMode, dryrun)
	}
	if config.ExportMode != export {
		t.Errorf(`TestGLCliConfig(ExportMode) = %t, want %t`, config.ExportMode, export)
	}
	if config.DeleteMode != delete {
		t.Errorf(`TestGLCliConfig(DeleteMode) = %t, want %t`, config.DeleteMode, delete)
	}
}

func TestGLCliExportProjects(t *testing.T) {
	glcli := GLCli{}
	glcli.Config.VerboseMode = false
	glcli.Config.GitlabUrl = "http://localhost:8080"
	glcli.Config.TokenFile = "/tmp/glcli.token"
	glcli.Config.ProjectsFile = "/tmp/glcli-projects.json"

	err := os.WriteFile(glcli.Config.TokenFile, []byte("token"), 0644)
	if err != nil {
		t.Errorf(`TestGLCliExportProjects(write token file) = %s`, err)
	}
	glcli.ExportProjects()

	// Check that export file have 48 projects(6 groups * 8 names in glsimulator)
	projects := gitlablib.NewGitlabProject(glcli.Config.GitlabUrl, "token", false)
	projects.ImportProjects(glcli.Config.ProjectsFile)
	if len(projects.Data) != 48 {
		t.Errorf(`TestGLCliExportProjects(count projects on export file) = %d, must be %d`, len(projects.Data), 48)
	}

	err = os.Remove(glcli.Config.TokenFile)
	if err != nil {
		t.Errorf(`TestGLCliExportProjects(delete token file) = %s`, err)
	}
	err = os.Remove(glcli.Config.ProjectsFile)
	if err != nil {
		t.Errorf(`TestGLCliExportProjects(delete project export file) = %s`, err)
	}
}

func TestGLCliExport(t *testing.T) {
	glcli := GLCli{}
	glcli.Config.VerboseMode = false
	glcli.Config.GitlabUrl = "http://localhost:8080"
	glcli.Config.TokenFile = "/tmp/glcli.token"
	glcli.Config.VarsFile = "/tmp/glcli-vars.json"
	glcli.Config.EnvsFile = "/tmp/glcli-envs.json"
	glcli.Config.GroupVarsFile = "/tmp/glcli-groupvars.json"
	glcli.Config.ExportMode = true
	glcli.ProjectId = "3"
	glcli.GroupId = "2"

	err := os.WriteFile(glcli.Config.TokenFile, []byte("token"), 0644)
	if err != nil {
		t.Errorf(`TestGLCliExport(write token file) = %s`, err)
	}
	glcli.Run()

	// Check that env export file have 2 envs
	envs := gitlablib.NewGitlabEnv(glcli.Config.GitlabUrl, "token", false)
	envs.ImportEnvs(glcli.Config.EnvsFile)
	if len(envs.FileData) != 2 {
		t.Errorf(`TestGLCliExportProjects(count vars on export file) = %d, must be %d`, len(envs.FileData), 2)
	}
	// Check that env export file have 3 vars
	vars := gitlablib.NewGitlabVar(glcli.Config.GitlabUrl, "token", false)
	vars.ImportVars(glcli.Config.VarsFile)
	if len(vars.FileData) != 3 {
		t.Errorf(`TestGLCliExportProjects(count vars on export file) = %d, must be %d`, len(vars.FileData), 3)
	}
	// Check that env export file have 3 vars
	groupvars := gitlablib.NewGitlabVar(glcli.Config.GitlabUrl, "token", false)
	vars.ImportGroupVars(glcli.Config.GroupVarsFile)
	if len(vars.FileGroupData) != 3 {
		t.Errorf(`TestGLCliExportProjects(count groupvars on export file) = %d, must be %d`, len(groupvars.FileGroupData), 3)
	}

	err = os.Remove(glcli.Config.TokenFile)
	if err != nil {
		t.Errorf(`TestGLCliExport(delete token file) = %s`, err)
	}
	err = os.Remove(glcli.Config.VarsFile)
	if err != nil {
		t.Errorf(`TestGLCliExport(delete var export file) = %s`, err)
	}
	err = os.Remove(glcli.Config.EnvsFile)
	if err != nil {
		t.Errorf(`TestGLCliExport(delete env export file) = %s`, err)
	}
}

func TestGLCliRun(t *testing.T) {
	glcli := GLCli{}
	glcli2 := GLCli{}

	glcli.Config.VerboseMode = false
	glcli.Config.GitlabUrl = "http://localhost:8080"
	glcli.Config.TokenFile = "/tmp/glcli.token"
	glcli.Config.VarsFile = "/tmp/glcli-vars.json"
	glcli.Config.EnvsFile = "/tmp/glcli-envs.json"
	glcli.Config.GroupVarsFile = "/tmp/glcli-groupvars.json"
	glcli.Config.ExportMode = true
	glcli.ProjectId = "3"
	glcli.GroupId = "2"

	err := os.WriteFile(glcli.Config.TokenFile, []byte("token"), 0644)
	if err != nil {
		t.Errorf(`TestGLCliExport(write token file) = %s`, err)
	}
	glcli.Run()
	glcli2.Config = glcli.Config
	glcli2.ProjectId = "3"
	glcli2.GroupId = "2"
	glcli2.Config.ExportMode = false
	glcli2.Run()

	err = os.Remove(glcli.Config.TokenFile)
	if err != nil {
		t.Errorf(`TestGLCliExport(delete token file) = %s`, err)
	}
	err = os.Remove(glcli.Config.VarsFile)
	if err != nil {
		t.Errorf(`TestGLCliExport(delete var export file) = %s`, err)
	}
	err = os.Remove(glcli.Config.EnvsFile)
	if err != nil {
		t.Errorf(`TestGLCliExport(delete env export file) = %s`, err)
	}
}
