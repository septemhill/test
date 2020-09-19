.PHONY: all downloadImgs lint migration env testdb proddb tests build 

glabel = "\033[92m$(1)\033[0m"

all: downloadImgs lint tests build

downloadImgs:
	@echo $(call glabel,"[Downloading images]")
	@docker pull golangci/golangci-lint
	@docker pull postgres:13
	@docker pull septemhill/liquibase:latest
	@docker pull redis

lint: downloadImgs
	@echo $(call glabel,"[Running golangci-lint check]")
	@docker run -v $(PWD):/app -w /app golangci/golangci-lint golangci-lint run

migration: downloadImgs env
	@echo $(call glabel,"[Running SQL migration]")
	@docker run --net=host -v $(PWD)/migration:/liquibase/changelog septemhill/liquibase --logLevel=debug --url=jdbc:postgresql://localhost:5433/postgres --changeLogFile=./changelog/dbchangelog.xml --username=postgres --password=postgres update

testdb:
	@echo $(call glabel,"Running test env database dockers")
	@docker run --name test_postgres -p 5433:5432 -e POSTGRES_PASSWORD=postgres -d postgres
	@docker run --name test_redis -p 6380:6379 -d redis

env: testdb
	@echo $(call glabel,"[Setup environment variable]")
	@export POSTGRES_USER=postgres
	@export POSTGRES_PASSWORD=postgres
	@export POSTGRES_PORT=5433

tests: migration
	@echo $(call glabel,"[Running test cases]")
	@go test -v ./...

build:
	@echo $(call glabel, "[Running build]")
	@go build -v

proddb:
	@echo $(call glabel, "[Running production env database dockers]")
	@docker run --name prod_postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
	@docker run --name prod_redis -p 6379:6379 -d redis

run: downloadImgs build proddb
	./test 3000