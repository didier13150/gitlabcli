package gitlabcli

import (
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

func readFromFile(filename string, kind string, verbose bool) string {
	if verbose {
		if len(kind) > 0 {
			log.Printf("Try to read %s from %s file", kind, filename)
		} else {
			log.Printf("Try to read %s file", filename)
		}
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		if len(kind) > 0 {
			log.Fatal(err, " and no "+kind+" specified in command line")
		} else {
			log.Fatal(err)
		}
		return ""
	}
	return strings.TrimSpace(string(content))
}

func writeFile(filename string, json []byte, verbose bool) {
	if verbose {
		log.Printf("Try to write %s file", filename)
	}
	err := os.WriteFile(filename, json, 0644)
	if err != nil {
		log.Println("Export to file", filename, err)
		return
	}
}

func getGitUrl(remoteName string, verbose bool) string {
	inidata, err := ini.Load(".git/config")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	section := inidata.Section("remote \"" + remoteName + "\"")
	url := section.Key("url").String()
	if verbose {
		log.Printf("On remote %s, found url %s", remoteName, url)
	}
	return url
}
