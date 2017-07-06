# Makefile for the Docker image upmcenterprises/elasticsearch-cron
# MAINTAINER: Steve Sloka <slokas@upmc.edu>

.PHONY: all build container push clean test

TAG ?= 0.0.1
PREFIX ?= upmcenterprises

all: container

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o _output/bin/elasticsearch-cron --ldflags '-w' ./main.go

container: build
	docker build -t $(PREFIX)/elasticsearch-cron:$(TAG) .

push:
	docker push $(PREFIX)/elasticsearch-cron:$(TAG)

clean:
	rm -f elasticsearch-cron

test: clean
	go test $$(go list ./... | grep -v /vendor/)