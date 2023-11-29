.PHONY : build build_twos.dev clean install_winter serve

build: clean test install_winter build_twos.dev

build_twos.dev: install_winter
	winter build $(WINTER_ARGS)

clean: install_winter
	winter clean

install_winter:
	go install twos.dev/winter@latest

serve: clean
	@echo Building.
	@gow run twos.dev/winter serve

test:
	gow test twos.dev/winter/...
