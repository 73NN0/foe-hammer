.PHONY: build

all: build

build:
	go build -o foe ./cmd/cli/ && chmod u+x ./foe

build-config:
	go build -o foe-config ./internal/config/ && chmod u+x ./foe-config