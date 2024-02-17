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
	"https://go.dev/dl/go1.22.0.windows-amd64.msi",
	"1.22.0",
}
var DockerInstaller = InstallerConfig{
	"Docker",
	".exe",
	"https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe",
	"4.27.2",
}

var NotionInstaller = InstallerConfig{
	"Notion",
	".exe",
	"https://www.notion.so/desktop/windows/download",
	"3.1.1",
}

var NodeInstaller = InstallerConfig{
	"Nodejs",
	".msi",
	"https://nodejs.org/dist/v20.11.1/node-v20.11.1-x64.msi",
	"v20.11.1",
}

var PostmanInstaller = InstallerConfig{
	"Postman",
	".exe",
	"https://dl.pstmn.io/download/latest/win64",
	"10.23",
}

var PythonInstaller = InstallerConfig{
	"Python",
	".exe",
	"https://www.python.org/ftp/python/3.12.2/python-3.12.2-amd64.exe",
	"3.12.2",
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
	"1.86.2",
}

var SlackInstaller = InstallerConfig{
	"Slack",
	".exe",
	"https://downloads.slack-edge.com/releases/windows/4.36.140/prod/x64/SlackSetup.exe",
	"4.36.140",
}
