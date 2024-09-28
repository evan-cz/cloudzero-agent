REPO_NAME=ghcr.io/cloudzero/cirrus-stress-app

# Help target to list all available targets with descriptions
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: docker-build
docker-build: ## Build and push the Docker image
	@docker buildx create --use
	@docker buildx build --platform linux/amd64,linux/arm64 -t $(REPO_NAME):latest --push .

.PHONY: build
build: ## Build the stessor locally
	@go build -o cirrus-stress-app .

.PHONY: fmt
fmt: ## Run go fmt against code
	@go fmt ./...

.PHONY: lint
lint: ## Run the linter 
	@golangci-lint run
	
.PHONY: clean
clean: ## Clean the build artifacts
	@rm -f stressor