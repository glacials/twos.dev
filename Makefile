download:
	@go mod download

install-tools: download
	@go list -f '{{range .Imports}}{{.}}{{end}}' winter/tools.go | xargs go install

build:
	@go build -o w twos.dev/winter/cmd

serve: install-tools
	@gow run twos.dev/winter/cmd serve

debug: install-tools
	@gow run twos.dev/winter/cmd serve --debug
