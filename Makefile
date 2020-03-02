#GOCMD=go
#GOBUILD=$(GOCMD) build
BINARY_NAME=mybinary
PATH=./cmd/server
  
build:
 go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)