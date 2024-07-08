include .envrc

run:
	go run ./cmd -spreadsheetid=${SPREADSHEET_ID} -credential=${CREDENTIAL}

test:
	go test ./...
