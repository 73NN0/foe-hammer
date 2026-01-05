.PHONY: build

all: build

build:
	go build -o foe ./cmd/cli/ && chmod u+x ./foe