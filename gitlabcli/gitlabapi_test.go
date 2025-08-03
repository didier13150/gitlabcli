package gitlabcli

import (
	"net/http"
	"testing"
)

func TestGLApiInit(t *testing.T) {
	url := "https://gitlab.com"
	token := "my_secret_token"
	verbose := true

	glapi := NewGLApi(url, token, verbose)

	if glapi.UrlBase != url {
		t.Errorf(`TestGLApiInit(urlbase) = %s, want %s`, glapi.UrlBase, url)
	}
	if glapi.Token != token {
		t.Errorf(`TestGLApiInit(token) = %s, want %s`, glapi.Token, token)
	}
	if glapi.Verbose != verbose {
		t.Errorf(`TestGLApiInit(verbose) = %t, want %t`, glapi.Verbose, verbose)
	}
}

func TestGLApiCall(t *testing.T) {
	url := "https://gitlab.com"
	token := ""
	verbose := true
	path := "/api/v4/projects"

	glapi := NewGLApi(url, token, verbose)
	resp, nbPage, err := glapi.CallGitlabApi(path, http.MethodGet, nil)
	if err != nil {
		t.Errorf(`CallGitlabApi(err) = %s`, err)
	}
	if nbPage == 0 {
		t.Errorf(`CallGitlabApi(nbPage) = %d`, nbPage)
	}
	if len(resp) == 0 {
		t.Errorf(`CallGitlabApi(resp len) = %s`, resp)
	}
}
