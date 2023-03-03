package filehandle

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const completeTestFile = "Docker.exe"
const partialTestFile = "Docker.crdownload"
const testDirName = "test_assets"

func getCompleteFileFullPath() (string, error) {
	curDir, err := os.Getwd()
	testDir := filepath.Join(curDir, "..", "..", testDirName)
	testFileFullPath := filepath.Join(testDir, completeTestFile)
	return testFileFullPath, err
}

func getPartialFileFullPath() (string, error) {
	curDir, err := os.Getwd()
	testDir := filepath.Join(curDir, "..", "..", testDirName)
	testFileFullPath := filepath.Join(testDir, partialTestFile)
	return testFileFullPath, err
}

func Test_OpenFile(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	fileFullPath := filepath.Join(currentDir, "..", "..", "test_assets", partialTestFile)
	NavigateToDir(fileFullPath)
}
