GOCMD=go
BINARY_NAME=server
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
PATH_TO_SERVER_MAIN=./cmd/server
PARAMETES_TO_SERVER=grpc-port 8080

all:	clean run
build: 
	@echo "building server..." 
	@$(GOBUILD) -i -v -o $(BINARY_NAME) $(PATH_TO_SERVER_MAIN)
	
clean: 
	@echo "cleaning the environment..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)

run:
	@echo "starting the server..."
	@$(GOBUILD) -i -v -o $(BINARY_NAME) $(PATH_TO_SERVER_MAIN)
	@./$(BINARY_NAME) --$(PARAMETES_TO_SERVER)