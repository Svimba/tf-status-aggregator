SHELL:=/bin/bash
PWD := $(shell pwd)
TAG :="willco/tf-status-aggregator"


.PHONY: build
build: ## Build tungsten-fabric-operator executable file in local go env
	echo "Building tf-status-aggregator bin"
	go build -o build/_output/bin/tf-status-aggregator src/main/*
	echo "Building container"
	docker build -t $(TAG) -f build/Dockerfile .

push:
	docker push $(TAG)


