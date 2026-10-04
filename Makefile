BINARY_NAME := gatekeeper
IMAGE_NAME  := gatekeeper
IMAGE_TAG   := latest

.PHONY: build run test fmt vet clean docker-build docker-run

build:
	CGO_ENABLED=1 go build -o bin/$(BINARY_NAME) ./cmd/app

run:
	go run ./cmd/app

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run:
	docker run --rm -p 8080:8080 \
		-v $(CURDIR)/configs/user.yaml:/app/configs/user.yaml:rw \
		-v $(CURDIR)/configs/app.local.yaml:/app/configs/app.local.yaml:ro \
		$(IMAGE_NAME):$(IMAGE_TAG)
