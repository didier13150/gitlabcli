package gitlabcli

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
)

type GLApi struct {
	UrlBase string
	Token   string
	Verbose bool
}

func NewGLApi(urlBase string, token string, verbose bool) GLApi {
	glapi := GLApi{urlBase, token, verbose}
	return glapi
}

func (glapi *GLApi) CallGitlabApi(uri string, method string, data []byte) ([]byte, int, error) {
	nbPage := 0
	gitlabUrl := glapi.UrlBase + uri
	//log.Printf("gitlabUrl: \"%s\"", glapi.UrlBase)
	if glapi.Verbose {
		log.Printf("Querying URL: \"%s\"", gitlabUrl)
	}
	req, err := http.NewRequest(method, gitlabUrl, bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if glapi.Token != "" {
		req.Header.Set("Private-Token", glapi.Token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, 0, err
	}
	if len(res.Header["X-Total-Pages"]) > 0 {
		nbPage, _ = strconv.Atoi(res.Header["X-Total-Pages"][0])
		currrentPage, _ := strconv.Atoi(res.Header["X-Page"][0])
		nbItems, _ := strconv.Atoi(res.Header["X-Total"][0])
		if glapi.Verbose {
			log.Printf("Pagination - Reading page %d/%d - There is %d item(s) to download", currrentPage, nbPage, nbItems)
		}
	} else if len(res.Header["X-Page"]) > 0 {
		nbPage, _ = strconv.Atoi(res.Header["X-Page"][0])
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln("Cannot close body", err)
		}
	}()

	json, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and body: %s\n", res.StatusCode, json)
	}
	if err != nil {
		log.Fatal(err)
	}
	return json, nbPage, nil
}
