#!/usr/bin/env bash

docker build -t loneexile/voidsync:latest .
docker push loneexile/voidsync:latest
