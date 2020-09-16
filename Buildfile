.PHONY: all downloadImgs runlint migration env db runtests build 

glabel = "\033[92m$(1)\033[0m"

all: downloadImgs runlint runTests build

downloadImgs:
	@echo $(call glabel,"[Downloading images]")
	@docker pull golangci/golangci-lint
	@docker pull postgres:13
	@docker pull septemhill/liquibase:latest

runlint: downloadImgs
	@echo $(call glabel,"[Running golangci-lint check]")
	@docker run -v /home/travis/gopath/src/github.com/septemhill/test:/app -w /app golangci/golangci-lint golangci-lint run

migration: downloadImgs env
	@echo $(call glabel,"[Running SQL migration]")
	@docker run --net=host -v $(PWD)/migration:/liquibase/changelog septemhill/liquibase --logLevel=debug --url=jdbc:postgresql://localhost:5433/postgres --changeLogFile=./changelog/dbchangelog.xml --username=postgres --password=postgres update

db:
	@echo $(call glabel,"Running Postgres docker")
	@docker run -p 5433:5432 -e POSTGRES_PASSWORD=postgres -d postgres

env: db
	@echo $(call glabel,"[Setup environment variable]")
	@export POSTGRES_USER=postgres
	@export POSTGRES_PASSWORD=postgres
	@export POSTGRES_PORT=5433

runtests: migration
	@echo $(call glabel,"[Running test cases]")
	@go test -v ./...

build:
	@echo $(call glabel, "[Running build]")
	@go build -v