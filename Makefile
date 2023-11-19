.PHONY : build build_twos.dev build_winter clean download serve

build: clean test build_winter build_twos.dev

build_twos.dev:

	./w build $(WINTER_ARGS)

build_winter:
	go build -o w .

clean: build_winter
	./w clean

download:
	go mod download

serve: clean
	@echo Building.
	@gow run . serve

test:
	gow test ./...
