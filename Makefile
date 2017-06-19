GO_PKGS=$(shell go list ./... | grep -v vendor)
REPO=aws-tag-dns
TAG=0.1.0

fmt: *.go
	go fmt $(GO_PKGS)

test: fmt
	go test $(GO_PKGS)


test_docker:
	docker run --rm -v $$(pwd):/go/src/github.com/skuid/$(REPO) -w /go/src/github.com/skuid/$(REPO) golang:1.8 go test $(GO_PKGS)

build:
	docker build -t $(REPO):build -f Dockerfile.build .
	docker create --name $(REPO)-extract $(REPO):build
	docker cp $(REPO)-extract:/go/src/github.com/skuid/$(REPO)/$(REPO) .
	docker rm $(REPO)-extract
	docker rmi $(REPO):build
	docker build -t quay.io/skuid/$(REPO):$(TAG) .

push:
	docker push quay.io/skuid/$(REPO):$(TAG)
