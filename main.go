package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	glcli := NewGLCli()

	var projectId = flag.String("id", "", "Gitlab project identifiant.")
	var groupId = flag.String("gid", "", "Gitlab group identifiant.")
	var projectIdFile = flag.String("idfile", glcli.Config.IdFile, "Gitlab project identifiant file.")
	var groupIdFile = flag.String("gidfile", glcli.Config.GroupIdFile, "Gitlab group identifiant file.")
	var varsFile = flag.String("varfile", glcli.Config.VarsFile, "File which contains vars.")
	var envsFile = flag.String("envfile", glcli.Config.EnvsFile, "File which contains envs.")
	var groupvarsFile = flag.String("groupvarfile", glcli.Config.GroupVarsFile, "File which contains group vars.")
	var projectsFile = flag.String("projectfile", glcli.Config.ProjectsFile, "File which contains projects.")
	var gitlabUrl = flag.String("url", glcli.Config.GitlabUrl, "Gitlab URL.")
	var gitlabTokenFile = flag.String("tokenfile", glcli.Config.TokenFile, "File which contains token to access Gitlab API.")
	var remoteName = flag.String("remote", glcli.Config.RemoteName, "Git remote name.")
	var verbose = flag.Bool("verbose", glcli.Config.VerboseMode, "Make application more talkative.")
	var debug = flag.Bool("debug", glcli.Config.DebugMode, "Enable debug mode")
	var dryrun = flag.Bool("dryrun", glcli.Config.DryrunMode, "Run in dry-run mode (read only).")
	var export = flag.Bool("export", glcli.Config.ExportMode, "Export current variables in var file.")
	var exportProjectsOnly = flag.Bool("export-projects", false, "Export current projects in project file.")
	var deleteIsActive = flag.Bool("delete", glcli.Config.DeleteMode, "Delete Gitlab var if not present in var file.")
	var allProjects = flag.Bool("all-projects", false, "Export all projects, not only projects where I'm a membership.")
	var simpleRequest = flag.Bool("full-projects-data", false, "Requesting full data about projects.")
	var bootstrapIsActive = flag.Bool("bootstrap", glcli.Config.BootstrapMode, "Bootstrap varsfile and envsfile with templates")
	var addVarIsActive = flag.Bool("add-var", false, "Add variable to varsfile in interactive mode")
	var addEnvIsActive = flag.Bool("add-env", false, "Add environment to envsfile in interactive mode")
	var duplicateVarsInEnvFrom = flag.String("duplicate-from", "", "Duplicate all vars from specified env (Must be set with duplicate-to option).")
	var duplicateVarsInEnvTo = flag.String("duplicate-to", "", "Duplicate all vars from env to specified env (Must be set with duplicate-from option).")
	var adminIsActive = flag.Bool("admin", false, "Admin mode")

	flag.Usage = func() {
		fmt.Print("Export variables from json file to project gitlab variables or vice versa\n\n")
		fmt.Printf("Usage: " + os.Args[0] + " [options]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var envFrom string
	var envTo string

	if *verbose {
		log.Print("Verbose mode is active")
		glcli.Config.VerboseMode = true
	}
	if *debug {
		log.Print("Verbose mode is active")
		glcli.Config.DebugMode = true
	}
	if *dryrun {
		log.Print("Dry run mode is active")
		glcli.Config.DryrunMode = true
	}
	if *export {
		log.Print("Export requested")
		glcli.Config.ExportMode = true
	}
	if *exportProjectsOnly {
		log.Print("Export projects requested")
	}
	if *deleteIsActive {
		log.Print("Delete mode is active")
		glcli.Config.DeleteMode = true
	}
	if *bootstrapIsActive {
		log.Print("Bootstrap mode is active")
		glcli.Bootstrap()
		return
	}
	if *addVarIsActive {
		log.Print("Add variable mode is active")
		glcli.AddVar()
		return
	}
	if *addEnvIsActive {
		log.Print("Add environment mode is active")
		glcli.AddEnv()
		return
	}
	if duplicateVarsInEnvFrom != nil {
		envFrom = *duplicateVarsInEnvFrom
	}
	if duplicateVarsInEnvTo != nil {
		envTo = *duplicateVarsInEnvTo
	}
	if projectIdFile != nil {
		glcli.Config.IdFile = *projectIdFile
	}
	if groupIdFile != nil {
		glcli.Config.GroupIdFile = *groupIdFile
	}
	if varsFile != nil {
		glcli.Config.VarsFile = *varsFile
	}
	if envsFile != nil {
		glcli.Config.EnvsFile = *envsFile
	}
	if projectsFile != nil {
		glcli.Config.ProjectsFile = *projectsFile
	}
	if groupvarsFile != nil {
		glcli.Config.GroupVarsFile = *groupvarsFile
	}
	if gitlabUrl != nil {
		glcli.Config.GitlabUrl = *gitlabUrl
	}
	if gitlabTokenFile != nil {
		glcli.Config.TokenFile = *gitlabTokenFile
	}
	if projectId != nil {
		glcli.ProjectId = *projectId
	}
	if groupId != nil {
		glcli.GroupId = *groupId
	}
	if remoteName != nil {
		glcli.RemoteName = *remoteName
	}
	if *adminIsActive {
		log.Print("Admin mode is active")
	}
	if envFrom != "" && envTo != "" {
		log.Printf("Copy all variables from %s environment to %s one\n", envFrom, envTo)
		glcli.CopyVars(envFrom, envTo)
		return
	}
	glcli.SetProjectParameters(*allProjects, *simpleRequest)
	glcli.Setup()
	if *exportProjectsOnly {
		glcli.ExportProjects()
	} else if *adminIsActive {
		glcli.AdminRun()
	} else {
		glcli.Run()
	}
}
