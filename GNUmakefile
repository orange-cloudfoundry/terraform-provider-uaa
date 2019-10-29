release:
	rm -fr bin
	mkdir -p bin
	GOARCH=amd64 GOOS=windows go build -o bin/terraform-provider-uaa_windows_amd64.exe
	GOARCH=amd64 GOOS=linux   go build -o bin/terraform-provider-uaa_linux_amd64
	GOARCH=amd64 GOOS=darwin  go build -o bin/terraform-provider-uaa_darwin_amd64
