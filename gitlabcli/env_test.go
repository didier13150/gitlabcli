package gitlabcli

import (
	"net/http"
	"testing"
)

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
