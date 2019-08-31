dev-dependencies:
	go get -u github.com/canthefason/go-watcher
	go install github.com/canthefason/go-watcher/cmd/watcher

watcher: dev-dependencies
	watcher # github.com/canthefason/go-watcher

mock:
	go get github.com/vektra/mockery/...
	mockery -output testdata/mocks -dir ./lib -all
	mockery -output testdata/mocks -dir ./analyzer -all

test:
	echo "" > coverage.txt
	for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -v -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out; \
	done

build-docker:
	docker build -t rodrigobrito/gocity .

deploy:
	heroku container:push -a go-city web
	heroku container:release -a go-city web