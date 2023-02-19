BINARY_NAME = falcon_markdown.app
LOGO_PATH = assets/app_logo.png
APP_NAME = Markdown.app
VERSION = 1.0.0

# build: build binary and package app
build:
	rm -rf ${APP_NAME}
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -icon ${LOGO_PATH} -release

# run: builds and runs the application
run:
	go run .

# clean: Runs go clean and deletes binaries
clean:
	@echo "Cleaning.."
	@go clean
	@rm -rf ${APP_NAME}
	@echo "Cleaned!"

# test: Runs all tests
test:
	go test -v ./...