fmt:
	@gofmt -s -w ./
lint:
	@golangci-lint run --disable-all -E errcheck