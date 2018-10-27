dev-dependencies:
	go get -u github.com/canthefason/go-watcher
	go install github.com/canthefason/go-watcher/cmd/watcher

watcher: dev-dependencies
	GOOGLE_APPLICATION_CREDENTIALS=`pwd`/gcs-credentials.json watcher # github.com/canthefason/go-watcher

build-static:
	cd ui && yarn build

deploy-static:
	cd ui/build && git init || echo "git ok"
	cd ui/build && git remote add deploy git@github.com:go-city/go-city.github.io.git || echo "remote ok"
	cd ui/build && git add .
	cd ui/build && git commit -v --no-edit --amend || git commit -m "deploy" || echo "changes ok"
	cd ui/build && git push deploy master -f

mock:
	go get github.com/vektra/mockery/...
	mockery -output testdata/mocks -dir ./lib -all
	mockery -output testdata/mocks -dir ./analyzer -all

run:
	docker run -ti -v`pwd`:/go/src/github.com/rodrigo-brito/gocity -e "GOOGLE_APPLICATION_CREDENTIALS=/go/src/github.com/rodrigo-brito/gocity/gcs-credentials.json" -p80:4000 -d -w /go/src/github.com/rodrigo-brito/gocity golang go run main.go