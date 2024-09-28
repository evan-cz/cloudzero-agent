# Directories
VIRTUALENV := venv
REQUIREMENTS := requirements.txt

# Colors
ERROR_COLOR = \033[1;31m
INFO_COLOR = \033[1;32m
WARN_COLOR = \033[1;33m
NO_COLOR = \033[0m

# Help target to list all available targets with descriptions
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Checks your active version of python
.PHONY: check-python-version
PYTHON_VERSION_SUPPORTED := 3.9 3.10 3.11 3.12
PYTHON_VERSION_SHELL := $(shell python --version | cut -f2 -d' ' | cut -f1,2 -d'.')
check-python-version: 
ifeq ($(findstring $(PYTHON_VERSION_SHELL),$(PYTHON_VERSION_SUPPORTED)),)
	@printf "$(WARN_COLOR)You have python $(ERROR_COLOR)$(PYTHON_VERSION_SHELL)$(WARN_COLOR); you need python one of $(INFO_COLOR)($(PYTHON_VERSION_SUPPORTED))$(NO_COLOR).\n"
	@exit 1
endif

# Create a virtual environment
.PHONY: $(VIRTUALENV)
$(VIRTUALENV): check-python-version ## creates a virtual environment
	@python -m venv $(VIRTUALENV)

# Initialize the virtual environment
PYTHON_ARGS := . $(VIRTUALENV)/bin/activate && PYTHONPATH=.
LIBS := $(VIRTUALENV)/lib/python$(PYTHON_VERSION_SHELL)/site-packages
.PHONY: init
init: $(LIBS) ## initializes a virtual environment, and installs the requirements
$(LIBS): $(VIRTUALENV) $(REQUIREMENTS)
	@$(PYTHON_ARGS) pip install -r $(REQUIREMENTS)
	@touch $(LIBS)


# Clean up any dangling images
.PHONY: clean
clean: ## Clean up dangling Docker images
	@rm -fr $(VIRTUALENV) __pycache__ .ruff_cache

# Run unit tests
.PHONY: fmt
fmt: init ## runs code formatter
	@$(PYTHON_ARGS) black -l 100 ./*.py

.PHONY: fmt-check
fmt-check: fmt ## Check if the code is properly formatted
	$(eval CHANGED := $(shell git status -s | grep "\.py" | wc -c))
	@[ $(CHANGED) -eq 0 ] || { echo "changed files"; exit 1; }

.PHONY: lint
lint: init ## runs code linter
	@$(PYTHON_ARGS) ruff check .  --config ../ruff.toml

.PHONY: check
check: fmt-check lint ## Run code linter and unit tests