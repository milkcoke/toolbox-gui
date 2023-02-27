package app

type installer struct {
	Name    string
	Ext     string
	Url     string
	version string
}

var GoInstaller = installer{
	"Go",
	".msi",
	"https://go.dev/dl/go1.20.1.windows-amd64.msi",
	"1.20.1",
}
var DockerInstaller = installer{
	"Docker",
	".exe",
	"https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe",
	"23.0",
}

var NotionInstaller = installer{
	"Notion",
	".exe",
	"https://www.notion.so/desktop/windows/download",
	"2.21",
}

var NodeInstaller = installer{
	"Nodejs",
	".msi",
	"https://nodejs.org/dist/v18.14.2/node-v18.14.2-x64.msi",
	"v18.14.2",
}
