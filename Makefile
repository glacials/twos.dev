download:
	@go mod download

install-tools: download
	@go list -f '{{range .Imports}}{{.}}{{end}}' winter/tools.go | xargs go install

build:
	@rm -rf dist
	@go build -o w twos.dev/winter/cmd
	@./w build --author "Benjamin Carlsson <ben@twos.dev>" --desc "misc thoughts" --domain twos.dev --name twos.dev --since 2021

serve: install-tools
	@gow run twos.dev/winter/cmd build --serve --author "Benjamin Carlsson <ben@twos.dev>" --desc "misc thoughts" --domain twos.dev --name twos.dev --since 2021
