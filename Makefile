UNAME := $(shell uname -s)

download:
	@go mod download

build:
	@rm -rf dist
	@go build -o w twos.dev/winter/cmd
	@./w build --author "Benjamin Carlsson <ben@twos.dev>" --desc "misc thoughts" --domain twos.dev --name twos.dev --since 2021

serve:
	go build -o w twos.dev/winter/cmd
	./w build --serve
