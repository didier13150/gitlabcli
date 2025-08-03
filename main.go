package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"tartarefr.eu/gitlabcli"
)

func main() {

	glcli := gitlabcli.NewGLCli()

	var projectId = flag.String("id", "", "Gitlab project identifiant.")
	var projectIdFile = flag.String("idfile", glcli.Config.IdFile, "Gitlab project identifiant file.")
	var varsFile = flag.String("varfile", glcli.Config.VarsFile, "File which contains vars.")
	var envsFile = flag.String("envfile", glcli.Config.EnvsFile, "File which contains envs.")
	var projectsFile = flag.String("projectfile", glcli.Config.ProjectsFile, "File which contains projects.")
	var gitlabUrl = flag.String("url", glcli.Config.GitlabUrl, "Gitlab URL.")
	var gitlabTokenFile = flag.String("tokenfile", glcli.Config.TokenFile, "File which contains token to access Gitlab API.")
	var remoteName = flag.String("remote", glcli.Config.RemoteName, "Git remote name.")
	var verbose = flag.Bool("verbose", glcli.Config.VerboseMode, "Make application more talkative.")
	var debug = flag.Bool("debug", glcli.Config.DebugMode, "Enable debug mode")
	var dryrun = flag.Bool("dryrun", glcli.Config.DryrunMode, "Run in dry-run mode (read only).")
	var export = flag.Bool("export", glcli.Config.ExportMode, "Export current variables in var file.")
	var exportProjectsOnly = flag.Bool("export-projects", glcli.Config.ExportProjectMode, "Export current projects in project file.")
	var deleteIsActive = flag.Bool("delete", glcli.Config.DeleteMode, "Delete Gitlab var if not present in var file.")

	flag.Usage = func() {
		fmt.Printf("Usage: " + os.Args[0] + " [options]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

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
		glcli.Config.ExportProjectMode = true
	}
	if *deleteIsActive {
		log.Print("Delete mode is active")
		glcli.Config.DeleteMode = true
	}
	if projectIdFile != nil {
		glcli.Config.IdFile = *projectIdFile
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
	if gitlabUrl != nil {
		glcli.Config.GitlabUrl = *gitlabUrl
	}
	if gitlabTokenFile != nil {
		glcli.Config.TokenFile = *gitlabTokenFile
	}
	if projectId != nil {
		glcli.ProjectId = *projectId
	}
	if remoteName != nil {
		glcli.RemoteName = *remoteName
	}
	glcli.Run()
}
