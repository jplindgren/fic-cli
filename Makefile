PROJECT_NAME=fic-cli
BINARY_NAME=fic

run:
	go run ./cmd -spreadsheetid=${SPREADSHEET_ID} -credential=${CREDENTIAL}

test:
	go test ./...

build:
	@echo "Building $(PROJECT_NAME)"
	go build -o $(BINARY_NAME) ./cmd
	chmod +x $(BINARY_NAME)

