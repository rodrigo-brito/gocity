dev-dependencies:
	go get -u github.com/canthefason/go-watcher
	go install github.com/canthefason/go-watcher/cmd/watcher

watcher: dev-dependencies
	watcher # github.com/canthefason/go-watcher

mock:
	go get github.com/vektra/mockery/...
	mockery -output testdata/mocks -dir ./lib -all
	mockery -output testdata/mocks -dir ./analyzer -all

build-docker:
	docker build -t rodrigobrito/gocity .

deploy:
	heroku container:push -a go-city web
	heroku container:release -a go-city web