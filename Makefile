.PHONY : download build_winter prep_twos.dev build_twos.dev build serve

UNAME := $(shell uname -s)
WINTER_ARGS := --author "Benjamin Carlsson <ben@twos.dev>" --desc "misc thoughts" --domain twos.dev --name twos.dev --since 2021 --source ~/.config/emacs/config.org

download:
	go mod download

build_winter:
	go build -o w .

prep_twos.dev:
	@rm -rf dist

build_twos.dev:
	./w build $(WINTER_ARGS)

build: build_winter prep_twos.dev build_twos.dev

serve: prep_twos.dev
	gow run . build --serve $(WINTER_ARGS)
