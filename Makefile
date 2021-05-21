.DEFAULT_TARGET=help

.PHONY: help
help:
	@echo
	@echo "Available commands: "
	@echo
	@echo "lint:    lint all files with golangci-lint"
	@echo "test:    test all files"
	@echo "cover:   coverage test all files"
	@echo "example: run an example with APP argument e.g. make example APP=basic"
	@echo

.PHONY: lint
lint:
	$(call blue, "# running linter...")
	@golangci-lint run ./...

.PHONY: test
test:
	$(call blue, "# running tests...")
	@go test -race ./...

.PHONY: cover
cover:
	$(call blue, "# running coverage tests...")
	@go test -cover ./...

.PHONY: example
example:
	$(call blue, "# running example $(APP)...")
	@go run ./_examples/$(APP)

define blue
	@tput setaf 4
	@echo $1
	@tput sgr0
endef