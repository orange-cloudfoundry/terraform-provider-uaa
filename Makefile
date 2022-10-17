.PHONY: test

test:
	 go test -v -timeout 10m ./test/...
