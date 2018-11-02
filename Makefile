dev-dependencies:
	go get -u github.com/canthefason/go-watcher
	go install github.com/canthefason/go-watcher/cmd/watcher

watcher: dev-dependencies
	GOOGLE_APPLICATION_CREDENTIALS=`pwd`/gcs-credentials.json watcher # github.com/canthefason/go-watcher

mock:
	go get github.com/vektra/mockery/...
	mockery -output testdata/mocks -dir ./lib -all
	mockery -output testdata/mocks -dir ./analyzer -all

run:
	docker run -ti -v`pwd`:/go/src/github.com/rodrigo-brito/gocity -e "GOOGLE_APPLICATION_CREDENTIALS=/go/src/github.com/rodrigo-brito/gocity/gcs-credentials.json" -p80:4000 -d -w /go/src/github.com/rodrigo-brito/gocity golang go run main.go

test:
	echo "" > coverage.txt
	for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -v -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out; \
	done