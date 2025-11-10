.PHONY: build clean install deb test

# Build variables
BINARY_NAME=multiflexi-tui
BUILD_DIR=bin
MAIN_PATH=./cmd/multiflexi-tui

# Go build settings
GO_BUILD_FLAGS=-trimpath -ldflags="-s -w"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f ../$(BINARY_NAME)_*.deb
	@rm -f ../$(BINARY_NAME)_*.changes
	@rm -f ../$(BINARY_NAME)_*.buildinfo

# Install binary locally (for testing)
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build Debian package
deb: clean
	@echo "Building Debian package..."
	dpkg-buildpackage -us -uc

# Development build (with debug info)
dev:
	@echo "Building $(BINARY_NAME) for development..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Check for required tools
check-deps:
	@which go > /dev/null || (echo "Go is not installed" && exit 1)
	@which dpkg-buildpackage > /dev/null || (echo "dpkg-dev is not installed" && exit 1)
	@echo "All required dependencies are available"