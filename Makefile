run:
	@go run github.com/rafaelsq/wtc

build-docker:
	docker build -t rodrigobrito/gocity .

deploy:
	heroku container:push -a go-city web
	heroku container:release -a go-city web

test-release:
	goreleaser release --snapshot
