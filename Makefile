.PHONY: test

test:
	 go test -v -timeout 120m \
 		./uaa_test \
 		./uaa_test/provider_test
