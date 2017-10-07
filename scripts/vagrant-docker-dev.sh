#!/usr/bin/env bash

docker-compose up -d
docker-compose exec sessionservice ./docker-entrypoint
