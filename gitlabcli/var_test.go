package gitlabcli

import (
	"net/http"
	"testing"
)

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
