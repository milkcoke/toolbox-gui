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
	"https://go.dev/dl/go1.20.4.windows-amd64.msi",
	"1.20.4",
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
	"https://nodejs.org/dist/v18.16.0/node-v18.16.0-x64.msi",
	"v18.16.0",
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

//TODO : mac git installer is different way
//var GitInstaller = InstallerConfig{
//	"Git",
//	".exe",
//	"https://github.com/git-for-windows/git/releases/download/v2.40.1.windows.1/Git-2.40.1-64-bit.exe",
//	"2.40.1",
//}

var VSCodeInstaller = InstallerConfig{
	"VisualStudioCode",
	".exe",
	"https://code.visualstudio.com/sha/download?build=stable&os=win32-x64-user",
	"1.78.1",
}

var SlackInstaller = InstallerConfig{
	"Slack",
	".exe",
	"https://downloads.slack-edge.com/releases/windows/4.32.122/prod/x64/SlackSetup.exe",
	"4.32.122",
}
