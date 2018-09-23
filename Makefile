run:
	docker run -ti -v`pwd`:/go/src/github.com/rodrigo-brito/gocity -p80:3000 -d -w /go/src/github.com/rodrigo-brito/gocity golang go run main.go