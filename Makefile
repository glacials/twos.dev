download:
	@go mod download

install-tools: download
	@go list -f '{{range .Imports}}{{.}} {{end}}' cmd/tools.go | xargs go install

build:
	@go build

serve: install-tools
	@gow run . serve

debug: install-tools
	@gow run . serve --debug