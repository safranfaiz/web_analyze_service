#! /bin/bash

artifact=$1
coverfile=coverage.out
consolefile=console.log

# Format code using golang stand
echo "Format code using golang stand..."
go fmt ./...
echo "Formatting done."

# Clean the code
echo "Cleaning the codebase..."
go clean
echo "Codebase cleaning done."

# Remove the existing generated coverage.out file if exists
echo "Removing existing $coverfile file if exists..."
if [ -f "$coverfile" ]; then
  rm -f "$coverfile"
  echo "$coverfile removed"
else
  echo "$coverfile not exists"
fi

# Remove the existing generated console.log file if exists
echo "Removing existing $consolefile file if exists..."
if [ -f "$consolefile" ]; then
  rm -f "$consolefile"
  echo "$consolefile removed"
else
  echo "$consolefile not exists"
fi

# Running unit test cases and write result to console.log file...
echo "Running unit test cases and write result to $consolefile file..."
go test ./... -v -cover -coverpkg=./... -coverprofile=./$coverfile ./... > $consolefile
echo "Unit test cases execution done."

# Build golang application
echo "Build the application..."

# Checking the already created artifact is exist the remove...
if [ -f "$artifact" ]; then
  rm -f "$artifact"
  echo "$artifact removed"
else
  echo "$artifact not exists"
fi

go build -o $artifact
if [[ $? != 0 ]]; then
  echo "Build Failed."
  exit 1
fi
echo "Build passed with no errors..."