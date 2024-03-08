help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

functions := $(shell find functions -name \*main.go | awk -F'/' '{print $$2}')

build: ## Test build golang binary - actual builds performed by serverless-go-plugin
	@for function in $(functions) ; do \
		env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bin/$$function/bootstrap functions/$$function/main.go ; \
	done