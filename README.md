# Vocabulary app API

This is REST API for Vocabulary application written on Go. Also it uses Elasticsearch and Postgresql as storages. The main application file is `main.go` file. Deployment branch is `new`.

## App link

This is demo application link: https://d2aps8c8tddxti.cloudfront.net

## Technologies

- Go
- Gin
- Elasticsearch
- Postgresql

## Install
    
    go mod download
You also should add `docker-compose.yml` file for running elasticsearch and postgres localy. Here is an example of my docker-compose:
```
version: '3'

services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch-oss:7.10.2
    environment:
      - discovery.type=single-node
    ports:
      - 9200:9200
    volumes:
      - esdata:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana-oss:7.10.2
    ports:
      - 5601:5601
    depends_on:
      - elasticsearch
    
  postgres:
    image: postgres
    environment:
      POSTGRES_DB: 
      POSTGRES_USER: 
      POSTGRES_PASSWORD: 
    ports:
      - 5432:5432
    volumes:
      - pgdata:/var/lib/postgresql/data/

volumes:
  esdata:
  pgdata:
```

## Run the app localy

    go run main.go


## Create migration

    migrate create -ext sql -dir ./migrations/postgres <migration_name>

## Run migration

    migrate -path ./migrations/postgres -database postgres://<username>:<password>@localhost:5432/<dbname>?sslmode=disable up