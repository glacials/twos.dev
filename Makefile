.PHONY : download build_winter prep_twos.dev build_twos.dev build serve
UNAME := $(shell uname -s)
WINTER_ARGS :=  --author "Benjamin Carlsson <ben@twos.dev>" --desc "misc thoughts" --domain twos.dev --name twos.dev --since 2021

download:
	go mod download

build_winter:
	go build -o w twos.dev/winter/cmd

prep_twos.dev:
	@rm -rf dist
	@curl --silent https://raw.githubusercontent.com/glacials/dotfiles/main/dot_config/emacs/config.org > src/warm/config.org

build_twos.dev:
	./w build $(WINTER_ARGS)

build: build_winter prep_twos.dev build_twos.dev

serve: build_winter prep_twos.dev
	./w build --serve $(WINTER_ARGS)
