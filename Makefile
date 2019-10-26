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
	sudo cp docker-compose@.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo cp docker-compose.yml /etc/docker/compose/overmyhouse/
	sudo systemctl enable docker-compose@overmyhouse
	sudo systemctl start docker-compose@overmyhouse
image:
	docker build -t $(IMAGE):$(VERSION) .
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
push-image:
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest
