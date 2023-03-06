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
	"https://go.dev/dl/go1.20.1.windows-amd64.msi",
	"1.20.1",
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
	"2.21",
}

var NodeInstaller = InstallerConfig{
	"Nodejs",
	".msi",
	"https://nodejs.org/dist/v18.14.2/node-v18.14.2-x64.msi",
	"v18.14.2",
}

var PostmanInstaller = InstallerConfig{
	"Postman",
	".exe",
	"https://dl.pstmn.io/download/latest/win64",
	"10.10.8",
}

var PythonInstaller = InstallerConfig{
	"Python",
	".exe",
	"https://www.python.org/ftp/python/3.11.2/python-3.11.2-amd64.exe",
	"3.11.2",
}
