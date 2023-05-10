//go:build darwin && arm64

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
	"https://go.dev/dl/go1.20.4.darwin-arm64.pkg",
	"1.20.4",
}
var DockerInstaller = InstallerConfig{
	"Docker",
	".dmg",
	"https://desktop.docker.com/mac/main/arm64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-arm64",
	"23.0",
}

var NotionInstaller = InstallerConfig{
	"Notion",
	".dmg",
	"https://www.notion.so/desktop/apple-silicon/download",
	"2.1.13",
}

var NodeInstaller = InstallerConfig{
	"Nodejs",
	".pkg",
	// pkg supports both intel, apple-silicon
	// It contains a universal binary that includes both architectures.
	"https://nodejs.org/dist/v18.16.0/node-v18.16.0.pkg",
	"v18.16.0",
}

var PostmanInstaller = InstallerConfig{
	"Postman",
	".zip",
	"https://dl.pstmn.io/download/latest/osx_arm64",
	"10.11.1",
}

var PythonInstaller = InstallerConfig{
	"Python",
	".pkg",
	"https://www.python.org/ftp/python/3.11.3/python-3.11.3-macos11.pkg",
	"3.11.3",
}
var VSCodeInstaller = InstallerConfig{
	"VisualStudioCode",
	".zip",
	"https://code.visualstudio.com/sha/download?build=stable&os=darwin-arm64",
	"1.78.1",
}
var SlackInstaller = InstallerConfig{
	"Slack",
	".dmg",
	"https://downloads.slack-edge.com/releases/macos/4.32.122/prod/universal/Slack-4.32.122-macOS.dmg",
	"4.32.122",
}
