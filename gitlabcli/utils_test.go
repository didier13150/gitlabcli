package gitlabcli

import (
	"os"
	"strings"
	"testing"
)

func TestUtilsReadFile(t *testing.T) {

	expected := "  \nline 1\nline 2\nline 3\n "
	filename := "/tmp/glcli-readfile.txt"
	err := os.WriteFile(filename, []byte(expected), 0644)
	if err != nil {
		t.Errorf(`TestUtilsReadFile(write test file) = %s`, err)
	}
	content := readFromFile(filename, "test", false)
	if content != strings.TrimSpace(expected) {
		t.Errorf(`TestUtilsReadFile(read test file) = %s`, err)
	}
	err = os.Remove(filename)
	if err != nil {
		t.Errorf(`TestUtilsReadFile(delete test file) = %s`, err)
	}
}

func TestUtilsWriteFile(t *testing.T) {

	expected := "  \nline 1\nline 2\nline 3\n "
	filename := "/tmp/glcli-readfile.txt"

	writeFile(filename, []byte(expected), false)

	content, err := os.ReadFile(filename)
	if string(content) != expected {
		t.Errorf(`TestUtilsWriteFile(read test file) = %s`, err)
	}
	err = os.Remove(filename)
	if err != nil {
		t.Errorf(`TestUtilsWriteFile(delete test file) = %s`, err)
	}
}
