.PHONY:

lint:
	docker run --rm -v "$(CURDIR):/app" -w /app golangci/golangci-lint:v1.47.2 golangci-lint run
