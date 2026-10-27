.PHONY : build serve tools

WINTER_VERSION ?= v0.7.0

build: tools
	winter build $(WINTER_ARGS)

serve: tools
	winter serve

tools:
	go install twos.dev/winter@$(WINTER_VERSION)
