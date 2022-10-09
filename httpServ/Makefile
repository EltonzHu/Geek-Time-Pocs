GOCMD=GO
#GOGET=$(GOCMD) get
#GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build

BINDIR=bin
CMDDIR=./cmd

all: build 
build:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "-s -w" -o $(BINDIR)/httpServ -v $(CMDDIR)
clean:
	@rm -rf ./$(BINDIR)