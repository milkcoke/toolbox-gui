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
	"https://go.dev/dl/go1.25.6.darwin-arm64.pkg",
	"1.26.2",
}
var DockerInstaller = InstallerConfig{
	"Docker",
	".dmg",
	"https://desktop.docker.com/mac/main/arm64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-arm64",
	"29.4.0",
}

var NotionInstaller = InstallerConfig{
	"Notion",
	".dmg",
	"https://www.notion.so/desktop/mac-universal/download",
	"7.1.0",
}

var NodeInstaller = InstallerConfig{
	"Nodejs",
	".pkg",
	// pkg supports both intel, apple-silicon
	// It contains a universal binary that includes both architectures.
	"https://nodejs.org/dist/v24.15.0/node-v24.15.0.pkg",
	"v24.15.0",
}

var PostmanInstaller = InstallerConfig{
	"Postman",
	".zip",
	"https://dl.pstmn.io/download/latest/osx_arm64",
	"12.6.3",
}

var PythonInstaller = InstallerConfig{
	"Python",
	".pkg",
	"https://www.python.org/ftp/python/3.14.4/python-3.14.4-macos11.pkg",
	"3.14.4",
}
var VSCodeInstaller = InstallerConfig{
	"VisualStudioCode",
	".zip",
	"https://code.visualstudio.com/sha/download?build=stable&os=darwin-arm64",
	"1.116",
}
var SlackInstaller = InstallerConfig{
	"Slack",
	".dmg",
	"https://downloads.slack-edge.com/desktop-releases/mac/universal/4.49.81/Slack-4.49.81-macOS.dmg",
	"4.49.81",
}
