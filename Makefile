.DEFAULT_GOAL = build

# Builds the binary to cwd
build: 
	@go build -o qosm .

# Builds and installs the binary to /usr/loca/bin
install:
	@go build -ldflags="-s -w" -o qosm .
	@mv ./qosm /usr/local/bin

PROTO_DIR := ./internal/protobuf

PROTO_FILES := qosm.proto

PROTO_PATHS := $(addprefix $(PROTO_DIR)/, $(PROTO_FILES))

#Compiles .proto files and generates go code to work protocol buffers.
proto: $(PROTO_PATHS)
	@echo "Compiling .proto files."
	protoc --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go_opt=paths=source_relative $^



SERVICES_FILE := ./internal/service/iana_ports.go

#Generates the file that has services. (Currently not used)
services: | $(SERVICES_FILE)
	@go run ./internal/service/genservices.go

#Creates the service file
$(SERVICES_FILE):
	@touch $@

#Genereates tailwind css classes that are used.
css:
	@npm install tailwindcss
	@npx @tailwindcss/cli  -i web/static/css/input.css -o web/static/css/output.css
	@rm -rf package.json package-lock.json node_modules
