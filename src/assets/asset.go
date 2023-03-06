package assets

import (
	"embed"
)

//go:embed *
var assets embed.FS
var PythonBytes, _ = assets.ReadFile("python_logo.svg")
var NodeBytes, _ = assets.ReadFile("nodejs_logo.svg")
var GoBytes, _ = assets.ReadFile("gopher_logo.svg")
var DockerBytes, _ = assets.ReadFile("docker.svg")
var PostmanBytes, _ = assets.ReadFile("postman_logo.svg")
var NotionBytes, _ = assets.ReadFile("notion_logo.svg")
