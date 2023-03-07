package assets

import _ "embed"

//go:embed python_logo.svg
var PythonBytes []byte

//go:embed nodejs_logo.svg
var NodeBytes []byte

//go:embed gopher_logo.svg
var GoBytes []byte

//go:embed docker.svg
var DockerBytes []byte

//go:embed postman_logo.svg
var PostmanBytes []byte

//go:embed notion_logo.svg
var NotionBytes []byte
