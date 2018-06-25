.PHONY: build doc fmt lint dev test vet godep install bench

PKG_NAME=$(shell basename `pwd`)

deps:
	dep version || (curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh)
	dep ensure

build: deps
	go build -v -o ./bin/$(PKG_NAME)

sanitize: fmt lint vet

build_linux: deps
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o $(PKG_NAME) -a -installsuffix cgo .

docker: build_linux
	docker-compose -f docker-compose.yaml build khabar

fmt: deps
	go fmt ./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint: deps
	golint ./...

dev:
	DEBUG=* go get && go install && gin -p 8911 -i

test: deps
	go test ./...

# https://godoc.org/golang.org/x/tools/cmd/vet
vet: deps
	go vet ./...

docker_login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

docker_upload: docker_login
	docker-compose -f docker-compose.yaml push khabar
	docker tag $(REPO):latest $(REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker tag $(REPO):latest $(REPO):$(TRAVIS_BRANCH)-latest
	docker push $(REPO):$(TRAVIS_BRANCH)-latest

