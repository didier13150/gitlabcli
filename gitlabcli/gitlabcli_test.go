package gitlabcli

import (
	"testing"
)

func TestGLCliConfig(t *testing.T) {
	url := "https://gitlab.com"
	idfile := "project.id"
	varsfile := "project.vars"
	envsfile := "project.envs"
	projectsfile := "projects.list"
	debugfile := "debug.log"
	tokenfile := "project.token"
	remoteName := "origin"
	debug := true
	verbose := true
	dryrun := true
	export := true
	exportProject := true
	delete := true

	config := GLCliConfig{}
	config.GitlabUrl = url
	config.IdFile = idfile
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
	config.ExportProjectMode = exportProject
	config.DeleteMode = delete

	if config.GitlabUrl != url {
		t.Errorf(`TestGLCliConfig(GitlabUrl) = %s, want %s`, config.GitlabUrl, url)
	}
	if config.IdFile != idfile {
		t.Errorf(`TestGLCliConfig(IdFile) = %s, want %s`, config.IdFile, idfile)
	}
	if config.VarsFile != varsfile {
		t.Errorf(`TestGLCliConfig(VarsFile) = %s, want %s`, config.VarsFile, varsfile)
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
	if config.ExportProjectMode != exportProject {
		t.Errorf(`TestGLCliConfig(ExportProjectMode) = %t, want %t`, config.ExportProjectMode, exportProject)
	}
	if config.DeleteMode != delete {
		t.Errorf(`TestGLCliConfig(DeleteMode) = %t, want %t`, config.DeleteMode, delete)
	}
}
