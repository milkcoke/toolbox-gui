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
	"https://go.dev/dl/go1.20.1.darwin-arm64.pkg",
	"1.20.1",
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
	"https://nodejs.org/dist/v18.14.2/node-v18.14.2.pkg",
	"v18.14.2",
}

var PostmanInstaller = InstallerConfig{
	"Postman",
	".zip",
	"https://dl.pstmn.io/download/latest/osx_arm64",
	"10.10.8",
}

var PythonInstaller = InstallerConfig{
	"Python",
	".pkg",
	"https://www.python.org/ftp/python/3.11.3/python-3.11.3-macos11.pkg",
	"3.11.3",
}
