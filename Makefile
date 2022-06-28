download:
	@go mod download

install-tools: download
	@go list -f '{{range .Imports}}{{.}}{{end}}' winter/tools.go | xargs go install

build:
	@go build -o w twos.dev/winter

serve: install-tools
	@cd winter && gow run . serve

debug: install-tools
	@cd winter && gow run . serve --debug
