package savedstructures

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var TestingDir string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "saved_struct_testing*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	TestingDir = tmpDir
	defaults := TestStruct{}
	defaults.Field1 = "premadeString"
	defaults.Field2 = -1
	defaults.Field3 = true //all default values start at true

	jsonData, err := json.MarshalIndent(defaults, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshel test struct: %v", err)
	}
	err = os.WriteFile(genTmpPath("premade.json"), jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to create premade json: %v", err)
	}

	code := m.Run()

	err = os.RemoveAll(TestingDir)

	if err != nil {
		fmt.Println("Failed to remove testing dir, check /tmp/")
	}

	os.Exit(code)

}

func genTmpPath(path string) string {
	return filepath.Join(TestingDir, path)
}
