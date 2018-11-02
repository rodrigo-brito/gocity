build:
	yarn build

deploy: build
	cd build && git init || echo "git ok"
	cd build && git remote add deploy git@github.com:go-city/go-city.github.io.git || echo "remote ok"
	cd build && git add .
	cd build && git commit -v --no-edit --amend || git commit -m "deploy" || echo "changes ok"
	cd build && git push deploy master -f