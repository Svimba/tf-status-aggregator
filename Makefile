SHELL:=/bin/bash
PWD := $(shell pwd)
TAG :="willco/tf-status-agregator"


.PHONY: build
build: ## Build tungsten-fabric-operator executable file in local go env
	echo "Building tf-status-agregator bin"
	go build -o build/_output/bin/tf-status-agregator src/main/*
	echo "Building container"
	docker build -t $(TAG) -f build/Dockerfile .

push:
	docker push $(TAG)


