package gitlabcli

import (
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
