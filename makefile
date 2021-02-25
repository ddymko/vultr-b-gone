build:
	go build -o vultr-b-gone

format:
	@go fmt ./...

lint:
	@go list ./... | grep -v /vendor/ | xargs -n 1 golint