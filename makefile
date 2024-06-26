.PHONY: bench build clean-all clean-mod clean-test fmt pkg sbom test tidy

bench:
	go test ./... -bench=. -benchmem

build:
	CGO_ENABLED=0 go build -o deepwalk .

build-docker:
	docker buildx build -f Dockerfile . --load -t deepwalk:latest

clean-mod:
	go clean -modcache

clean-test:
	go clean -testcache

clean-all: clean-mod clean-test

fmt:
	go fmt ./...

pkg:
	curl https://sum.golang.org/lookup/github.com/egibs/deepwalk@v$(shell cat VERSION)

sbom:
	syft . -o json | jq . > deepwalk_sbom.json

test:
	go test ./... -v

tidy:
	go mod tidy
