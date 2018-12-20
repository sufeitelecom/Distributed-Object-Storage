#!/bin/bash

export RABBITMQ_SERVER='amqp://admin:admin@localhost:5672'
export ES_SERVER=localhost:9200

LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=/tmp/1 go run ./dataserver/main.go &
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=/tmp/2 go run ./dataserver/main.go &
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=/tmp/3 go run ./dataserver/main.go &

LISTEN_ADDRESS=10.29.2.1:12345 go run ./apiserver/main.go &
LISTEN_ADDRESS=10.29.2.2:12345 go run ./apiserver/main.go &