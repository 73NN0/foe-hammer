.PHONY: build

all: build

build: main.go
	go build -o foe main.go && chmod u+x ./foe
