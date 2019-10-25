# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=overmyhouse
IMAGE=jsmithedin/overmyhouse
VERSION=master

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
install:
	systemctl stop overmyhouse
	cp overmyhouse /usr/local/overmyhouse/
	systemctl start overmyhouse
	systemctl start overmyhouse-forwarder
image:
	docker build -t $(IMAGE):$(VERSION) .
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
push-image:
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest
