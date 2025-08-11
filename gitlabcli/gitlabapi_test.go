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

func TestGLApiCallAllProjects(t *testing.T) {
	url := "http://localhost:8080"
	verbose := false
	path := "/api/v4/projects?page=2&per_page=2"
	expected := `[{"id":3,"name":"Ludwig Van Beethoven","description":null,"path":"ludwig_van_beethoven","name_with_namespace":"CPP Language / Ludwig Van Beethoven","path_with_namespace":"cpp_language/ludwig_van_beethoven","ssh_url_to_repo":"git@localhost:cpp_language/ludwig_van_beethoven.git","http_url_to_repo":"http://localhost:8080/cpp_language/ludwig_van_beethoven.git","web_url":"http://localhost:8080/cpp_language/ludwig_van_beethoven","visibility":"public"},{"id":4,"name":"Johannes Brahms","description":null,"path":"johannes_brahms","name_with_namespace":"CPP Language / Johannes Brahms","path_with_namespace":"cpp_language/johannes_brahms","ssh_url_to_repo":"git@localhost:cpp_language/johannes_brahms.git","http_url_to_repo":"http://localhost:8080/cpp_language/johannes_brahms.git","web_url":"http://localhost:8080/cpp_language/johannes_brahms","visibility":"public"}]`

	glapi := NewGLApi(url, "", verbose)
	resp, nbPage, err := glapi.CallGitlabApi(path, http.MethodGet, nil)
	if err != nil {
		t.Errorf(`CallGitlabApi(err) = %s`, err)
	}
	if nbPage == 0 {
		t.Errorf(`CallGitlabApi(nbPage) = %d`, nbPage)
	}
	if len(resp) == 0 {
		t.Errorf(`CallGitlabApi(resp len) = %d`, len(resp))
	}
	if string(resp) != expected {
		t.Errorf(`CallGitlabApi(resp) = %s`, resp)
	}
}

func TestGLApiCallProject(t *testing.T) {
	url := "http://localhost:8080"
	verbose := false
	path := "/api/v4/projects/3"
	expected := `[{"id":3,"name":"Ludwig Van Beethoven","description":null,"path":"ludwig_van_beethoven","name_with_namespace":"CPP Language / Ludwig Van Beethoven","path_with_namespace":"cpp_language/ludwig_van_beethoven","ssh_url_to_repo":"git@localhost:cpp_language/ludwig_van_beethoven.git","http_url_to_repo":"http://localhost:8080/cpp_language/ludwig_van_beethoven.git","web_url":"http://localhost:8080/cpp_language/ludwig_van_beethoven","visibility":"public"}]`

	glapi := NewGLApi(url, "", verbose)
	resp, _, err := glapi.CallGitlabApi(path, http.MethodGet, nil)
	if err != nil {
		t.Errorf(`CallGitlabApi(err) = %s`, err)
	}
	if len(resp) == 0 {
		t.Errorf(`CallGitlabApi(resp len) = %d`, len(resp))
	}
	if string(resp) != expected {
		t.Errorf(`CallGitlabApi(resp) = %s`, resp)
	}
}

func TestGLApiCallVariables(t *testing.T) {
	url := "http://localhost:8080"
	verbose := false
	path := "/api/v4/projects/3/variables"
	expected := `[{"key":"DEBUG_ENABLED","value":"1","description":null,"environment_scope":"*","raw":true,"hidden":false,"protected":false,"masked":false},{"key":"DEBUG_ENABLED","value":"1","description":null,"environment_scope":"staging","raw":true,"hidden":false,"protected":false,"masked":false},{"key":"DEBUG_ENABLED","value":"1","description":null,"environment_scope":"production","raw":true,"hidden":false,"protected":false,"masked":false}]`

	glapi := NewGLApi(url, "", verbose)
	resp, _, err := glapi.CallGitlabApi(path, http.MethodGet, nil)
	if err != nil {
		t.Errorf(`CallGitlabApi(err) = %s`, err)
	}
	if len(resp) == 0 {
		t.Errorf(`CallGitlabApi(resp len) = %d`, len(resp))
	}
	if string(resp) != expected {
		t.Errorf(`CallGitlabApi(resp) = %s`, resp)
	}
}

func TestGLApiCallEnvironments(t *testing.T) {
	url := "http://localhost:8080"
	verbose := false
	path := "/api/v4/projects/3/environments"
	expected := `[{"id":5,"name":"Staging","state":"available","external_url":null,"description":null},{"id":6,"name":"Production","state":"available","external_url":null,"description":null}]`

	glapi := NewGLApi(url, "", verbose)
	resp, _, err := glapi.CallGitlabApi(path, http.MethodGet, nil)
	if err != nil {
		t.Errorf(`CallGitlabApi(err) = %s`, err)
	}
	if len(resp) == 0 {
		t.Errorf(`CallGitlabApi(resp len) = %d`, len(resp))
	}
	if string(resp) != expected {
		t.Errorf(`CallGitlabApi(resp) = %s`, resp)
	}
}
