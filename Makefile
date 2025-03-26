all: get-deps unit

SDK_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")

# Define PRINT_OK macro for success messages
PRINT_OK = echo "\033[0;32m✓\033[0m"

get-deps: get-deps-tests get-deps-verify
	@echo "go get SDK dependencies"
	@go get -v $(SDK_ONLY_PKGS)

get-deps-tests:
	@echo "go get SDK testing dependencies"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go get -v github.com/smartystreets/goconvey

get-deps-verify:
	@echo "go get SDK verification utilities"
	@if [ \( -z "${SDK_GO_1_4}" \) -a \( -z "${SDK_GO_1_5}" \) ]; then go get -v golang.org/x/lint/golint; else echo "skipped getting golint"; fi

lint:
	@echo "go lint SDK and vendor packages"
	@golangci-lint run --disable-all --enable=gofmt --enable=golint --enable=vet --enable=unconvert --deadline=4m ${SDK_ONLY_PKGS}
	@$(PRINT_OK)

unit:
	@echo "go test SDK and vendor packages"
	@go test -v $(SDK_ONLY_PKGS)
	@$(PRINT_OK)
