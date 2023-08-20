.PHONY:  bench clean-mod clean-test clean-all test tidy

bench:
	go test -bench=.

build:
	go build -race .

clean-mod:
	go clean -modcache

clean-test:
	go clean -testcache

clean-all: clean-mod clean-test

pkg:
	curl https://sum.golang.org/lookup/github.com/egibs/deepwalk@v$(shell cat VERSION)

sbom:
	syft . -o json > deepwalk_sbom.json

test:
	go test ./... -v

tidy:
	go mod tidy
