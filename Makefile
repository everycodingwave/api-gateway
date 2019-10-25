.PHONY: default dep build test

BINARY = api-gw
GOCMD = go
MAIN = cmd/main.go
#BUILDENV = GOOS=linux

default: dep build test

dep:
	$(GOCMD) get -d ./...

build:
	$(BUILDENV) $(GOCMD) build -o ${BINARY} ${MAIN}

test:
	$(GOCMD) vet ./...
	$(GOCMD) test -cover ./...
