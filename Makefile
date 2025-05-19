# Name of the binary file
BINARY_NAME=web-page-analyze-service
coverfile=coverage.out
consolefile=console.log

# Default target
all: clean run

# Build the Go binary
build:
	@scripts/build.sh ${BINARY_NAME}

# Clean up binary
clean:
	@echo "Cleaning up executed..."
	@rm -f $(BINARY_NAME)
	@echo "$(BINARY_NAME) removed..."
	@rm -f $(coverfile)
	@echo "$(coverfile) removed..."
	@rm -f $(consolefile)
	@echo "$(consolefile) removed..."

# Run the binary
run: build
	@./$(BINARY_NAME)