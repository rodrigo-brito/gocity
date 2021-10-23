run:
	@go run github.com/rafaelsq/wtc

mock:
	go get github.com/vektra/mockery/...
	mockery -output testdata/mocks -dir ./lib -all
	mockery -output testdata/mocks -dir ./analyzer -all

build-docker:
	docker build -t rodrigobrito/gocity .

deploy:
	heroku container:push -a go-city web
	heroku container:release -a go-city web


test-release:
	goreleaser release --snapshot

# This is a workaround for the fact that assets are currently being
# maintained in a different branch.  How to build:
#
# git clone $repo_url master
# git clone $repo_url front-end
# cd master
# git checkout $master_fix_branch
# cd ../front-end
# git checkout $front_end_fix_branch
# nvm-load
# npm install yarn 
# npm install 
# cd ../master
# make build-with-assets
#
build-with-assets:
	cd ../front-end && make -B build
	cp -a ../front-end/build/* handle/assets/
	go build
