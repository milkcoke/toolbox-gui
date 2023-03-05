package filehandle

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
)

func NavigateToDir(fileFullPath string) error {
	// Is it right approach conditional statement with runtime.GOOS ?
	switch runtime.GOOS {
	case "windows":
		cmdName := "explorer"

		cmd := exec.Command(cmdName, "/select,", fileFullPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}

	case "darwin":
		cmdName := "open"

		cmd := exec.Command(cmdName, "-R", fileFullPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	default:
		return errors.New("not supported OS")
	}

	return nil
}
