HANDLER ?= main
PACKAGE ?= $(HANDLER)

ifeq ($(OS),Windows_NT)
	GOPATH ?= $(USERPROFILE)/go
	GOPATH := /$(subst ;,:/,$(subst \,/,$(subst :,,$(GOPATH))))
	CURDIR := /$(subst :,,$(CURDIR))
	RM := del /q
else
	GOPATH ?= $(HOME)/go
	RM := rm -f
endif

MAKEFILE = $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))

build:
	GOOS=linux go build -ldflags='-w -s $(LDFLAGS)' -o $(HANDLER)

zip:
	$(RM) $(PACKAGE).zip

clean:
	$(RM) $(HANDLER) $(PACKAGE).zip
