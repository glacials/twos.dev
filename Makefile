.PHONY : download build_winter clean build_twos.dev build serve

download:
	go mod download

build_winter:
	go build -o w .

clean:
	@rm -rf dist

build_twos.dev: build_winter
	./w build $(WINTER_ARGS)

build: build_winter clean build_twos.dev

serve: clean
	gow run . serve
