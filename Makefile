# Define variables
APP_NAME := chat-app
MODULE_NAME := github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern
PROTO_DIR := proto
PROTO_FILE := /Users/bernardbaker/go/src/github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/$(PROTO_DIR)/chat.proto
GRPC_OUT_DIR := .
SRC_DIR := .

AWS_STACK_NAME := chat-app-stack
CLOUDFORMATION_TEMPLATE := cloudformation.yml

GO_FILES := $(shell find . -name '*.go')
BINARY := $(APP_NAME)

# Protoc
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Define the protoc version and download link
PROTOC_VERSION := 25.1
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip

ifeq ($(UNAME_S),Darwin)
    ifeq ($(UNAME_M),arm64)
        PROTOC_ZIP := protoc-$(PROTOC_VERSION)-osx-aarch_64.zip
    else
        PROTOC_ZIP := protoc-$(PROTOC_VERSION)-osx-x86_64.zip
    endif
else
    ifeq ($(UNAME_M),x86_64)
        PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
    else
        PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-aarch_64.zip
    endif
endif

PROTOC_URL := https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)
PROTOC_DIR := /usr/local/protoc

# Commands
GO := go
PROTOC := $(PROTOC_DIR)/bin/protoc
AWS := aws
CLOUDFORMATION := cloudformation

# Default target
.PHONY: all
all: build

# Initialize Go module if not already initialized
.PHONY: init
init:
	@if [ ! -f go.mod ]; then \
		echo "Initializing Go module..."; \
		$(GO) mod init $(MODULE_NAME); \
		$(GO) mod tidy; \
	else \
		echo "Go module already initialized."; \
	fi

# Check if protoc is installed
.PHONY: check-protoc
check-protoc:
	@if ! [ -x "$(PROTOC)" ]; then \
	    echo "protoc not installed. Installing..."; \
	    $(MAKE) install-protoc; \
	    $(MAKE) add-protoc-to-path; \
	else \
	    echo "protoc is already installed"; \
	fi

# Install protoc if not installed
.PHONY: install-protoc
install-protoc:
	@echo "Downloading protoc..."
	@curl -OL $(PROTOC_URL)
	@echo "Unzipping protoc..."
	@sudo mkdir -p $(PROTOC_DIR)
	@sudo chown bernardbaker $(PROTOC_DIR)
	@unzip -o $(PROTOC_ZIP) -d $(PROTOC_DIR)
	@rm -f $(PROTOC_ZIP)
	@echo "protoc installed at $(PROTOC_DIR)/bin"
	@sudo chmod +x $(PROTOC_DIR)/bin/protoc
 

# Add protoc binary to PATH (optional, to ensure it's available in the current session)
.PHONY: add-protoc-to-path
add-protoc-to-path:
	@echo "Adding protoc to PATH..."
	export PATH="$(PROTOC_DIR)/bin:$$PATH"

# Generate gRPC code from proto file
.PHONY: proto
proto: check-protoc
	@echo "Generating gRPC code..."
	@if command -v protoc >/dev/null 2>&1; then \
	$(PROTOC) --proto_path=/Users/bernardbaker/go/src/github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/proto/ --go-grpc_out=$(GRPC_OUT_DIR) $(PROTO_FILE); \
	else \
		echo "protoc not found, falling back to direct protoc call"; \
		go generate ./...; \
	fi
# protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chat.proto 
# $(PROTOC) --proto_path=/Users/bernardbaker/go/src/hexagonal.streamlit.chat/proto/chat.proto --go_out=plugins=grpc:$(GRPC_OUT_DIR) $(PROTO_FILE); \
# proto:
#  @echo "Generating gRPC code from $(PROTO_FILE)"
# $(PROTOC) --go_out=plugins=grpc:$(GRPC_OUT_DIR) $(PROTO_FILE)
# go generate ./...

# Build the Go application
.PHONY: build
build: init proto
	@echo "Building $(BINARY)"
	$(GO) build -o $(BINARY) $(SRC_DIR)

# Run the Go application
.PHONY: run
run: build
	@echo "Running $(BINARY)..."
	./$(BINARY)

# Test the Go application
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./...

# Clean the build and generated files
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY)
	rm -f $(GRPC_OUT_DIR)/*.pb.go
	@echo $(GRPC_OUT_DIR)
	sudo rm -rf $(PROTOC_DIR)

# Deploy AWS infrastructure using CloudFormation
.PHONY: deploy
deploy:
	@echo "Deploying CloudFormation stack..."
	$(AWS) $(CLOUDFORMATION) deploy \
		--template-file $(CLOUDFORMATION_TEMPLATE) \
		--stack-name $(AWS_STACK_NAME) \
		--capabilities CAPABILITY_NAMED_IAM

# Delete the AWS stack
.PHONY: delete-stack
delete-stack:
	@echo "Deleting CloudFormation stack..."
	$(AWS) $(CLOUDFORMATION) delete-stack --stack-name $(AWS_STACK_NAME)

# Show the outputs of the deployed CloudFormation stack
.PHONY: stack-outputs
stack-outputs:
	@echo "Fetching CloudFormation stack outputs..."
	$(AWS) $(CLOUDFORMATION) describe-stacks \
		--stack-name $(AWS_STACK_NAME) \
		--query 'Stacks[0].Outputs'

# Run code linting
.PHONY: lint
lint:
	@echo "Running Go lint..."
	golangci-l
