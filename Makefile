BINARY_NAME = Toolbox
LOGO_PATH = assets/app_logo.png
APP_ID = milkcoke.falcon-toolbox
APP_NAME = Toolbox
VERSION = 1.0.0

# build: build binary and package app
build:
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -icon ${LOGO_PATH} -appID ${APP_ID} -release

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