//go:build windows

package app

type InstallerConfig struct {
	Name    string
	Ext     string
	Url     string
	version string
}

var GoInstaller = InstallerConfig{
	"Go",
	".msi",
	"https://go.dev/dl/go1.20.2.windows-amd64.msi",
	"1.20.2",
}
var DockerInstaller = InstallerConfig{
	"Docker",
	".exe",
	"https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe",
	"23.0",
}

var NotionInstaller = InstallerConfig{
	"Notion",
	".exe",
	"https://www.notion.so/desktop/windows/download",
	"2.0.41",
}

var NodeInstaller = InstallerConfig{
	"Nodejs",
	".msi",
	"https://nodejs.org/dist/v18.15.0/node-v18.15.0-x64.msi",
	"v18.15.0",
}

var PostmanInstaller = InstallerConfig{
	"Postman",
	".exe",
	"https://dl.pstmn.io/download/latest/win64",
	"10.12.0",
}

var PythonInstaller = InstallerConfig{
	"Python",
	".exe",
	"https://www.python.org/ftp/python/3.11.3/python-3.11.3-amd64.exe",
	"3.11.3",
}
