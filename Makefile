start_host: build_host run_host

build_host:
	CGO_ENABLE=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -modfile go.mod -o bin/host-service -trimpath main/*.go

run_host:
	/bin/bash -c bin/host-service

clean_host:
	rm -rf bin/host-service