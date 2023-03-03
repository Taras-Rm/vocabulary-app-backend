# Vocabulary application API

This is REST API for Vocabulary application.
The main application file is `main.go` file.

While it is no tests in application.

## Install
    
    go mod download

## Run the app localy

    go run main.go

## Run the Dockerise app

    docker-compose up

## Create migration

    migrate create -ext sql -dir ./migrations/postgres <migration_name>

## Run migration

    migrate -path ./migrations/postgres -database postgres://postgres:post1235@localhost:5432/vocabulary?sslmode=disable up