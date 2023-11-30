.PHONY : build serve tools

build: tools
	winter build $(WINTER_ARGS)

serve: tools
	winter serve

tools:
	go install twos.dev/winter@latest
