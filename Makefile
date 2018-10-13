dev-dependencies:
	go get github.com/canthefason/go-watcher

watcher: dev-dependencies
	watcher # github.com/canthefason/go-watcher

run:
	docker run -ti -v`pwd`:/go/src/github.com/rodrigo-brito/gocity -p80:4000 -d -w /go/src/github.com/rodrigo-brito/gocity golang go run main.go