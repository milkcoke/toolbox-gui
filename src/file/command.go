package filehandle

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func NavigateToDir(fileFullPath string) {
	// Is it right approach conditional statement with runtime.GOOS ?
	switch runtime.GOOS {
	case "windows":
		cmdName := "explorer"

		cmd := exec.Command(cmdName, fileFullPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal("\nFailed to open downloaded test file! ", err)
		}
	case "darwin":
		cmdName := "open"

		cmd := exec.Command(cmdName, "-R", fileFullPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal("\nFailed to open downloaded test file! ", err)
		}
	default:
		log.Fatal("Not supported OS")
	}
}
