dev-dependencies:
	go get -u github.com/canthefason/go-watcher
	go install github.com/canthefason/go-watcher/cmd/watcher

watcher: dev-dependencies
	GOOGLE_APPLICATION_CREDENTIALS=`pwd`/gcs-credentials.json watcher # github.com/canthefason/go-watcher

run:
	docker run -ti -v`pwd`:/go/src/github.com/rodrigo-brito/gocity -p80:4000 -d -w /go/src/github.com/rodrigo-brito/gocity golang go run main.go