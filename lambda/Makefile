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
	GOOS=linux go build -tags dev -ldflags='-w -s $(LDFLAGS)' -o $(HANDLER)

build-qa:
	GOOS=linux go build -tags qa -ldflags='-w -s $(LDFLAGS)' -o $(HANDLER)

build-prod:
	GOOS=linux go build -tags prod -ldflags='-w -s $(LDFLAGS)' -o $(HANDLER)

zip:
	$(RM) $(PACKAGE).zip
	cp ../template/*.html .
	zip $(PACKAGE).zip $(HANDLER) *.html
	$(RM) *.html

clean:
	$(RM) $(HANDLER) $(PACKAGE).zip
