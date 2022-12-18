run:
	@go run github.com/rafaelsq/wtc

build-docker:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s"
	docker buildx build --platform linux/arm64 --push -t rodrigobrito/gocity .

test-release:
	goreleaser release --snapshot
