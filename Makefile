download:
	@go mod download

install-tools: download
	@go list -f '{{range .Imports}}{{.}} {{end}}' winter/tools.go | xargs go install

build:
	@cd winter; go build

serve: install-tools
	@gow run . serve

debug: install-tools
	@gow run . serve --debug
