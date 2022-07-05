download:
	@go mod download

install-tools: download
	@go list -f '{{range .Imports}}{{.}}{{end}}' winter/tools.go | xargs go install

build:
	@go build -o w twos.dev/winter/cmd
	@./w --author "Benjamin Carlsson <ben@twos.dev>" --desc "misc thoughts" --domain twos.dev --name twos.dev --since 2021 build

serve: install-tools
	@gow run twos.dev/winter/cmd --author "Benjamin Carlsson <ben@twos.dev>" --desc "misc thoughts" --domain twos.dev --name twos.dev --since 2021 --serve build
