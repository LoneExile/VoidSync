#!/usr/bin/env bash

docker build -t loneexile/voidsync:latest .
docker-compose up -d --build
