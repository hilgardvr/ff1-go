#!/bin/bash

source config.sh &&
docker build . -t backend-go &&
docker run --network ff1-go_default -p 3000:3000 backend-go
