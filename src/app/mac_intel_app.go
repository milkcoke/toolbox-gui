//go:build darwin && amd64

// runtime.GOARCH
// intel: amd64
// M1, M2 :arm64

package app

type InstallerConfig struct {
	Name    string
	Ext     string
	Url     string
	version string
}

var GoInstaller = InstallerConfig{
	"Go",
	".pkg",
	"https://go.dev/dl/go1.22.0.darwin-amd64.pkg",
	"1.22.0",
}
var DockerInstaller = InstallerConfig{
	"Docker",
	".dmg",
	"https://desktop.docker.com/mac/main/amd64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-amd64",
	"25.0.2",
}

var NotionInstaller = InstallerConfig{
	"Notion",
	".dmg",
	"https://www.notion.so/desktop/mac-universal/download",
	"3.1.1",
}

var NodeInstaller = InstallerConfig{
	"Nodejs",
	".pkg",
	"https://nodejs.org/dist/v20.11.1/node-v20.11.1.pkg",
	"v20.11.1",
}

var PostmanInstaller = InstallerConfig{
	"Postman",
	".zip",
	"https://dl.pstmn.io/download/latest/osx_64",
	"10.23",
}

var PythonInstaller = InstallerConfig{
	"Python",
	".pkg",
	"https://www.python.org/ftp/python/3.12.2/python-3.12.2-macos11.pkg",
	"3.12.2",
}

var VSCodeInstaller = InstallerConfig{
	"VisualStudioCode",
	".zip",
	"https://code.visualstudio.com/sha/download?build=stable&os=darwin",
	"1.86.2",
}

var SlackInstaller = InstallerConfig{
	"Slack",
	".dmg",
	"https://downloads.slack-edge.com/releases/macos/4.36.140/prod/universal/Slack-4.36.140-macOS.dmg",
	"4.36.140",
}
