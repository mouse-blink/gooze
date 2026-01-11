name := gooze
bin := ./.bin
# Versions
COBRA_CLI_VERSION := v1.3.0
GOLANGCI_LINT_VERSION := v2.8.0
MOCKERY_VERSION := v2.53.5

# Whitelisted packages (exclude examples explicitly)
PKG_WHITELIST :=  ./cmd/... ./internal/...

.PHONY: all install-tools build lint test clean run fmt mocks clean-mocks

all: build

install-tools:
	@echo "Installing development tools into $(bin)..."
	@mkdir -p $(bin)
	@echo "Installing cobra-cli $(COBRA_CLI_VERSION)..."
	@GOBIN=$(abspath $(bin)) go install github.com/spf13/cobra-cli@$(COBRA_CLI_VERSION)
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(abspath $(bin)) $(GOLANGCI_LINT_VERSION)
	@echo "Installing mockery $(MOCKERY_VERSION)..."
	@GOBIN=$(abspath $(bin)) go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)


build:
	@go build -o $(bin)/$(name) main.go
	@echo "Built $(name) binary at $(PWD)/$(bin)/$(name)"

lint:
	@echo "Running golangci-lint..."
	@$(bin)/golangci-lint run $(PKG_WHITELIST)

test:
	@go test -v $(PKG_WHITELIST)

clean:
	@rm -rf $(bin)

run: build
	@$(bin)/$(name) $$(echo "$(filter-out $@,$(MAKECMDGOALS))" | sed 's/^-/-/')

%:
	@:

fmt:
	@go fmt $(PKG_WHITELIST)
	@$(bin)/golangci-lint fmt $(PKG_WHITELIST)

clean-mocks:
	@echo "Cleaning mocks..."
	@rm -rf internal/*/mocks

mocks: clean-mocks
	@echo "Generating mocks..."
	@$(bin)/mockery --all --config .mockery.yaml